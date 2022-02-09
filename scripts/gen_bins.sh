#!/bin/sh

export CGO_ENABLED=0

GOOS=linux GOARCH=amd64 go build -o ./bin/ycloud-cm-downloader-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o ./bin/ycloud-cm-downloader-linux-arm64 .
GOOS=darwin GOARCH=amd64 go build -o ./bin/ycloud-cm-downloader-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o ./bin/ycloud-cm-downloader-darwin-arm64 .