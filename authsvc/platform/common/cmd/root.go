package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var RootCmd *cobra.Command

var (
	Address        string
	Http           string
	Https          string
	Consul         string
	Token          string
	ServiceConf    string
	Sdt			   string
	Sda			   string
	Workerid       int
	Debug          bool
	IsChild        bool
	SocketOrder    string
	DebugAddr      string
	ZipkinAddr     string
	KafkaAddr      string
	AppdashAddr    string
	LightstepToken string
)

func init() {
	_, f := filepath.Split(os.Args[0])
	RootCmd = &cobra.Command{Use: f}
	RootCmd.PersistentFlags().StringVar(&Address, "addr", "", "config this service address format(ip:port).")
	RootCmd.PersistentFlags().StringVar(&Http, "http", "", "config this service http address format(ip:port).")
	RootCmd.PersistentFlags().StringVar(&Https, "https", "", "config this service https address format(ip:port).")
	RootCmd.PersistentFlags().StringVar(&Consul, "consul", "", "config consul address format(ip:port).")
	RootCmd.PersistentFlags().StringVar(&Token, "token", "", "config consul acl token.")
	RootCmd.PersistentFlags().StringVar(&ServiceConf, "file", "", "service config file.")
	RootCmd.PersistentFlags().StringVar(&Sda, "sda", "", "service discovery address.")
	RootCmd.PersistentFlags().StringVar(&Sdt, "sdt", "", "service discovery type,can by cansul,etcd.")

	RootCmd.PersistentFlags().IntVar(&Workerid, "worker", 1, "config worker id for distribute .")
	RootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "enable debug mode true or false")
	RootCmd.PersistentFlags().StringVar(&DebugAddr, "debug.addr", ":18080", "Debug and metrics listen address")
	//RootCmd.PersistentFlags().StringVar(&HttpAddr,"http.addr", ":8081", "HTTP listen address")
	//RootCmd.PersistentFlags().StringVar(&RecordAddr,"collect.addr", "", "zipkin collect address")
	RootCmd.PersistentFlags().StringVar(&ZipkinAddr, "zipkin.addr", "", "Enable Zipkin tracing via a Kafka server host:port")
	RootCmd.PersistentFlags().StringVar(&KafkaAddr, "kafka.addr", "", "Enable Kafka server host:port")
	RootCmd.PersistentFlags().StringVar(&AppdashAddr, "appdash.addr", "", "Enable Appdash tracing via an Appdash server host:port")
	RootCmd.PersistentFlags().StringVar(&LightstepToken, "lightstep.token", "", "Enable LightStep tracing via a LightStep access token")

	RootCmd.PersistentFlags().BoolVar(&IsChild, "graceful", false, "listen on open fd (after forking)")
	RootCmd.PersistentFlags().StringVar(&SocketOrder, "socketorder", "", "previous initialization order - used when more than one listener was started")

	//RootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
