package grace

import (
	"github.com/libra9z/httprouter"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	// PreSignal is the position to add filter before signal
	PreSignal = iota
	// PostSignal is the position to add filter after signal
	PostSignal
	// StateInit represent the application inited
	StateInit
	// StateRunning represent the application is running
	StateRunning
	// StateShuttingDown represent the application is shutting down
	StateShuttingDown
	// StateTerminate represent the application is killed
	StateTerminate
)

var (
	regLock              *sync.Mutex
	runningServers       map[string]*MicroService
	runningServersOrder  []string
	socketPtrOffsetMap   map[string]uint
	runningServersForked bool

	// DefaultReadTimeOut is the HTTP read timeout
	DefaultReadTimeOut time.Duration
	// DefaultWriteTimeOut is the HTTP Write timeout
	DefaultWriteTimeOut time.Duration
	// DefaultMaxHeaderBytes is the Max HTTP Herder size, default is 0, no limit
	DefaultMaxHeaderBytes int
	// DefaultTimeout is the shutdown server's timeout. default is 60s
	DefaultTimeout = 60 * time.Second

	hookableSignals []os.Signal
)

func init() {

	regLock = &sync.Mutex{}
	runningServers = make(map[string]*MicroService)
	runningServersOrder = []string{}
	socketPtrOffsetMap = make(map[string]uint)

	hookableSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
	}
}

// NewServer returns a new graceServer.
func NewServer(ischild bool, socketorder, addr string) (srv *MicroService) {
	regLock.Lock()
	defer regLock.Unlock()

	if len(socketorder) > 0 {
		for i, addr := range strings.Split(socketorder, ",") {
			socketPtrOffsetMap[addr] = uint(i)
		}
	} else {
		socketPtrOffsetMap[addr] = uint(len(runningServersOrder))
	}

	srv = &MicroService{
		Router:  httprouter.New(),
		Server:  &http.Server{},
		wg:      sync.WaitGroup{},
		sigChan: make(chan os.Signal),
		isChild: ischild,
		SignalHooks: map[int]map[os.Signal][]func(){
			PreSignal: {
				syscall.SIGHUP:  {},
				syscall.SIGINT:  {},
				syscall.SIGTERM: {},
			},
			PostSignal: {
				syscall.SIGHUP:  {},
				syscall.SIGINT:  {},
				syscall.SIGTERM: {},
			},
		},
		state:   StateInit,
		Network: "tcp",
	}
	srv.Server.Addr = addr
	srv.Server.ReadTimeout = DefaultReadTimeOut
	srv.Server.WriteTimeout = DefaultWriteTimeOut
	srv.Server.MaxHeaderBytes = DefaultMaxHeaderBytes
	srv.Server.Handler = srv.Router

	runningServersOrder = append(runningServersOrder, addr)
	runningServers[addr] = srv

	return
}
