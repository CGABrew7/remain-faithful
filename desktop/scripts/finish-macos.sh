#!/bin/bash
# Sign and notarize the Mac app when GitHub Actions secrets are present,
# then always write build/bin/RemainFaithful-macos.dmg.
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

DMG="build/bin/RemainFaithful-macos.dmg"
STATUS="build/bin/signing-status.txt"
ENTITLEMENTS="macos/entitlements.plist"
HELPER="${APP}/Contents/MacOS/RFMonitorHelper"

if [[ ! -f "$ENTITLEMENTS" ]]; then
  echo "entitlements file is missing" >&2
  exit 1
fi
if [[ ! -f "$HELPER" ]]; then
  echo "RFMonitorHelper is missing from the app" >&2
  exit 1
fi

note() {
  echo "$1"
  if [[ -n "${GITHUB_STEP_SUMMARY:-}" ]]; then
    echo "$1" >> "$GITHUB_STEP_SUMMARY"
  fi
}

package_dmg() {
  local stage
  stage=$(mktemp -d)
  ditto "$APP" "$stage/$(basename "$APP")"
  ln -s /Applications "$stage/Applications"
  rm -f "$DMG"
  hdiutil create \
    -volname "Remain Faithful" \
    -srcfolder "$stage" \
    -ov \
    -format UDZO \
    "$DMG"
  rm -rf "$stage"
  echo "packaged ${DMG}"
}

required_secrets=(
  MACOS_CERT_P12_BASE64
  MACOS_CERT_PASSWORD
  APPLE_API_KEY_ID
  APPLE_API_ISSUER_ID
  APPLE_API_KEY_P8_BASE64
  APPLE_TEAM_ID
)
missing=()
for name in "${required_secrets[@]}"; do
  if [[ -z "${!name:-}" ]]; then
    missing+=("$name")
  fi
done

if [[ ${#missing[@]} -gt 0 ]]; then
  note "Mac signing skipped. Missing secrets: ${missing[*]}"
  note "The disk image is unsigned and this job stays green."
  package_dmg
  printf 'mac=unsigned\n' > "$STATUS"
  exit 0
fi

P12=""
P8=""
KEYCHAIN=""
cleanup() {
  if [[ -n "$KEYCHAIN" ]]; then
    security delete-keychain "$KEYCHAIN" >/dev/null 2>&1 || true
  fi
  if [[ -n "$P12" ]]; then
    rm -f "$P12"
  fi
  if [[ -n "$P8" ]]; then
    rm -f "$P8"
  fi
}
trap cleanup EXIT

decode_secret() {
  local var_name="$1"
  local dest="$2"
  SECRET_NAME="$var_name" python3 - "$dest" <<'PY'
import os, pathlib, sys
name = os.environ["SECRET_NAME"]
raw = "".join(os.environ.get(name, "").split())
pathlib.Path(sys.argv[1]).write_bytes(__import__("base64").b64decode(raw))
PY
  chmod 600 "$dest"
}

P12="${RUNNER_TEMP:-/tmp}/remain-faithful-dev-id.p12"
P8="${RUNNER_TEMP:-/tmp}/AuthKey.p8"
decode_secret MACOS_CERT_P12_BASE64 "$P12"
decode_secret APPLE_API_KEY_P8_BASE64 "$P8"
if ! grep -q "PRIVATE KEY" "$P8"; then
  echo "APPLE_API_KEY_P8_BASE64 did not decode to an App Store Connect API key." >&2
  exit 1
fi

KEYCHAIN="${RUNNER_TEMP:-/tmp}/remain-faithful-signing.keychain-db"
KEYCHAIN_PASSWORD=$(openssl rand -base64 32)
security create-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN"
security set-keychain-settings -lut 21600 "$KEYCHAIN"
security unlock-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN"
security import "$P12" -k "$KEYCHAIN" -P "$MACOS_CERT_PASSWORD" -A -t cert -f pkcs12 -T /usr/bin/codesign
security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "$KEYCHAIN_PASSWORD" "$KEYCHAIN" >/dev/null
security list-keychains -d user -s "$KEYCHAIN" $(security list-keychains -d user | tr -d '"')
security default-keychain -s "$KEYCHAIN"

identities=$(security find-identity -v -p codesigning "$KEYCHAIN")
IDENTITY=$(printf '%s\n' "$identities" | sed -n "s/.*\"\(Developer ID Application:.*(${APPLE_TEAM_ID})\)\"/\1/p" | sed -n '1p')
if [[ -z "$IDENTITY" ]]; then
  echo "No Developer ID Application certificate matched APPLE_TEAM_ID." >&2
  exit 1
fi
echo "Developer ID Application certificate loaded."

sign_macho() {
  local bin="$1"
  if [[ "$(basename "$bin")" == "RFMonitorHelper" ]]; then
    codesign --force --options runtime --timestamp --sign "$IDENTITY" "$bin"
  else
    codesign --force --options runtime --timestamp \
      --entitlements "$ENTITLEMENTS" \
      --sign "$IDENTITY" \
      "$bin"
  fi
}

xattr -cr "$APP"
while IFS= read -r bin; do
  kind=$(file -b "$bin" || true)
  case "$kind" in
    Mach-O*) sign_macho "$bin" ;;
  esac
done < <(find "$APP" -type f)
codesign --force --options runtime --timestamp \
  --entitlements "$ENTITLEMENTS" \
  --sign "$IDENTITY" \
  "$APP"
codesign --verify --deep --strict "$APP"
codesign --verify --strict "$HELPER"

details=$(codesign -d --verbose=4 "$APP" 2>&1 || true)
case "$details" in
  *runtime*) ;;
  *)
    echo "Hardened runtime is missing from the app signature." >&2
    exit 1
    ;;
