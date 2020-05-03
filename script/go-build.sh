#!/usr/bin/env bash

set -o errexit

projectDir=$(cd "$(dirname "${0}")/.." && pwd)
# shellcheck source=script/util.sh
source "${projectDir}/script/util.sh" || source ./util.sh

if [ "${ENV_TYPE}" == "mobile" ] && [ "${OS_FAMILY}" == "macos" ]; then
  "$(go env GOPATH)/bin/gomobile" build -target ios "${@}"
elif [ "${ENV_TYPE}" == "mobile" ]; then
  "$(go env GOPATH)/bin/gomobile" build -target android "${@}"
elif [ "${ENV_TYPE}" == "browser" ]; then
  GOOS=js GOARCH=wasm go build "${@}"
else
  go build "${@}"
fi
