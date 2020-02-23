#!/usr/bin/env bash

UNAME=$(uname)

if [[ $UNAME == "Linux" ]]; then
    apt-get install -y libasound2-dev
fi

mkdir -p bin && GO111MODULE=on go build -o ./bin/gopom ./cmd/gopom