# Copyright 2020 The golang.design Initiative authors.
# All rights reserved. Use of this source code is governed
# by a MIT license that can be found in the LICENSE file.

FROM golang:1.14.9
WORKDIR /app
COPY . .
RUN go build -o redir -ldflags \
  "-X main.Version=v$(git describe --always --tags) \
   -X main.BuildTime=$(date +%FT%T%z)"
# @changkun: could use alpine instead, maybe in the future.
EXPOSE 8080
CMD ["/app/redir", "-s"]