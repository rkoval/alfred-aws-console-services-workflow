#!/usr/bin/env bash
set -e
source env.sh
./generate.sh
export AWS_SDK_LOAD_CONFIG=true
go run main.go "$@"