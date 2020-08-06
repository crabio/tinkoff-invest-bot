#!/bin/bash

echo "Build application with gom"
go build -o build/main | exit 1

echo "Run tests"
go test -i | exit 1