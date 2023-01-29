#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."
RELEASE_DIR="$(pwd)/release/"
CURRENT_VERSION=$(/usr/libexec/PlistBuddy -c "Print :version" info.plist)
echo -e "Current version: $CURRENT_VERSION\nInput new version: "
read -r VERSION
VERSION="v${VERSION//v/}"

test() {
  ./test.sh
}

install_package() {
  ./build.sh
}

sign_binary() {
  # must cd to release directory because signing takes into account directory contents at time of signing.
  # if directory contents change between now and notarization (e.g., because we've packaged into an .alfredworkflow file), the signature verification will fail
  cd "$RELEASE_DIR"
  gon ../release_tools/sign-binary.hcl
  cd -
}

bump_version_and_tag() {
  /usr/libexec/PlistBuddy -c "Set :version ${VERSION//v/}" info.plist
  git add info.plist
  git commit -m "$VERSION"
  git push
  git tag "$VERSION"
  git push origin "$VERSION"
}

copy_to_release_dir() {
  echo "Using directory $RELEASE_DIR to stage release files ..."
  rm -rf "$RELEASE_DIR"
  mkdir -p "$RELEASE_DIR"
  cp -R images alfred-aws-console-services-workflow console-services.yml icon.png info.plist LICENSE README.md "$RELEASE_DIR"
}

PACKAGE_NAME="AWS Console Services.alfredworkflow"
package_release() {
  ditto -ck "$RELEASE_DIR" "$PACKAGE_NAME"
}

notarize_package() {
  gon release_tools/package.hcl
  rm -f "$PACKAGE_NAME"
}

add_version_to_package_name() {
  mv "$PACKAGE_NAME.zip" "AWS Console Services ${VERSION}.alfredworkflow.zip"
}

create_dummy_awgo_updater_file() {
  echo -e "please ignore and/or discard this file! it only exists to make it so awgo's auto-update will detect a new version is available.\nsee here: https://github.com/deanishe/awgo/blob/11f767b094816cd865fa3b396d09023baeaa8ff5/update/github.go#L93-L97" \
      > do-not-download-this-file.alfredworkflow
}

open_github() {
  echo "opening github releases page ..."
  open "https://github.com/rkoval/alfred-aws-console-services-workflow/releases/new?tag=${VERSION}&title=${VERSION}&body=%23%23%20Changes%0A%0AUser-facing%0A-%20TODO"
}

open_finder() {
  open .
}

test
install_package
copy_to_release_dir
sign_binary
package_release
notarize_package
add_version_to_package_name
bump_version_and_tag
create_dummy_awgo_updater_file
open_github
open_finder
