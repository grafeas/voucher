# Go parameters
GOCMD=go
GODEP=dep
DOCKER=docker
GOPATH?=`echo $$GOPATH`
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
PACKAGES := cli server
CODE=./cmd/
SERVER_NAME=voucher_server
CLI_NAME=voucher_cli
CLIENT_NAME=voucher_client
IMAGE_NAME?=voucher
BINARY_UNIX=_unix

.PHONY: clean setup test build install voucher_cli

all: clean build

install:
	for BINARY_NAME in $(PACKAGES); do cp -v voucher_$$BINARY_NAME $(GOPATH)/bin/voucher_$$BINARY_NAME; done

show-coverage: test
	go tool cover -html=coverage.txt

test:
	./test.sh

clean:
	$(GOCLEAN)
	for PACKAGE in $(PACKAGES); do \
		rm -vrf voucher_$$PACKAGE voucher_$$PACKAGE$(BINARY_UNIX); \
	done

update-deps:
	$(GOCMD) get github.com/golang/dep/cmd/dep
	$(GODEP) ensure

build: voucher_cli voucher_server voucher_client

voucher_cli:
	$(GOBUILD) -o $(CLI_NAME) -v $(CODE)$(CLI_NAME)

voucher_client: $(wildcard cmd/voucher_client/*.go)
	$(GOBUILD) -o $(CLIENT_NAME) -v $(CODE)$(CLIENT_NAME)

voucher_server:
	$(GOBUILD) -o $(SERVER_NAME) -v $(CODE)$(SERVER_NAME)

container:
	$(DOCKER) build -t $(IMAGE_NAME) .

# Cross Compilation
server-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(SERVER_NAME)$(BINARY_UNIX) -v $(CODE)$(SERVER_NAME)

cli-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(CLI_NAME)$(BINARY_UNIX) -v $(CODE)$(CLI_NAME)
