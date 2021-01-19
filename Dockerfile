FROM golang:latest
ADD . /go/src/github.com/vincent/sharedlists
WORKDIR /go/src/github.com/vincent/sharedlists
RUN go install -v
ENTRYPOINT ["sharedlists", "--database /sharedlists.db"]
