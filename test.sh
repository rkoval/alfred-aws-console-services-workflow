#!/usr/bin/env bash
set -e
export TEST=1
source env.sh
go generate ./...
UPDATE_SNAPSHOTS=true go test ./... $@