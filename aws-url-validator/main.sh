#!/usr/bin/env bash
cd "$(dirname "$0")"

# this holds session info in firefox like the currently opened tabs
./lz4jsoncat $HOME/Library/Application\ Support/Firefox/Profiles/*default*/sessionstore-backups/recovery.jsonlz4 \
  | ./diff-urls-with-last-opened-urls-json.js
