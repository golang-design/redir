# Copyright 2021 The golang.design Initiative Authors.
# All rights reserved. Use of this source code is governed
# by a MIT license that can be found in the LICENSE file.

FROM alpine
WORKDIR /app
COPY . .
EXPOSE 8080
CMD ["/app/redir", "-s"]
