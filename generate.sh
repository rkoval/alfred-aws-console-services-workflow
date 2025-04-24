#!/usr/bin/env bash
set -e
command -v goimports >> /dev/null || go install golang.org/x/tools/cmd/goimports@v0.32.0
go generate ./...
