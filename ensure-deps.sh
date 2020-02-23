#!/usr/bin/env bash

UNAME=$(uname)

echo Installing deps for $UNAME

if [[ $UNAME == "Linux" ]]; then
    apt-get install -y libasound2-dev
fi