#!/usr/bin/env bash
set -e
go run generators/searcher/*.go $@
go fmt ./...
./generate.sh
go mod tidy