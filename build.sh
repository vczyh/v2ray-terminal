#!/bin/bash

# Linux
export GOOS=linux
export GOARCH=amd64
go build -o releases/v2rayT-linux-amd64


# Mac
export GOOS="darwin"
export GOARCH="amd64"
go build -o releases/v2rayT-darwin-amd64

# Mac ARM
export GOOS="darwin"
export GOARCH="arm64"
go build -o releases/v2rayT-darwin-arm64

