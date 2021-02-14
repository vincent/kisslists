#!/bin/sh

echo "publishing docker image" \
    && docker build -t allyouneedisgnu/sharedlists . \
    && docker push allyouneedisgnu/sharedlists
