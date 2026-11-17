#!/usr/bin/env bash
# Packs a Progen Desktop build from cmd/desktop/build/ into its release archive.
#
#   .github/scripts/package-desktop.sh <target> <output-dir>
#
# <target> is macos, windows-amd64, windows-arm64, linux-amd64 or linux-arm64; build it first with
# go tool mygo build -platform <darwin/universal|windows/amd64|windows/arm64|linux/amd64|linux/arm64> ./cmd/desktop
# The archive is named after the version in cmd/desktop/mygo.json, which the release
# workflow sets from the tag:
#
#   progen-desktop-vX.Y.Z-macos.zip           Progen Desktop.app
#   progen-desktop-vX.Y.Z-windows-<arch>.zip  Progen Desktop/ with Progen Desktop.exe and its files
#   progen-desktop-vX.Y.Z-linux-<arch>.tar.gz progen-desktop-vX.Y.Z-linux-<arch>/ with the executable,
#                                             desktop entry and icon, install.sh, and the archive
#                                             install.sh installs from
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <macos|windows-amd64|windows-arm64|linux-amd64|linux-arm64> <output-dir>" >&2
  exit 2
fi
target=$1
mkdir -p "$2"
out=$(cd "$2" && pwd)

repo=$(cd "$(dirname "$0")/../.." && pwd)
desktop="$repo/cmd/desktop"
version=$(sed -n 's/^ *"version": *"\([^"]*\)".*$/\1/p' "$desktop/mygo.json")
if [[ -z "$version" ]]; then
  echo "no version in $desktop/mygo.json" >&2
  exit 1
fi
archive="progen-desktop-v${version}-${target}"

# copy_app_files copies the app's files, which mygo build lists in .mygo-files, from $1 to $2.
copy_app_files() {
  local from=$1 to=$2
  [[ -f "$from/.mygo-files" ]] || { echo "no app in $from: build it first" >&2; exit 1; }
  mkdir -p "$to"
  cp "$from/.mygo-files" "$to/"
  while IFS= read -r name; do
    [[ -n "$name" ]] && cp -R "$from/$name" "$to/"
  done < "$from/.mygo-files"
}

staging=$(mktemp -d)
trap 'rm -rf "$staging"' EXIT

case "$target" in
  macos)
    app="$desktop/build/darwin-universal/Progen Desktop.app"
    [[ -d "$app" ]] || { echo "no app at $app: build it first" >&2; exit 1; }
    rm -f "$out/$archive.zip"
    ditto -c -k --sequesterRsrc --keepParent "$app" "$out/$archive.zip"
    echo "$out/$archive.zip"
    ;;
  windows-amd64 | windows-arm64)
    copy_app_files "$desktop/build/$target" "$staging/Progen Desktop"
    rm -f "$out/$archive.zip"
    (cd "$staging" && zip -qr "$out/$archive.zip" "Progen Desktop")
    echo "$out/$archive.zip"
    ;;
  linux-amd64 | linux-arm64)
    build="$desktop/build/$target"
    copy_app_files "$build" "$staging/$archive"
    # install.sh installs from the archive mygo build made, when it is next to the script.
    cp "$build/install.sh" "$build/progen-desktop-${version}-${target}.tar.gz" "$staging/$archive/"
    rm -f "$out/$archive.tar.gz"
    tar -czf "$out/$archive.tar.gz" -C "$staging" "$archive"
    echo "$out/$archive.tar.gz"
    ;;
  *)
    echo "unknown target $target: use macos, windows-amd64, windows-arm64, linux-amd64 or linux-arm64" >&2
    exit 2
    ;;
esac
