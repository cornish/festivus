#!/bin/sh
# Textivus installer
#
# User install (default):
#   curl -fsSL https://raw.githubusercontent.com/cornish/textivus-editor/main/install.sh | sh
#
# System-wide install:
#   curl -fsSL https://raw.githubusercontent.com/cornish/textivus-editor/main/install.sh | sh -s -- --bin-dir /usr/local/bin
#
# Options:
#   --version VERSION    Install specific version (default: latest)
#   --bin-dir DIR        Installation directory (default: ~/.local/bin)

set -e

REPO="cornish/textivus-editor"
BINARY_NAME="textivus"
DEFAULT_BIN_DIR="$HOME/.local/bin"

# Colors (disabled if not a terminal)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
fi

info() {
    printf "${BLUE}==>${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}Warning:${NC} %s\n" "$1"
}

error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux) OS="linux" ;;
        darwin) OS="darwin" ;;
        *) error "Unsupported operating system: $OS" ;;
    esac

    case "$ARCH" in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac

    PLATFORM="${OS}-${ARCH}"
}

# Get latest version from GitHub API
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Download file
download() {
    url="$1"
    dest="$2"

    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$url" -o "$dest"
    elif command -v wget >/dev/null 2>&1; then
        wget -q "$url" -O "$dest"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Verify checksum
verify_checksum() {
    binary_path="$1"
    checksum_path="$2"

    if command -v sha256sum >/dev/null 2>&1; then
        expected=$(cat "$checksum_path" | awk '{print $1}')
        actual=$(sha256sum "$binary_path" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        expected=$(cat "$checksum_path" | awk '{print $1}')
        actual=$(shasum -a 256 "$binary_path" | awk '{print $1}')
    else
        warn "sha256sum/shasum not found, skipping checksum verification"
        return 0
    fi

    if [ "$expected" != "$actual" ]; then
        error "Checksum verification failed!\nExpected: $expected\nActual: $actual"
    fi
}

# Parse arguments
BIN_DIR="$DEFAULT_BIN_DIR"
VERSION=""

while [ $# -gt 0 ]; do
    case "$1" in
        --version|-v)
            VERSION="$2"
            shift 2
            ;;
        --bin-dir|-d)
            BIN_DIR="$2"
            shift 2
            ;;
        --help|-h)
            echo "Textivus installer"
            echo ""
            echo "Usage: install.sh [options]"
            echo ""
            echo "Options:"
            echo "  --version, -v VERSION  Install specific version (default: latest)"
            echo "  --bin-dir, -d DIR      Installation directory (default: ~/.local/bin)"
            echo "  --help, -h             Show this help message"
            exit 0
            ;;
        *)
            error "Unknown option: $1"
            ;;
    esac
done

# Main installation
main() {
    info "Installing Textivus..."

    detect_platform
    info "Detected platform: $PLATFORM"

    # Get version
    if [ -z "$VERSION" ]; then
        info "Fetching latest version..."
        VERSION=$(get_latest_version)
        if [ -z "$VERSION" ]; then
            error "Could not determine latest version. Try specifying with --version"
        fi
    fi
    info "Version: $VERSION"

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download binary
    BINARY_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${PLATFORM}"
    CHECKSUM_URL="${BINARY_URL}.sha256"

    info "Downloading binary..."
    download "$BINARY_URL" "$TMP_DIR/$BINARY_NAME"

    info "Downloading checksum..."
    download "$CHECKSUM_URL" "$TMP_DIR/$BINARY_NAME.sha256"

    info "Verifying checksum..."
    verify_checksum "$TMP_DIR/$BINARY_NAME" "$TMP_DIR/$BINARY_NAME.sha256"
    success "Checksum verified"

    # Install binary
    mkdir -p "$BIN_DIR"
    mv "$TMP_DIR/$BINARY_NAME" "$BIN_DIR/$BINARY_NAME"
    chmod +x "$BIN_DIR/$BINARY_NAME"

    # Create txv shortcut symlink
    ln -sf "$BIN_DIR/$BINARY_NAME" "$BIN_DIR/txv"

    success "Installed to $BIN_DIR/$BINARY_NAME"
    success "Created shortcut: txv -> textivus"

    # Check if BIN_DIR is in PATH
    case ":$PATH:" in
        *":$BIN_DIR:"*) ;;
        *)
            echo ""
            warn "$BIN_DIR is not in your PATH"
            echo ""
            echo "Add it to your shell configuration:"
            echo ""
            echo "  # For bash (~/.bashrc)"
            echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
            echo ""
            echo "  # For zsh (~/.zshrc)"
            echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
            echo ""
            echo "  # For fish (~/.config/fish/config.fish)"
            echo "  set -gx PATH \$HOME/.local/bin \$PATH"
            echo ""
            ;;
    esac

    echo ""
    success "Textivus $VERSION installed successfully!"
    echo ""
    echo "Run 'textivus --help' to get started."
}

main
