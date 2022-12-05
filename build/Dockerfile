FROM golang:1.19

WORKDIR /go/src/force

COPY . /go/src/force

ARG protocVersion=3.16.3
ARG golangciVersion=1.50.1
ARG bufVersion=1.9.0

ENV GO111MODULE=on

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get -y install git unzip build-essential autoconf libtool iputils-ping net-tools dnsutils

RUN curl -ssL https://github.com/protocolbuffers/protobuf/releases/download/v$protocVersion/protoc-$protocVersion-linux-x86_64.zip -o /tmp/protoc-$protocVersion-linux-x86_64.zip && \
    unzip -d '/usr/local' /tmp/protoc-$protocVersion-linux-x86_64.zip 'bin/*' 'include/*' && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v$golangciVersion && \
    curl -sSL https://github.com/bufbuild/buf/releases/download/v$bufVersion/buf-Linux-x86_64 -o "/usr/local/bin/buf" && \
    chmod +x "/usr/local/bin/buf"

RUN MAKEFILE_ENV=local make install && \
    MAKEFILE_ENV=local make download && \
    MAKEFILE_ENV=local make vendor && \
    MAKEFILE_ENV=local make gen

CMD go run ./cmd/migrate/main.go
CMD air -c build/air.conf
