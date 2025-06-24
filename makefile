GOCMD=go


GOBUILD=${GOCMD} build -gcflags=all='-l -N' -ldflags "-s -w"



.PHONY: all clean build

all: clean build

clean:

build:
	${GOBUILD} -o qnoracle main.go