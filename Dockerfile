FROM golang:alpine

LABEL maintainer="Vincent Lark <vincent.lark@gmail.com>"

RUN apk add build-base

ADD . /go/src/github.com/vincent/kisslists
WORKDIR /go/src/github.com/vincent/kisslists
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .

FROM alpine:latest
VOLUME /kisslists

WORKDIR /root/
COPY --from=0 /go/src/github.com/vincent/kisslists .
ENTRYPOINT /root/kisslists -database /kisslists/kisslists.sqlite
