#!/bin/sh

echo "publishing docker image" \
    && docker build -t allyouneedisgnu/kisslists . \
    && docker push allyouneedisgnu/kisslists
