#!/bin/bash

set -e # immediately fail on error

bash ./build.sh
cdk deploy
