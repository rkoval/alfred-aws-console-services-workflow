#!/usr/bin/env bash

root=".alfred"
export alfred_version=
export alfred_workflow_bundleid="com.ryankoval.awsconsoleservices"
export alfred_workflow_data="${root}/data"
export alfred_workflow_cache="${root}/cache"

go run main.go "$1"