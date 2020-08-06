#!/bin/bash

echo "Build application with gom"
GOOS=linux CGO_ENABLED=0 go build -o build/main | exit 1

echo "Run tests"
go test -i | exit 1

echo "Build docker"
docker build --force-rm -t candles-loader .