.PHONY: build docker

TARGET_DIR ?= ./

# TEST_ARGS = -v | grep -c RUN
DOCKER_OPTS ?= --rm
DOCKER_REPO ?=github.com/influenzanet/survey-repository

build:
	go build -o ./survey-repository ./cmd/survey-repository

run:
	go run ./cmd/survey-repository 

server:
	go run ./cmd/survey-repository server