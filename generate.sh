#!/usr/bin/env bash
set -e
which goimports
command -v goimports >> /dev/null || go install golang.org/x/tools/cmd/goimports@v0.24.0 golang.org/x/tools/cmd/goimports
go generate ./...
