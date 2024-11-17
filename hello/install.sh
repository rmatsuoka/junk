#!/bin/sh

set -ux

PREFIX="${PREFIX-/usr/local}"

do_install() {
  src="$1"
  distdir="$2"

  mkdir -p "$distdir" || exit
  cp "$src" "$distdir" || exit
}

do_install hello "${PREFIX}/bin"
