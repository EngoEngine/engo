#!/usr/bin/env bash

set -o errexit

projectDir=$(cd "$(dirname "${0}")/.." && pwd)
# shellcheck source=script/util.sh
source "${projectDir}/script/util.sh" || source ./util.sh

ANDROID_NDK_VERSION="r21b"

usage() {
  println "POSIX-compliant bash script to manage toolchain for develop project"
  println "Usage: ${0} <option>"
  println "Options:"
  println "  -h this help"
  println "  -x enable debug mode (trace per command line in scripts)"
  println "  -c check requirements for environment"
  println "  -s setup environment ENV_TYPE=${ENV_TYPE}"
}

checkingHash() {
  filename="${1}"
  grep <"${projectDir}/script/checksum.txt" "${filename}" | sha256sum -c
}

installNodeJsLTS() {
  if ! [ -x "$(command -v npm)" ]; then
    info "getting and installing NodeJs"
    curl -fSL https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh -o get-nvm.sh
    bash get-nvm.sh
    rm get-nvm.sh
    # shellcheck disable=SC1090
    source "${HOME}/.nvm/nvm.sh"
    nvm install --lts
  fi
}

installAndroidNDK() {
  if [ "${ANDROID_NDK_HOME}" != "" ] && [ -d "${ANDROID_NDK_HOME}" ]; then
    info "already exist Android NDK"
    return 0
  fi

  if [ "${ANDROID_HOME}" != "" ] && [ -d "${ANDROID_HOME}/ndk-bundle" ]; then
    info "already exist Android NDK"
    ANDROID_NDK_HOME="${ANDROID_HOME}/ndk-bundle"
  else
    info "downloading Android NDK"
    filename="android-ndk-${ANDROID_NDK_VERSION}-${OS}-x86_64.zip"
    curl -fsSL "https://dl.google.com/android/repository/${filename}" -o "${filename}"
    checkingHash "${filename}"
    mkdir -p /usr/local/lib/android/ndk
    unzip "${filename}" -d /usr/local/lib/android/ndk
    rm "${filename}"
    ANDROID_NDK_HOME="/usr/local/lib/android/ndk/android-ndk-${ANDROID_NDK_VERSION}"
  fi
  echo "export ANDROID_NDK_HOME=${ANDROID_NDK_HOME}" | tee -a "${HOME}/.bashrc"
}

commonSetup() {
  # for compares benchmarking
  info "install tools for compares benchmarking"
  GO111MODULE=off go get golang.org/x/perf/cmd/...
  # for test coverage report
  info "install tools for test coverage report"
  GO111MODULE=off go get golang.org/x/tools/cmd/cover
  # for analyzes to identify unnecessary type conversions
  info "install tools for analyzes to identify unnecessary type conversions"
  GO111MODULE=off go get github.com/mdempsky/unconvert
}

