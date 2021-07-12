#!/usr/bin/env bash
set -e
command -v genny >> /dev/null || go get github.com/cheekybits/genny
command -v goimports >> /dev/null || go get golang.org/x/tools/cmd/goimports
go generate ./...
go run generators/gen-caching_importer/main.go
goimports -w caching/gen-caching.go
