#!/bin/bash

# propagate errors, and print commands
set -ex

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# ubuntu stacks need to install libgit2 package manually, as required version is not in apt repo
# macos stacks need to install libgit2 from older homebrew formula
if [ $(uname -s) == "Linux" ]; then
    DEB_FILES=("${THIS_SCRIPT_DIR}/ubuntu/products/"*.deb)
    dpkg -i "${DEB_FILES[0]}"
    ldconfig
elif [ $(uname -s) == "Darwin" ]; then
    FORMULA_COMMIT_HASH="8df6c18106e3e1ad425e01b8ce3d5c8d177a79ce"
    FORMULA_PATH="${THIS_SCRIPT_DIR}/libgit2.rb"
    curl "https://raw.githubusercontent.com/Homebrew/homebrew-core/${FORMULA_COMMIT_HASH}/Formula/libgit2.rb" -o "${FORMULA_PATH}"
    HOMEBREW_NO_INSTALLED_DEPENDENTS_CHECK=1 HOMEBREW_NO_AUTO_UPDATE=1 brew install "${FORMULA_PATH}"
else
    echo "Unsupported OS: $(uname -s)"
    exit 1
fi

# run step in golang
TMP_GOPATH_DIR="$(mktemp -d)"

GO_PACKAGE_NAME="github.com/Holdapp/bitrise-step-jira-build"
FULL_PACKAGE_PATH="${TMP_GOPATH_DIR}/src/${GO_PACKAGE_NAME}"
mkdir -p "${FULL_PACKAGE_PATH}"

rsync -avh --quiet "${THIS_SCRIPT_DIR}/" "${FULL_PACKAGE_PATH}/"

export GOPATH="${TMP_GOPATH_DIR}"
go run "${FULL_PACKAGE_PATH}/main.go"
