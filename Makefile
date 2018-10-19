# Go parameters
GOCMD=go
GODEP=dep
GORELEASER=goreleaser
DOCKER=docker
GOPATH?=`echo $$GOPATH`
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
PACKAGES := voucher_cli voucher_server voucher_client
CODE=./cmd/
SERVER_NAME=voucher_server
CLI_NAME=voucher_cli
CLIENT_NAME=voucher_client
IMAGE_NAME?=voucher

.PHONY: clean update-deps system-deps \
	test show-coverage \
	build release snapshot container \
	$(PACKAGES)

all: clean build

# System Dependencies
system-deps:
ifeq ($(shell $(GOCMD) version 2> /dev/null) , "")
	$(error "go is not installed")
endif
ifeq ($(shell $(DOCKER) -v dot 2> /dev/null) , "")
	$(error "docker is not installed")
endif
ifeq ($(shell $(GODEP) version dot 2> /dev/null) , "")
	$(error "dep is not installed")
endif
ifeq ($(shell $(GORELEASER) version dot 2> /dev/null) , "")
	$(error "goreleaser is not installed")
endif
	$(info "No missing dependencies")

show-coverage: test
	go tool cover -html=coverage.txt

test:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

clean:
	$(GOCLEAN)
	@for PACKAGE in $(PACKAGES); do \
		rm -vrf build/$$PACKAGE; \
	done

update-deps:
	$(GOCMD) get github.com/golang/dep/cmd/dep
	$(GODEP) ensure

build: $(PACKAGES)

voucher_cli:
	$(GOBUILD) -o build/$(CLI_NAME) -v $(CODE)$(CLI_NAME)

voucher_client:
	$(GOBUILD) -o build/$(CLIENT_NAME) -v $(CODE)$(CLIENT_NAME)

voucher_server:
	$(GOBUILD) -o build/$(SERVER_NAME) -v $(CODE)$(SERVER_NAME)

container:
	$(DOCKER) build -t $(IMAGE_NAME) .

release:
	$(GORELEASER)

snapshot:
	$(GORELEASER) --snapshot
