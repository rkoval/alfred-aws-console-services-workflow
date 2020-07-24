#!/usr/bin/env bash
set -e
export TEST=1
source env.sh
export AWS_REGION=us-west-2
go generate ./...
goimports -w caching/gen-caching.go
UPDATE_SNAPSHOTS=true go test ./... $@