#!/bin/bash
env GOOS=linux GOARCH=arm GOARM=5 go build -o cmd/server/server cmd/server/*.go