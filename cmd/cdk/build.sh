#!/bin/bash

GOOS=linux GOARCH=arm64 go build -o ./getProgressLambda/bootstrap ./getProgressLambda/main.go