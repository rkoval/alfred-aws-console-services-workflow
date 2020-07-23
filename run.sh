#!/usr/bin/env bash
source env.sh
go generate ./...
go run main.go "$@"