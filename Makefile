# Go parameters
GOCMD=go
GORELEASER=goreleaser
GOLANGCI-LINT=golangci-lint
DOCKER=docker
GOPATH?=`echo $$GOPATH`
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
PACKAGES := voucher_server voucher_client
CODE=./cmd/
SERVER_NAME=voucher_server
CLIENT_NAME=voucher_client
IMAGE_NAME?=voucher

export GO111MODULE=on

.PHONY: clean ensure-deps update-deps system-deps \
	test show-coverage \
	build release snapshot container \
	$(PACKAGES)

all: clean ensure-deps build

# System Dependencies
system-deps:
ifeq ($(shell $(GOCMD) version 2> /dev/null) , "")
	$(error "go is not installed")
endif
ifeq ($(shell $(DOCKER) -v dot 2> /dev/null) , "")
	$(error "docker is not installed")
endif
ifeq ($(shell $(GOLANGCI-LINT) version 2> /dev/null) , "")
	$(error "golangci-lint is not installed")
endif
ifeq ($(shell $(GORELEASER) --version dot 2> /dev/null) , "")
	$(error "goreleaser is not installed")
endif
	$(info "No missing dependencies")

show-coverage: test
	go tool cover -html=coverage.txt

test:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

test-integrations:
	docker ps -a | grep grafeasos || docker run -d -p 8080:8080 --name grafeasos us.gcr.io/grafeas/grafeas-server:v0.1.4
	go test ./grafeasos -grafeasos="http://localhost:8080"
	docker rm -f grafeasos

lint:
	golangci-lint run

lint-new:
	golangci-lint run --new-from-rev master

clean:
	$(GOCLEAN)
	@for PACKAGE in $(PACKAGES); do \
		rm -vrf build/$$PACKAGE; \
	done

ensure-deps:
	$(GOCMD) mod tidy
	$(GOCMD) mod download
	$(GOCMD) mod verify

update-deps:
	$(GOCMD) get -u -t all
	$(GOCMD) mod tidy

build: $(PACKAGES)

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
