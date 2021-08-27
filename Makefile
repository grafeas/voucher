# Go parameters
GOCMD=cd v2 && go
GORELEASER=goreleaser
GOLANGCI-LINT=cd v2 && golangci-lint
DOCKER=docker
GOPATH?=`echo $$GOPATH`
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
PACKAGES := voucher_server voucher_subscriber voucher_client
CODE=./cmd/
SERVER_NAME=voucher_server
SUBSCRIBER_NAME=voucher_subscriber
CLIENT_NAME=voucher_client
IMAGE_NAME?=voucher

export GO111MODULE=on

.PHONY: clean ensure-deps update-deps system-deps \
	test show-coverage \
	build release snapshot container mocks \
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
	$(GOCMD) tool cover -html=coverage.txt

test:
	$(GOCMD) test ./... -race -coverprofile=coverage.txt -covermode=atomic

lint:
	$(GOLANGCI-LINT) run

lint-new:
	$(GOLANGCI-LINT) run --new-from-rev main

clean:
	$(GOCLEAN)
	@for PACKAGE in $(PACKAGES); do \
		rm -vrf build/$$PACKAGE; \
	done

ensure-deps:
	$(GOCMD) mod download
	$(GOCMD) mod verify

update-deps:
	$(GOCMD) get -u -t all
	$(GOCMD) mod tidy

build: $(PACKAGES)

voucher_client:
	$(GOBUILD) -o ../build/$(CLIENT_NAME) -v $(CODE)$(CLIENT_NAME)

voucher_subscriber:
	$(GOBUILD) -o ../build/$(SUBSCRIBER_NAME) -v $(CODE)$(SUBSCRIBER_NAME)

voucher_server:
	$(GOBUILD) -o ../build/$(SERVER_NAME) -v $(CODE)$(SERVER_NAME)

container:
	$(DOCKER) build -t $(IMAGE_NAME) .

release:
	$(GORELEASER)

snapshot:
	$(GORELEASER) --snapshot

mocks:
	mockgen -source=grafeas/grafeas_service.go -destination=grafeas/mocks/grafeas_service_mock.go package=mocks

test-in-docker:
	docker run -v $(PWD):/go/src/github.com/grafeas/voucher -w /go/src/github.com/grafeas/voucher -e CGO_ENABLED=0  -it golang:1.15.6-alpine go test ./...
