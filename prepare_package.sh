#!/usr/bin/env bash
set -e
cd "$(dirname "$0")/src"
npm install
./generate_items.js
rm -rf node_modules