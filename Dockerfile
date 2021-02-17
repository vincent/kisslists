FROM golang:alpine AS build

LABEL maintainer="Vincent Lark <vincent.lark@gmail.com>"

RUN apk add build-base

ADD . /go/src/github.com/vincent/kisslists
WORKDIR /go/src/github.com/vincent/kisslists

RUN go get ./... && go mod vendor
RUN CGO_ENABLED=1 go build -a -ldflags '-linkmode external -extldflags "-static"' .

FROM alpine:latest

WORKDIR /root/
COPY --from=build /go/src/github.com/vincent/kisslists .
ENTRYPOINT /root/kisslists -database /kisslists/kisslists.sqlite
