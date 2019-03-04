# This how we want to name the binary output
BINARY=`basename \`pwd\``

# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"
VERSION=`git describe --tags|sed -e "s/\-/\./g"`
BUILD=`date +%FT%T%z`
COMMIT=`git rev-parse HEAD`
TARGETDIR=$(HOME)/bin
TAGS=-tags "reuseport quic kcp zookeeper etcd consul ping"
TAGSDOCKER=-tags "reuseport quic kcp zookeeper etcd consul ping docker"

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS_DEB=-ldflags=all="-X main.Version=${VERSION} -X main.Date=${BUILD} -X main.Commit=${COMMIT}"
LDFLAGS_REL=-ldflags=all="-w -s -X main.Version=${VERSION} -X main.Date=${BUILD} -X main.Commit=${COMMIT}"

# Builds the project
lrelease:
	env GOOS=linux GOARCH=amd64 go build ${TAGS} ${LDFLAGS_REL} -o ${BINARY}_linux_release
release:
	go build ${TAGS} ${LDFLAGS_REL} -o ${BINARY}
	env GOOS=linux GOARCH=amd64 go build ${TAGS} ${LDFLAGS_REL} -o ${BINARY}_linux_release

ldebug:
	env GOOS=linux GOARCH=amd64 go build ${TAGS} ${LDFLAGS_DEB} -o ${BINARY}_linux_debug

debug:
	go build ${TAGS} ${LDFLAGS_DEB} -o ${BINARY}
	env GOOS=linux GOARCH=amd64 go build ${TAGS} ${LDFLAGS_DEB} -o ${BINARY}_linux_debug

docker:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${TAGSDOCKER} ${LDFLAGS_REL} -o ${BINARY}
	sudo docker build -t libra9z/${BINARY}:${VERSION} .
	
# Installs our project: copies binaries
install:
	go install ${TAGS} ${LDFLAGS_REL}

fb:
	go build ${TAGS} ${LDFLAGS_DEB} -o ${BINARY}
	mv -f ${BINARY} ${TARGETDIR}

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${BINARY}_linux_* ] ; then rm ${BINARY}_linux_*; fi

.PHONY: clean install
