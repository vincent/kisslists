# Use golang:alpine using the given platform from buildx
FROM golang:alpine AS build

LABEL maintainer="Vincent Lark <vincent.lark@gmail.com>"

# We'll need make, gcc and friends
RUN apk add build-base

# Use the whole project directory
ADD . /go/src/github.com/vincent/kisslists
WORKDIR /go/src/github.com/vincent/kisslists

# Compile binary
RUN make full

FROM alpine:latest

# Copy only the binary on the final layer
WORKDIR /root/
COPY --from=build /go/src/github.com/vincent/kisslists/dist/kisslists .

# Use the binary as entrypoint
ENTRYPOINT /root/kisslists -database /kisslists/kisslists.sqlite
