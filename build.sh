#!/usr/bin/env bash

UNAME=$(uname)

echo Running build for $UNAME

if [[ $UNAME == "Linux" ]]; then
    apt-get install -y libasound2-dev
fi

mkdir -p bin && GO111MODULE=on go build -o ./bin/gopom ./cmd/gopom