#!/usr/bin/env bash
set -e
./build.sh

# mock out AWS env vars in CI so that aws-sdk-go is happy
if [ ! -f "$HOME/.aws/config" ]; then
  mkdir -p "$HOME/.aws"
  echo "[default]" > "$HOME/.aws/config"
fi
export AWS_ACCESS_KEY_ID=AAAAAAAAAAAAAAAAAAAA
export AWS_SECRET_ACCESS_KEY=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
./test.sh