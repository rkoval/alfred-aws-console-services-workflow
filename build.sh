#!/usr/bin/env bash
set -e
./generate.sh
GOOS=darwin go build $@
