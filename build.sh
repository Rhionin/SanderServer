#!/bin/bash

set -e # immediately fail on error

GOOS=linux GOARCH=arm64 go build -o ./getProgressLambda/bootstrap ./getProgressLambda/main.go
GOOS=linux GOARCH=arm64 go build -o ./pushUpdatesLambda/bootstrap ./pushUpdatesLambda/main.go