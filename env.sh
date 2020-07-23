#!/usr/bin/env bash

export alfred_workflow_bundleid="com.ryankoval.awsconsoleservices"
export alfred_version=

if [ "$(uname)" == "Darwin" ]; then
  # make this mirror where alfred stores cache/data
  if [ -d "$HOME/Library/Application Support/Alfred 3" ]; then
    alfred_workflow_cache="$HOME/Library/Caches/com.runningwithcrayons.Alfred-3/Workflow Data/$alfred_workflow_bundleid"
    alfred_workflow_data="$HOME/Library/Application Support/Alfred 3/Workflow Data/$alfred_workflow_bundleid"
  else
    echo "TODO what is alfred 4 cache directory?"
    exit 1
  fi
  export alfred_workflow_data
  export alfred_workflow_cache
else
  # CI won't have alfred directories, so just create in repo root
  root="$(git rev-parse --show-toplevel)/.alfred"
  export alfred_workflow_data="${root}/data"
  export alfred_workflow_cache="${root}/cache"
fi
