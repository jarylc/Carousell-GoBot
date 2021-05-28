#!/bin/ash
apk --no-cache add npm go-bindata
cd chrono && ./prepare.sh && cd ..

GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o carousell-gobot.windows-amd64.exe
GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o carousell-gobot.linux-amd64
GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o carousell-gobot.linux-arm64
GOOS=linux GOARCH=arm go build -ldflags="-w -s" -o carousell-gobot.linux-arm-v7
