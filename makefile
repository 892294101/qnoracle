GOCMD=go


GO_MAC_BUILD=${GOCMD} build -gcflags=all='-l -N' -ldflags "-s -w"
GO_LINUX_BUILD=GOOS=linux GOARCH=amd64 ${GOCMD} build -gcflags=all='-l -N' -ldflags "-s -w"


.PHONY: all clean build

all: clean build

clean:

build:
	${GO_MAC_BUILD} -o qnoracle-mac main.go
	${GO_LINUX_BUILD} -o qnoracle-linux main.go