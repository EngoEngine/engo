#!/usr/bin/env bash

set -o errexit

projectDir=$(cd "$(dirname "${0}")/.." && pwd)
# shellcheck source=script/util.sh
source "${projectDir}/script/util.sh" || source ./util.sh

if [ "${ENV_TYPE}" == "browser" ] && [ "${OS_FAMILY}" == "windows" ]; then
  ls "$(go env GOROOT)\misc\wasm"
  GOOS=js GOARCH=wasm go test -exec="$(go env GOROOT)\misc\wasm\go_js_wasm_exec" "${@}"
elif [ "${ENV_TYPE}" == "browser" ]; then
  GOOS=js GOARCH=wasm go test -exec="$(go env GOROOT)/misc/wasm/go_js_wasm_exec" "${@}"
else
  go test "${@}"
fi
