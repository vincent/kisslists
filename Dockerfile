FROM golang:alpine

LABEL maintainer="Vincent Lark <vincent.lark@gmail.com>"

ADD . /go/src/github.com/vincent/kisslists
WORKDIR /go/src/github.com/vincent/kisslists

RUN apk add build-base
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' .

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/src/github.com/vincent/kisslists .
ENTRYPOINT /root/kisslists -database /kisslists.sqlite
