#!/bin/bash
# clang wrapper so Wails can compile the Intel slice on an Apple Silicon runner.
# Go invokes this as CC. When Go has already chosen an architecture, leave it alone.
set -e
if [[ " $* " == *" -arch "* || " $* " == *" -target "* ]]; then
  exec clang "$@"
fi
case "${GOARCH:-}" in
  amd64) exec clang -arch x86_64 "$@" ;;
  arm64) exec clang -arch arm64 "$@" ;;
  *) exec clang "$@" ;;
esac
