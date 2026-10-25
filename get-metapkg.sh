#!/usr/bin/env bash
# Download the latest metapkg release binary for this machine.
#
# Usage: ./get-metapkg.sh [install-dir]      (default: current directory)
# Pin a version with: METAPKG_VERSION=v0.1.0 ./get-metapkg.sh
set -euo pipefail

REPO="kenliu/metapkg"
DEST_DIR="${1:-.}"
VERSION="${METAPKG_VERSION:-latest}"

case "$(uname -s)" in
  Linux)  os=linux ;;
  Darwin) os=darwin ;;
  *) echo "Unsupported OS: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64)  arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo "Unsupported architecture: $(uname -m)" >&2; exit 1 ;;
esac

asset="metapkg_${os}_${arch}"
if [[ "${VERSION}" == "latest" ]]; then
  base="https://github.com/${REPO}/releases/latest/download"
else
  base="https://github.com/${REPO}/releases/download/${VERSION}"
fi

tmp="$(mktemp -d)"
trap 'rm -rf "${tmp}"' EXIT

echo "Downloading ${asset} (${VERSION})..."
curl -fsSL -o "${tmp}/${asset}" "${base}/${asset}"
curl -fsSL -o "${tmp}/SHA256SUMS" "${base}/SHA256SUMS"

# Verify against the release's SHA256SUMS (sha256sum on Linux, shasum on macOS).
expected="$(awk -v f="${asset}" '$2 == f {print $1}' "${tmp}/SHA256SUMS")"
if [[ -z "${expected}" ]]; then
  echo "No checksum for ${asset} in SHA256SUMS" >&2
  exit 1
fi
if command -v sha256sum >/dev/null 2>&1; then
  actual="$(sha256sum "${tmp}/${asset}" | awk '{print $1}')"
else
  actual="$(shasum -a 256 "${tmp}/${asset}" | awk '{print $1}')"
fi
if [[ "${actual}" != "${expected}" ]]; then
  echo "Checksum mismatch for ${asset}" >&2
  exit 1
fi

mkdir -p "${DEST_DIR}"
install -m 0755 "${tmp}/${asset}" "${DEST_DIR}/metapkg"

echo "Installed ${DEST_DIR}/metapkg"
"${DEST_DIR}/metapkg" version || true
