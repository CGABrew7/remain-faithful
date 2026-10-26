#!/bin/bash
# Build checksums and publish the desktop GitHub Release for a desktop-v* tag.
set -euo pipefail

if [[ $# -ne 5 ]]; then
  echo "usage: publish-release.sh <mac-dmg> <mac-status> <win-setup> <win-status> <out-dir>" >&2
  exit 1
fi

mac_dmg=$1
mac_status=$2
win_setup=$3
win_status=$4
out=$5

for input in "$mac_dmg" "$mac_status" "$win_setup" "$win_status"; do
  if [[ ! -f "$input" ]]; then
    echo "Missing release input: $input" >&2
    exit 1
  fi
done

# GitHub does not let a step replace GITHUB_REF_NAME, so tests pass DESKTOP_RELEASE_TAG.
tag=${DESKTOP_RELEASE_TAG:-${GITHUB_REF_NAME:-}}
version=${tag#desktop-v}
if [[ "$tag" == "$version" || ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+([-.][0-9A-Za-z.]+)?$ ]]; then
  echo "Refusing to publish. The tag must look like desktop-v1.0.0." >&2
  exit 1
fi

mkdir -p "$out"
cp "$mac_dmg" "$out/RemainFaithful-macos.dmg"
cp "$win_setup" "$out/RemainFaithful-windows-setup.exe"
(
  cd "$out"
  sha256sum RemainFaithful-macos.dmg RemainFaithful-windows-setup.exe > SHA256SUMS.txt
  sha256sum -c SHA256SUMS.txt
)

mac_state=$(tr -d '\r' < "$mac_status")
win_state=$(tr -d '\r' < "$win_status")
mac_line="RemainFaithful-macos.dmg (unsigned, signing secrets were not set)"
if [[ "$mac_state" == "mac=signed" ]]; then
  mac_line="RemainFaithful-macos.dmg (signed)"
fi
win_line="RemainFaithful-windows-setup.exe (unsigned, signing secrets were not set)"
case "$win_state" in
  windows=azure)
    win_line="RemainFaithful-windows-setup.exe (signed with Azure Artifact Signing)"
    ;;
  windows=pfx)
    win_line="RemainFaithful-windows-setup.exe (signed with a PFX certificate)"
    ;;
esac

title="Remain Faithful Desktop ${version}"
cat > "$out/NOTES.txt" <<EOF
Remain Faithful desktop ${version}.

Mac: ${mac_line}
Windows: ${win_line}

SHA256SUMS.txt lists the checksum for each file.
EOF

if [[ "${PUBLISH_DRY_RUN:-}" == "1" ]]; then
  echo "Dry run wrote ${out}"
  exit 0
fi

if [[ -z "${GH_TOKEN:-}${GITHUB_TOKEN:-}" ]]; then
  echo "GH_TOKEN is not set." >&2
  exit 1
fi

assets=(
  "$out/RemainFaithful-macos.dmg"
  "$out/RemainFaithful-windows-setup.exe"
  "$out/SHA256SUMS.txt"
)
if gh release view "$tag" >/dev/null 2>&1; then
  if [[ "$version" == *-* ]]; then
    gh release edit "$tag" --prerelease --title "$title" --notes-file "$out/NOTES.txt"
  else
    gh release edit "$tag" --title "$title" --notes-file "$out/NOTES.txt"
  fi
  gh release upload "$tag" --clobber "${assets[@]}"
else
  if [[ "$version" == *-* ]]; then
    gh release create "$tag" --prerelease --title "$title" --notes-file "$out/NOTES.txt" "${assets[@]}"
  else
    gh release create "$tag" --title "$title" --notes-file "$out/NOTES.txt" "${assets[@]}"
  fi
fi
