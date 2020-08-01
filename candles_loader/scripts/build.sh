#!/bin/bash

echo "Build application with gom"
go build -o build/candles_loader | exit 1

echo "Run tests"
go test -i | exit 1