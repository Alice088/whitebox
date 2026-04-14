#!/bin/bash

APP=whitebox

mkdir -p dist

echo "Building..."

GOOS=linux GOARCH=amd64 go build -o dist/$APP-linux
GOOS=darwin GOARCH=arm64 go build -o dist/$APP-mac
GOOS=windows GOARCH=amd64 go build -o dist/$APP.exe

echo "Done → dist/"