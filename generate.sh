#!/usr/bin/env bash
set -e
command -v goimports >> /dev/null || go install golang.org/x/tools/cmd/goimports@latest golang.org/x/tools/cmd/goimports
go generate ./...
