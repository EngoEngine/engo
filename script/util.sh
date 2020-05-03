#!/usr/bin/env bash

set -o errexit

lowercase() { printf "%s" "${1}" | tr '[:upper:]' '[:lower:]'; }
println() { printf "%s\n" "${*}"; }
info() { println "[INFO] ${1}"; }
error() { println "[ERROR] ${1}"; }
warning() { println "[WARNING] ${1}"; }
unsupported() { warning "${1} not supported"; }
notImplemented() { warning "${1} is not yet implemented on ${OS_FAMILY}/${ENV_TYPE} environment"; }

fatal() {
  if [ -n "${1}" ]; then
    println "[FATAL] ${1}"
  fi
  exit 1
}

notReady() {
  notImplemented "${1}"
  exit 1
}

checkCommand() {
  if ! [ -x "$(command -v "${1}")" ]; then
    fatal "'${1}' is not installed."
  fi
}

checkPOSIX() {
  checkCommand uname
  checkCommand id
  checkCommand cut
  checkCommand grep
  checkCommand tr
}

tryCommand() {
  if ! [ -x "$(command -v "${1}")" ]; then
    println "'${1}' is not installed."
  fi
}

printPlatform() {
  println "OS: ${OS}"
  println "OS_FAMILY: ${OS_FAMILY}"
}

printENV() {
  printPlatform
  println "ENV_TYPE: ${ENV_TYPE}"
  println "GOOS: $(go env GOOS)"
  println "GOARCH: $(go env GOARCH)"
  println "GOPATH: $(go env GOPATH)"
  ls "$(go env GOROOT)/misc/wasm"
  tryCommand git version
  tryCommand go version
  tryCommand node --version
  tryCommand npm --version
}

detectPlatform() {
  ENV_TYPE=${ENV_TYPE:-destop}
  KERNEL_NAME=$(lowercase "$(uname -s 2>/dev/null)" || printf unknown)
  OS=unknown
  OS_FAMILY=unknown

  case "$KERNEL_NAME" in
  darwin*)
    OS=darwin
    OS_FAMILY=macos
    ;;
  linux*)
    OS=linux
    if [ -f /etc/os-release ]; then
      OS_FAMILY=$(grep </etc/os-release '^ID_LIKE' | cut -d= -f2)
      if [ -z "${OS_FAMILY}" ]; then
        OS_FAMILY=$(grep </etc/os-release '^ID' | cut -d= -f2)
      fi
    fi
    ;;
  cygwin* | mingw32* | msys* | mingw* | win*)
    OS=windows
    OS_FAMILY=windows
    ;;
  *)
    OS=$(uname -a 2>/dev/null || printf unknown)
    ;;
  esac
  OS_FAMILY=$(lowercase "$OS_FAMILY")

  export ENV_TYPE
  export OS
  export OS_FAMILY
}

checkPOSIX
detectPlatform
