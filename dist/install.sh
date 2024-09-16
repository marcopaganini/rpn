#!/bin/sh
# Generic installer program for binaries.
# See https://github.com/marcopaganini/installer for details.
# This script is bundled with the releases and should not be executed directly.

set -eu

readonly PREFIX="/usr/local"
readonly BIN="rpn"

main() {
  uid="$(id -u)"
  if [ "${uid}" -ne 0 ]; then
    die >&2 "Please run this program as root (using sudo)."
    exit 1
  fi

  bindir="${PREFIX}/bin"

  mkdir -p "${bindir}"
  cp "${BIN}" "${bindir}"
  chmod 755 "${bindir}/${BIN}"
}

main "${@}"
