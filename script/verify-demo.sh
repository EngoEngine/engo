#!/usr/bin/env bash

set -o errexit
# set -o xtrace # for debug uncomment

projectDir=$(cd "$(dirname "${0}")/.." && pwd)
# shellcheck source=script/util.sh
source "${projectDir}/script/util.sh" || source ./util.sh

echo "VERIFYING demos..."

# These can fail without us minding it
blacklist="demoutils"

for target in "${projectDir}"/demos/*/; do
  # Ignore the directory if it's in the blacklist
  if [ $blacklist == "${target##*/}" ]; then
    println "SKIP ${target}"
    continue
  fi

  println "VERIFYING ${target}..."
  "${projectDir}/script/go-build.sh" -tags demo "${target}"
  go clean
done
