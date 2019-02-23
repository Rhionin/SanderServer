#!/bin/bash
env GOOS=linux GOARM=7 go build cmd/server/*.go