debianSetup() {
  info "setup for platform debian"
  sudo apt-get update
  # base tools
  if ! [ -x "$(command -v make)" ]; then
    info "Install make"
    sudo apt-get install -qq -y --no-install-recommends make
  fi
  if ! [ -x "$(command -v curl)" ]; then
    info "Install curl"
    sudo apt-get install -qq -y --no-install-recommends curl
  fi
  # for testing
  if ! [ -x "$(command -v Xvfb)" ]; then
    info "Install xvfb"
    sudo apt-get install -qq -y --no-install-recommends xvfb
  fi
  if [ "${DISPLAY}" == "" ]; then
    info "Starting xvfb"
    Xvfb :1 &
    echo "export DISPLAY=:1" | tee -a "${HOME}/.bashrc"
  fi
  # for build project
  info "Installing libs for build project"
  sudo apt-get install -qq -y \
    libasound2-dev \
    libglu1-mesa-dev \
    freeglut3-dev \
    mesa-common-dev \
    xorg-dev \
    libgl1-mesa-dev
#  [ -x "$(command -v pkg-config)" ] || sudo apt-get install -qq -y --no-install-recommends pkg-config
#  [ "$(ldconfig -p | grep libgl)" != "" ] || sudo apt-get install -qq -y --no-install-recommends  mesa-common-dev libgl1-mesa-dev libglu1-mesa-dev
#  [ "$(ldconfig -p | grep libX11)" != "" ] || sudo apt-get install -qq -y --no-install-recommends xorg-dev
#  [ "$(ldconfig -p | grep libXcursor)" != "" ] || sudo apt-get install -qq -y --no-install-recommends libxcursor-dev
#  [ "$(ldconfig -p | grep libXrandr)" != "" ] || sudo apt-get install -qq -y --no-install-recommends libxrandr-dev
#  [ "$(ldconfig -p | grep libXinerama)" != "" ] || sudo apt-get install -qq -y --no-install-recommends libxinerama-dev
#  [ "$(ldconfig -p | grep 'libXi\.')" != "" ] || sudo apt-get install -qq -y --no-install-recommends libxi-dev
#  [ "$(ldconfig -p | grep libasound)" != "" ] || sudo apt-get install -qq -y --no-install-recommends libasound2-dev
#  [ "$(ldconfig -p | grep libglut)" != "" ] || sudo apt-get install -qq -y --no-install-recommends freeglut3-dev

  commonSetup
}

macosSetup() {
  info "setup for platform macos"
  if [ "${XCODE_11_DEVELOPER_DIR}" != "" ] && [ -d "${XCODE_11_DEVELOPER_DIR}" ]; then
    info "select Xcode dir {XCODE_11_DEVELOPER_DIR}"
    sudo Xcode-select --switch "${XCODE_11_DEVELOPER_DIR}"
  fi
  # TODO for X11
  commonSetup
}

windowsSetup() {
  info "setup for platform windows"
  # base tools
  tryCommand make || choco install -y make
  tryCommand curl || choco install -y curl
  commonSetup
}

desktopEnvironmentSetup() {
  echo "nothing more to add for environment desktop"
}

mobileEnvironmentSetup() {
  info "setup for environment mobile"
  if [ "${OS_FAMILY}" != "macos" ]; then
    installAndroidNDK
  fi
  # for build mobile
  info "install tools for build mobile"
  GO111MODULE=off go get golang.org/x/mobile/cmd/gomobile
  "$(go env GOPATH)/bin/gomobile" init
}

browserEnvironmentSetup() {
  info "setup for environment browser"
  # for build wasm
  info "install tools for build wasm"
  installNodeJsLTS
}

developerEnvironmentSetup() {
  desktopEnvironmentSetup
  browserEnvironmentSetup
  mobileEnvironmentSetup
}

checkRequirements() {
  printPlatform
  # TODO implement check requirements
}

checkEnvironment() {
  checkRequirements
  println "ENV_TYPE: ${ENV_TYPE}"
  tryCommand git version
  tryCommand go version
  tryCommand node || println "node $(node --version)"
  tryCommand npm || println "npm $(npm --version)"
  tryCommand make || println "make $(make --version | grep Make | cut -d" " -f3)"
  # TODO implement check environment
}

setupEnvironment() {
  checkRequirements

  case "$OS_FAMILY" in
  debian) debianSetup ;;
  macos) macosSetup ;;
  windows) windowsSetup ;;
  *) notReady "setup toolchain" ;;
  esac

  case "$ENV_TYPE" in
  browser) browserEnvironmentSetup ;;
  mobile) mobileEnvironmentSetup ;;
  developer) developerEnvironmentSetup ;;
  esac
}

main() {
  if [ "$(id -u)" == "0" ]; then fatal "Not running as root"; fi
  if [ -z "$*" ]; then usage; fi

  cmd=
  while getopts ":hxsc" flag; do
    case "${flag}" in
    x) set -o xtrace ;;
    s) cmd=setupEnvironment ;;
    c) cmd=checkEnvironment ;;
    ?) usage ;;
    esac
  done

  ${cmd}
}

main "$*"
