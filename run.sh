#!/usr/bin/env bash
set -e
source env.sh
go generate ./...
goimports -w caching/gen-caching.go
export AWS_SDK_LOAD_CONFIG=true
go run main.go "$@"