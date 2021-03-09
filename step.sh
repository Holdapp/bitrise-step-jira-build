#!/bin/bash

# propagate errors, and print commands
set -ex

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# ubuntu stacks need to install libgit2 package manually, as required version is not in apt repo
if [ $(uname -s) == "Linux" ]; then
    DEB_FILES=("${THIS_SCRIPT_DIR}/ubuntu/products/"*.deb)
    dpkg -i "${DEB_FILES[0]}"
    ldconfig
fi

# run step in golang
TMP_GOPATH_DIR="$(mktemp -d)"

GO_PACKAGE_NAME="github.com/Holdapp/bitrise-step-jira-build"
FULL_PACKAGE_PATH="${TMP_GOPATH_DIR}/src/${GO_PACKAGE_NAME}"
mkdir -p "${FULL_PACKAGE_PATH}"

rsync -avh --quiet "${THIS_SCRIPT_DIR}/" "${FULL_PACKAGE_PATH}/"

export GOPATH="${TMP_GOPATH_DIR}"
go run "${FULL_PACKAGE_PATH}/main.go"
