#!/bin/bash

set -e

CLI_URI="https://github.com/sivaprasadreddy/progen/releases/download"
CLI_DIR="$HOME/.progen"
CLI_VERSION="v0.0.1"

CLI_ARCHIVES_FOLDER="${CLI_DIR}/archives"

mkdir -p "$CLI_ARCHIVES_FOLDER"
rm -rf "$CLI_DIR/$CLI_VERSION"
mkdir -p "$CLI_DIR/$CLI_VERSION"
mkdir -p "$CLI_DIR/current"

ostype="linux"
case "$(uname -s)" in
    Darwin*)
        ostype="Darwin"
        ;;
    Linux*)
        ostype="Linux"
        ;;
esac

DOWNLOAD_URI="${CLI_URI}/${CLI_VERSION}/progen_${ostype}_arm64.tar.gz"
echo "Downloading $ostype binary from ${DOWNLOAD_URI}"
cli_zip_file="${CLI_ARCHIVES_FOLDER}/progen-${CLI_VERSION}_${ostype}_arm64.tar.gz"
curl --location --progress-bar "${DOWNLOAD_URI}" > "$cli_zip_file"

tar -xf $cli_zip_file -C $CLI_ARCHIVES_FOLDER "progen"
mv "$CLI_ARCHIVES_FOLDER/progen" "$CLI_DIR/$CLI_VERSION"
ln -sf "$CLI_DIR/$CLI_VERSION/progen" "$CLI_DIR/current/progen"
