#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $DIR/cmd/labkube
CGO_ENABLED=0 go build -o labkube -a -installsuffix cgo -ldflags "-s" *.go

docker build -t dbenque/labkube:v1 .
docker push dbenque/labkube:v1