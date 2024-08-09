#!/usr/bin/env bash
set -e
command -v goimports >> /dev/null || go install golang.org/x/tools/cmd/goimports@0.24.0 golang.org/x/tools/cmd/goimports
go generate ./...
