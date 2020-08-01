#!/bin/bash

echo "Remove Gomfile"
rm -rf Gomfile

echo "Generate dom file in cmd folder"
cd cmd && gom gen gomfile && mv Gomfile ../ | exit 1
# Go back to root folder
cd ../

echo "Install packages"
gom install | exit 1

echo "Build application with gom"
gom build -o build/candles_loader -i ./cmd/ | exit 1

echo "Run tests"
gom test -i ./cmd/ | exit 1