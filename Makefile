# Go parameters
GOCMD=go
GODEP=dep
GOPATH?=`echo $$GOPATH`
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTESTARGS?=-race -covermode=atomic
PACKAGES := cli server
CODE=./cmd/
SERVER_NAME=voucher_server
CLI_NAME=voucher_cli
BINARY_UNIX=_unix

.PHONY: clean setup deps test build install

all: clean deps build

install:
	for BINARY_NAME in $(PACKAGES); do cp -v voucher_$$BINARY_NAME $(GOPATH)/bin/voucher_$$BINARY_NAME; done

test:
	$(GOTEST) ./... $(GOTESTARGS)

clean:
	$(GOCLEAN)
	for PACKAGE in $(PACKAGES); do \
		rm -vrf voucher_$$PACKAGE voucher_$$PACKAGE$(BINARY_UNIX); \
	done

deps: setup
	wget -P hack/ https://storage.googleapis.com/container-analysis-v1alpha1/containeranalysis-go.tar.gz
	tar xzvf hack/containeranalysis-go.tar.gz -C vendor

setup:
	$(GOCMD) get github.com/golang/dep/cmd/dep
	$(GODEP) ensure

build: voucher_cli voucher_server

voucher_cli:
	$(GOBUILD) -o $(CLI_NAME) -v $(CODE)$(CLI_NAME)

voucher_server:
	$(GOBUILD) -o $(SERVER_NAME) -v $(CODE)$(SERVER_NAME)

# Cross Compilation
server-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(SERVER_NAME)$(BINARY_UNIX) -v $(CODE)$(SERVER_NAME)

cli-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(CLI_NAME)$(BINARY_UNIX) -v $(CODE)$(CLI_NAME)
