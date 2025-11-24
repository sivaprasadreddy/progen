#!/bin/bash

set -e

APP="progen"
CLI_URI="https://github.com/sivaprasadreddy/${APP}/releases/download"
CLI_DIR="$HOME/.${APP}"
CLI_VERSION="v1.0.0"
CLI_ARCHIVES_FOLDER="${CLI_DIR}/archives"

mkdir -p "${CLI_ARCHIVES_FOLDER}"
rm -rf "${CLI_DIR}/${CLI_VERSION}"
mkdir -p "${CLI_DIR}/${CLI_VERSION}"
mkdir -p "${CLI_DIR}/current"

ostype="linux"
case "$(uname -s)" in
    Darwin*)
        ostype="darwin"
        ;;
    Linux*)
        ostype="linux"
        ;;
esac

DOWNLOAD_URI="${CLI_URI}/${CLI_VERSION}/${APP}_${ostype}_arm64.tar.gz"
echo "Downloading $ostype binary from ${DOWNLOAD_URI}"
CLI_ZIP_FILE="${CLI_ARCHIVES_FOLDER}/${APP}-${CLI_VERSION}_${ostype}_arm64.tar.gz"
curl --location --progress-bar "${DOWNLOAD_URI}" > "${CLI_ZIP_FILE}"

tar -xf ${CLI_ZIP_FILE} -C ${CLI_ARCHIVES_FOLDER} "${APP}"
mv "${CLI_ARCHIVES_FOLDER}/${APP}" "${CLI_DIR}/${CLI_VERSION}"
ln -sf "${CLI_DIR}/${CLI_VERSION}/${APP}" "${CLI_DIR}/current/${APP}"
