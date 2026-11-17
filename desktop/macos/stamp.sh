#!/bin/bash
# Stamp the Mac bundle id and copy a universal monitor helper into the app.
set -euo pipefail
cd "$(dirname "$0")/.."

APP=""
while IFS= read -r candidate; do
  if [[ -z "$APP" ]]; then
    APP=$candidate
  fi
done < <(find build/bin -maxdepth 2 -name '*.app' -type d)
if [[ -z "${APP}" ]]; then
  echo "no .app was produced" >&2
  exit 1
fi

PLIST="${APP}/Contents/Info.plist"
/usr/libexec/PlistBuddy -c "Set :CFBundleIdentifier com.remainfaithful.app.mac" "$PLIST" \
  || /usr/libexec/PlistBuddy -c "Add :CFBundleIdentifier string com.remainfaithful.app.mac" "$PLIST"
/usr/libexec/PlistBuddy -c "Set :NSAppleEventsUsageDescription Remain Faithful closes apps you chose to block." "$PLIST" \
  || /usr/libexec/PlistBuddy -c "Add :NSAppleEventsUsageDescription string Remain Faithful closes apps you chose to block." "$PLIST"

EXE_NAME=$(/usr/libexec/PlistBuddy -c "Print :CFBundleExecutable" "$PLIST")
MAIN_BIN="${APP}/Contents/MacOS/${EXE_NAME}"
if [[ ! -f "$MAIN_BIN" ]]; then
  echo "app executable is missing: ${EXE_NAME}" >&2
  exit 1
fi

SDK=$(xcrun --sdk macosx --show-sdk-path)
workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT
swiftc -O -sdk "$SDK" -target arm64-apple-macos11.0 -framework Cocoa \
  -o "$workdir/helper-arm64" macos/RFMonitorHelper/main.swift
swiftc -O -sdk "$SDK" -target x86_64-apple-macos10.13 -framework Cocoa \
  -o "$workdir/helper-amd64" macos/RFMonitorHelper/main.swift
mkdir -p "${APP}/Contents/MacOS"
lipo -create -output "${APP}/Contents/MacOS/RFMonitorHelper" \
  "$workdir/helper-arm64" "$workdir/helper-amd64"
chmod +x "${APP}/Contents/MacOS/RFMonitorHelper"
rm -rf "$workdir"
trap - EXIT

require_universal() {
  local label="$1"
  local bin="$2"
  local info
  info=$(lipo -info "$bin")
  case "$info" in
    *x86_64*arm64*|*arm64*x86_64*) ;;
    *)
      echo "${label} is not a universal binary (arm64 + x86_64)." >&2
      exit 1
      ;;
  esac
}

require_universal "RFMonitorHelper" "${APP}/Contents/MacOS/RFMonitorHelper"
if [[ "${DESKTOP_REQUIRE_UNIVERSAL:-}" == "1" ]]; then
  require_universal "Remain Faithful" "$MAIN_BIN"
fi

echo "stamped ${APP}"
