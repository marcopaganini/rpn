#!/bin/sh
# Generic installer program for binaries.
# See https://github.com/marcopaganini/installer for details.
# This script is bundled with the releases and should not be executed directly.

set -eu

readonly DEFAULT_INSTALL_DIR="/usr/local/bin"
readonly BIN="rpn"

main() {
  # Only argument is the install directory. If not provided,
  # the program will use DEFAULT_INSTALL_DIR.
  install_dir="${DEFAULT_INSTALL_DIR}"
  if [ $# -eq 1 ]; then
    install_dir="${1}"
  fi

  mkdir -p "${install_dir}"
  cp "${BIN}" "${install_dir}"
  chmod 755 "${install_dir}/${BIN}"
}

main "${@}"
