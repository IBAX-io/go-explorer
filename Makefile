export GOPROXY=https://goproxy.io
export GO111MODULE=on

HOMEDIR := $(shell pwd)

all: mod build

mod:
	go mod tidy -v

build:
	bash $(HOMEDIR)/build.sh

initdb:
	go-ibax-explorer initDatabase
start:
	go-ibax-explorer start

startup: initdb start
