#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."
CURRENT_VERSION=$(/usr/libexec/PlistBuddy -c "Print :version" info.plist)
echo -e "Current version: $CURRENT_VERSION\nInput new version: "
read -r VERSION
VERSION="v${VERSION//v/}"

install_package() {
  go build
}

bump_version_and_tag() {
  /usr/libexec/PlistBuddy -c "Set :version ${VERSION//v/}" info.plist
  git add info.plist
  git commit -m "$VERSION"
  git push
  git tag "$VERSION"
  git push origin "$VERSION"
}

package_release() {
  local tmpdir
  tmpdir=$(mktemp -d)
  echo "Using tmp dir $tmpdir to stage release files ..."
  cp -R images alfred-aws-console-services-workflow console-services.yml icon.png info.plist LICENSE README.md "$tmpdir/"
  ditto -ck "$tmpdir" "AWS Console Services ${VERSION}.alfredworkflow"
}

open_github() {
  echo "opening github releases page ..."
  open "https://github.com/rkoval/alfred-aws-console-services-workflow/releases/new?tag=${VERSION}&title=${VERSION}&body=Fill out the release"
}

install_package
bump_version_and_tag
package_release
open_github