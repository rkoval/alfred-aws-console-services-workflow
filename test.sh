#!/usr/bin/env bash
source env.sh
go generate ./...
UPDATE_SNAPSHOTS=true go test ./... $@