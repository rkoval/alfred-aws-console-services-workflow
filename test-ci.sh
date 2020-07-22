#!/usr/bin/env bash

# mock out AWS env vars in CI so that aws-sdk-go is happy
export AWS_ACCESS_KEY_ID=AAAAAAAAAAAAAAAAAAAA
export AWS_SECRET_ACCESS_KEY=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
export AWS_DEFAULT_REGION=us-west-2
./test.sh