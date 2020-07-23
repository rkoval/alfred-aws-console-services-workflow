#!/usr/bin/env bash
set -e
source env.sh
go generate ./...
UPDATE_SNAPSHOTS=true go test ./... $@