#!/bin/bash
rm getProgressLambda.zip
env GOOS=linux GOARCH=arm64 go build -o bootstrap cmd/getProgressLambda/*.go
zip getProgressLambda.zip bootstrap
rm bootstrap
