DIR=$(PWD)
NAME ?= survey-repository
GOBINS ?= .
GO=go
DOCKER_IMAGE ?= survey-repository
DOCKER_TAG ?= latest
GOPATH ?= /go
VERSION ?= $(shell go run build/meta.go version)
GO_PKG ?= github.com/influenzanet/$(NAME)
DOCKER_NAME ?= $(DOCKER_IMAGE):$(DOCKER_TAG)
GO_LDFLAGS ?=

.PHONY: build docker

TARGET_DIR ?= ./

# TEST_ARGS = -v | grep -c RUN
DOCKER_OPTS ?= --rm
DOCKER_REPO ?=github.com/influenzanet/survey-repository

build:
	go build -o ./survey-repository ./cmd/survey-repository

run:
	go run ./cmd/survey-repository 

version:
	@echo $(VERSION)

server:
	go run ./cmd/survey-repository server

_docker_install: 
	CGO_ENABLED=1 $(GO) build -o $(GOPATH)/$(NAME) -ldflags '-extldflags "-static" $(GO_LDFLAGS)' -tags netgoex $(DIR)/cmd/$(NAME) 

docker:
	go run build/meta.go
	docker build -f build/debian/Dockerfile --build-arg NAME=$(NAME) -t $(DOCKER_NAME) . 

docker-export:
	mkdir -p artifacts
	go run build/meta.go
	docker buildx build --output type=local,dest=artifacts/$(VERSION) -f build/debian/Dockerfile .