esac
entitlement_dump=$(codesign -d --entitlements :- "$APP" 2>/dev/null || true)
case "$entitlement_dump" in
  *com.apple.security.network.client*) ;;
  *)
    echo "The app signature is missing the existing entitlements." >&2
    exit 1
    ;;
esac
unset details entitlement_dump

notarize() {
  local file="$1"
  local result err id status
  result=$(mktemp)
  err=$(mktemp)
  set +e
  xcrun notarytool submit "$file" \
    --key "$P8" \
    --key-id "$APPLE_API_KEY_ID" \
    --issuer "$APPLE_API_ISSUER_ID" \
    --wait \
    --output-format json >"$result" 2>"$err"
  local code=$?
  set -e
  if [[ $code -ne 0 ]]; then
    echo "Notarization failed for $(basename "$file")." >&2
    if [[ -s "$err" ]]; then
      cat "$err" >&2
    fi
    if [[ -s "$result" ]]; then
      id=$(python3 -c 'import json,sys
try:
    print(json.load(open(sys.argv[1])).get("id",""))
except Exception:
    print("")' "$result")
      if [[ -n "$id" ]]; then
        xcrun notarytool log "$id" \
          --key "$P8" \
          --key-id "$APPLE_API_KEY_ID" \
          --issuer "$APPLE_API_ISSUER_ID" >&2 || true
      fi
    fi
    rm -f "$result" "$err"
    exit 1
  fi
  status=$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1])).get("status",""))' "$result")
  rm -f "$result" "$err"
  if [[ "$status" != "Accepted" ]]; then
    echo "Notarization status was not Accepted." >&2
    exit 1
  fi
  echo "Notarization accepted for $(basename "$file")."
}

staple() {
  local target="$1"
  local attempt=1
  while [[ $attempt -le 5 ]]; do
    if xcrun stapler staple "$target"; then
      xcrun stapler validate "$target"
      return 0
    fi
    echo "Staple attempt ${attempt} did not finish. Waiting before trying again."
    sleep 20
    attempt=$((attempt + 1))
  done
  echo "Stapling failed for $(basename "$target")." >&2
  return 1
}

notary_dir=$(mktemp -d)
notary_zip="$notary_dir/RemainFaithful-macos.zip"
ditto -c -k --keepParent "$APP" "$notary_zip"
notarize "$notary_zip"
rm -rf "$notary_dir"
staple "$APP"

package_dmg
codesign --force --timestamp --sign "$IDENTITY" "$DMG"
codesign --verify --strict "$DMG"
notarize "$DMG"
staple "$DMG"

printf 'mac=signed\n' > "$STATUS"
note "Mac disk image: signed."
