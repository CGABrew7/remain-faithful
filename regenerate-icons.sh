#!/usr/bin/env bash
set -euo pipefail

SRC="RemainFaithful/Assets.xcassets/AppIcon.appiconset/AppIcon-1024.png"
DEST="RemainFaithful/Assets.xcassets/AppIcon.appiconset"

if [[ ! -f "$SRC" ]]; then
  echo "ERROR: master icon not found at $SRC" >&2
  exit 1
fi

resize() {
  local px=$1 name=$2
  sips -z "$px" "$px" "$SRC" --out "$DEST/$name" >/dev/null
  echo "  ${px}x${px}  ->  $name"
}

echo "Generating icons from $SRC …"

resize  40  "AppIcon-20@2x.png"
resize  60  "AppIcon-20@3x.png"
resize  58  "AppIcon-29@2x.png"
resize  87  "AppIcon-29@3x.png"
resize  76  "AppIcon-38@2x.png"
resize 114  "AppIcon-38@3x.png"
resize  80  "AppIcon-40@2x.png"
resize 120  "AppIcon-40@3x.png"
resize 120  "AppIcon-60@2x.png"
resize 180  "AppIcon-60@3x.png"
resize  76  "AppIcon-76@1x.png"
resize 152  "AppIcon-76@2x.png"
resize 167  "AppIcon-83.5@2x.png"

echo "Done — 13 sizes written."
