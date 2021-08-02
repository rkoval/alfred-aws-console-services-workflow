#!/usr/bin/env bash
set -e
export TEST=1
source env.sh
export AWS_SHARED_CREDENTIALS_FILE="../tests/test_aws_credentials_file"
export AWS_CONFIG_FILE="../tests/test_aws_config_file"
./generate.sh
UPDATE_SNAPSHOTS=true go test ./... $@