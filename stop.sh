#!/bin/bash
CWD=$(cd "$(dirname $0)";pwd)
"$CWD"/easydarwin stop
"$CWD"/easydarwin uninstall 
#   wget https://go.dev/dl/go1.18.3.linux-amd64.tar.gz
#    11  tar -C /usr/local -xzf go1.18.3.linux-amd64.tar.gz 
# GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build .
# docker build -t arm-easydarwin:1.0.6 .
# docker save -o arm-easydarwin1.0.6.tar arm-easydarwin:1.0.6