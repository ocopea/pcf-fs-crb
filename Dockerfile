# Copyright (c) [2017] Dell Inc. or its subsidiaries. All Rights Reserved.
FROM golang:1.8

RUN curl -o /usr/local/bin/swagger -L'#' https://github.com/go-swagger/go-swagger/releases/download/0.9.0/swagger_$(echo `uname`|tr '[:upper:]' '[:lower:]')_amd64
RUN chmod +x /usr/local/bin/swagger
RUN swagger version
RUN go get github.com/tools/godep
RUN mkdir $GOPATH/src/crb
ADD build_crb.sh /usr/local/bin/build_crb.sh
RUN chmod +x /usr/local/bin/build_crb.sh
ENTRYPOINT build_crb.sh
