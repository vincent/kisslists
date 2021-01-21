FROM golang:latest

LABEL maintainer="Vincent Lark <vincent.lark@gmail.com>"

ADD . /go/src/github.com/vincent/sharedlists

WORKDIR /go/src/github.com/vincent/sharedlists

RUN go install -v

ENV THEME=""

ENTRYPOINT sharedlists -database /sharedlists.sqlite -theme ${THEME:-default}
