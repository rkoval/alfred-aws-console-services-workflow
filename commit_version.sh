#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"
VERSION=v$(/usr/libexec/PlistBuddy -c "Print :version" info.plist)
git add info.plist
git commit -m "$VERSION"
git push
git tag "$VERSION"
git push origin "$VERSION"
echo "opening github releases page ..."
open "https://github.com/rkoval/alfred-aws-console-services-workflow/releases/new?tag=${VERSION}&title=${VERSION}&body=Fill out the release"