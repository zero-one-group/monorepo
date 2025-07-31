#!/usr/bin/env bash
set -euo pipefail

# ------------------------------------------------------------------
# usage / args
# ------------------------------------------------------------------
if [ $# -lt 1 ]; then
  echo "Usage: $0 <version> [install_dir]"
  echo "Example: $0 v3.5.1 /usr/local/bin"
  exit 1
fi

VERSION_RAW=$1
INSTALL_DIR=${2:-/usr/local/bin}
VERSION=${VERSION_RAW#v}   # strip leading "v"

# ------------------------------------------------------------------
# dependency checks
# ------------------------------------------------------------------
for cmd in curl mktemp find; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Error: '$cmd' is required but not installed." >&2
    exit 1
  fi
done

# tar is needed for .tar.gz, unzip for .zip
HAS_TAR=true
command -v tar >/dev/null 2>&1 || HAS_TAR=false
HAS_UNZIP=true
command -v unzip >/dev/null 2>&1 || HAS_UNZIP=false

# ------------------------------------------------------------------
# detect OS & ARCH, set archive EXT
# ------------------------------------------------------------------
case "$(uname -s)" in
  Linux)
    OS="Linux"
    EXT="tar.gz"
    ;;
  Darwin)
    OS="Darwin"
    EXT="tar.gz"
    ;;
  MINGW*|MSYS*|CYGWIN*)
    OS="Windows"
    EXT="zip"
    ;;
  *)
    echo "Unsupported OS: $(uname -s)" >&2
    exit 1
    ;;
esac

case "$(uname -m)" in
  x86_64|amd64)  ARCH="x86_64" ;;
  arm64|aarch64) ARCH="arm64"  ;;
  *)
    echo "Unsupported ARCH: $(uname -m)" >&2
    exit 1
    ;;
esac

# ------------------------------------------------------------------
# decide binary name, check for existing install
# ------------------------------------------------------------------
BIN="mockery"
[ "$OS" = "Windows" ] && BIN+=".exe"

if [ -x "$INSTALL_DIR/$BIN" ]; then
  echo "✅  '$BIN' already exists in $INSTALL_DIR; nothing to do."
  exit 0
fi

# ------------------------------------------------------------------
# download & unpack
# ------------------------------------------------------------------
FILE="mockery_${VERSION}_${OS}_${ARCH}.${EXT}"
URL="https://github.com/vektra/mockery/releases/download/${VERSION_RAW}/${FILE}"

WORKDIR=$(mktemp -d)
trap 'rm -rf "$WORKDIR"' EXIT

ARCHIVE="$WORKDIR/$FILE"
echo "Downloading $URL …"
curl -fsSL -o "$ARCHIVE" "$URL"

echo "Unpacking $FILE …"
if [ "$EXT" = "zip" ]; then
  if ! $HAS_UNZIP; then
    echo "Error: unzip required to extract .zip files" >&2
    exit 1
  fi
  unzip -q "$ARCHIVE" -d "$WORKDIR"
else
  if ! $HAS_TAR; then
    echo "Error: tar required to extract .tar.gz files" >&2
    exit 1
  fi
  tar -xzf "$ARCHIVE" -C "$WORKDIR"
fi

# ------------------------------------------------------------------
# locate & install
# ------------------------------------------------------------------
BPATH=$(find "$WORKDIR" -type f -name "$BIN" -print -quit)
if [ -z "$BPATH" ]; then
  echo "Error: $BIN not found in the archive" >&2
  exit 1
fi

echo "Installing $BIN to $INSTALL_DIR …"
# if 'install' isn’t available, fall back to cp+chmod
if command -v install >/dev/null; then
  sudo install -m 0755 "$BPATH" "$INSTALL_DIR/$BIN"
else
  cp "$BPATH" "$INSTALL_DIR/$BIN"
  chmod 0755 "$INSTALL_DIR/$BIN"
fi

echo "✅  Installed mockery $VERSION to $INSTALL_DIR/$BIN"
