# Copyright 2021 The golang.design Initiative Authors.
# All rights reserved. Use of this source code is governed
# by a MIT license that can be found in the LICENSE file.

VERSION = $(shell git describe --always --tags)
NAME = redir
BUILD_FLAGS = -o $(NAME) -mod=vendor

all:
	go build $(BUILD_FLAGS)
run:
	./$(NAME) -s
build:
	CGO_ENABLED=0 GOOS=linux go build $(BUILD_FLAGS)
	docker build -f docker/Dockerfile -t $(NAME):latest .
up:
	docker-compose -f docker/docker-compose.yml up -d
down:
	docker-compose -f docker/docker-compose.yml down
status:
	docker-compose -f docker/docker-compose.yml ps -a
clean:
	rm -rf $(NAME)
	docker rmi -f $(shell docker images -f "dangling=true" -q) 2> /dev/null; true
	docker rmi -f $(NAME):latest 2> /dev/null; true
.PHONY: all run build up down clean