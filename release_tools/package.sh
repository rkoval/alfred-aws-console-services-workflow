#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."
CURRENT_VERSION=$(/usr/libexec/PlistBuddy -c "Print :version" info.plist)
echo -e "Current version: $CURRENT_VERSION\nInput new version: "
read -r VERSION
VERSION="v${VERSION//v/}"

install_package() {
  ./build.sh
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
  open "https://github.com/rkoval/alfred-aws-console-services-workflow/releases/new?tag=${VERSION}&title=${VERSION}&body=**_If%20you%20missed%20it%20and%20you're%20upgrading%20from%202.x%2C%20you'll%20need%20to%20see%20%5Bthe%203.0.0%20release%5D(https%3A%2F%2Fgithub.com%2Frkoval%2Falfred-aws-console-services-workflow%2Freleases%2Ftag%2Fv3.0.0)%20for%20upgrade%20information_**%0A%0A%23%23%20Changes%0A%0AUser-facing%0A-%20TODO"
}

install_package
bump_version_and_tag
package_release
open_github