# Copyright 2020 The golang.design Initiative authors.
# All rights reserved. Use of this source code is governed
# by a MIT license that can be found in the LICENSE file.

VERSION = $(shell git describe --always --tags)
BUILDTIME = $(shell date +%FT%T%z)
GOPATH=$(shell go env GOPATH)
IMAGE = redir
BINARY = redir.app
TARGET = -o $(BINARY)
BUILD_SETTINGS = -ldflags="-X main.Version=$(VERSION) -X main.BuildTime=$(BUILDTIME)"
BUILD_FLAGS = $(TARGET) $(BUILD_SETTINGS) -mod=vendor

all: native
native:
	go build $(BUILD_FLAGS)
run:
	./$(BINARY) -s
build:
	GOOS=linux go build $(BUILD_FLAGS)
	docker build -t $(IMAGE):$(VERSION) -t $(IMAGE):latest -f docker/Dockerfile .
up: down
	docker-compose -f docker/docker-compose.yml up -d
down:
	docker-compose -f docker/docker-compose.yml down
test:
	mkdir -p build
	go test -cover -coverprofile=build/cover.test -v ./...
	go tool cover -html=build/cover.test -o build/cover.html
clean: down
	rm redir.app
	docker rmi -f $(shell docker images -f "dangling=true" -q) 2> /dev/null; true
	docker rmi -f $(IMAGE):latest $(IMAGE):$(VERSION) 2> /dev/null; true
.PHONY: all native run build up down test clean