#!/bin/bash
# Stamp the Mac bundle id and copy the monitor helper into the app.
set -euo pipefail
cd "$(dirname "$0")/.."
APP=$(find build/bin -maxdepth 2 -name '*.app' -type d | head -1)
if [[ -z "${APP}" ]]; then
  echo "no .app was produced" >&2
  exit 1
fi
PLIST="${APP}/Contents/Info.plist"
/usr/libexec/PlistBuddy -c "Set :CFBundleIdentifier com.remainfaithful.app.mac" "$PLIST" \
  || /usr/libexec/PlistBuddy -c "Add :CFBundleIdentifier string com.remainfaithful.app.mac" "$PLIST"
/usr/libexec/PlistBuddy -c "Set :NSAppleEventsUsageDescription Remain Faithful closes apps you chose to block." "$PLIST" \
  || /usr/libexec/PlistBuddy -c "Add :NSAppleEventsUsageDescription string Remain Faithful closes apps you chose to block." "$PLIST"
swiftc -O -framework Cocoa -o "${APP}/Contents/MacOS/RFMonitorHelper" macos/RFMonitorHelper/main.swift
chmod +x "${APP}/Contents/MacOS/RFMonitorHelper"
ditto -c -k --keepParent "$APP" "build/bin/RemainFaithful-mac.zip"
echo "stamped ${APP}"
