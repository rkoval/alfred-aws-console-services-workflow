#!/usr/bin/env bash

export alfred_workflow_bundleid="com.ryankoval.awsconsoleservices"
export alfred_version=3.8
export alfred_workflow_version=1.0

if [ "$(uname)" == "Darwin" ] && [ -z "$TEST" ]; then
  # make this mirror where alfred stores cache/data
  if [ -d "$HOME/Library/Application Support/Alfred" ]; then
    data_dir="com.runningwithcrayons.Alfred"
    cache_dir="Alfred"
  else
    data_dir="com.runningwithcrayons.Alfred-3"
    cache_dir="Alfred 3"
  fi
  export alfred_workflow_data="$HOME/Library/Caches/$cache_dir/Workflow Data/$alfred_workflow_bundleid"
  export alfred_workflow_cache="$HOME/Library/Application Support/$data_dir/Workflow Data/$alfred_workflow_bundleid"
else
  # CI won't have alfred directories, so just create in repo root
  root="$(git rev-parse --show-toplevel)/.alfred"
  export alfred_workflow_data="${root}/data"
  export alfred_workflow_cache="${root}/cache"
fi
