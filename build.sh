#!/bin/bash

set -e # immediately fail on error

GOOS=linux GOARCH=arm64 go build -o ./cmd/getProgressLambda/bootstrap ./cmd/getProgressLambda/main.go
GOOS=linux GOARCH=arm64 go build -o ./cmd/pushUpdatesLambda/bootstrap ./cmd/pushUpdatesLambda/main.go