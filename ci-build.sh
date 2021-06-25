#!/bin/ash
GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o "$(basename "${PWD}").windows-amd64.exe"
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o "$(basename "${PWD}").linux-amd64"
GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o "$(basename "${PWD}").linux-arm64"
GOOS=linux GOARCH=arm go build -ldflags="-w -s" -o "$(basename "${PWD}").linux-arm-v7"
