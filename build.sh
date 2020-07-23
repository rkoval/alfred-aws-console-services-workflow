#!/usr/bin/env bash
set -e
command -v genny >> /dev/null || go get github.com/cheekybits/genny
go generate ./...
go build
