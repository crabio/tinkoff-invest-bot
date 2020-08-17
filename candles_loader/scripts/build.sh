#!/bin/bash

echo "Build application with gom"
go build -o build/main | exit 1

echo "Run tests"
go test -v -coverpkg=./... -coverprofile=profile.cov ./...

echo "Run test coverage"
go tool cover -func profile.cov