#!/bin/sh
set -e

BINARY="coding-type"
INSTALL_DIR="${CODING_TYPE_INSTALL_DIR:-$HOME/.local/bin}"
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/coding-type"

echo "Uninstalling ${BINARY}..."

# Remove binary
BINARY_PATH="${INSTALL_DIR}/${BINARY}"
if [ -f "$BINARY_PATH" ]; then
  rm -f "$BINARY_PATH"
  echo "Removed binary: ${BINARY_PATH}"
else
  echo "Binary not found at ${BINARY_PATH} (already removed?)"
fi

# Ask before removing history
if [ -d "$CONFIG_DIR" ]; then
  printf "Remove history and config at %s? [y/N] " "$CONFIG_DIR"
  read -r answer
  case "$answer" in
    [Yy]*)
      rm -rf "$CONFIG_DIR"
      echo "Removed config: ${CONFIG_DIR}"
      ;;
    *)
      echo "Config kept at: ${CONFIG_DIR}"
      ;;
  esac
fi

echo ""
echo "Uninstalled successfully."
