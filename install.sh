#!/bin/bash

# Ultimate SDD Framework Installer
# Automatically detects platform and installs the appropriate binary

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Repository information
REPO="iamcoolvivek007/Ultimate-SDD-Framework"
GITHUB_API="https://api.github.com/repos/$REPO/releases/latest"

# Print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Detect platform
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $OS in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            print_error "Unsupported OS: $OS"
            exit 1
            ;;
    esac

    case $ARCH in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    BINARY_NAME="sdd-$OS-$ARCH"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="$BINARY_NAME.exe"
        ARCHIVE_EXT="zip"
    else
        ARCHIVE_EXT="tar.gz"
    fi

    ARCHIVE_NAME="$BINARY_NAME.$ARCHIVE_EXT"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Download file with progress
download() {
    local url="$1"
    local output="$2"

    if command_exists curl; then
        curl -L -o "$output" "$url" --progress-bar
    elif command_exists wget; then
        wget -O "$output" "$url"
    else
        print_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

# Get latest release information
get_latest_release() {
    print_info "Fetching latest release information..."

    if command_exists jq; then
        # Use jq if available
        LATEST_VERSION=$(curl -s "$GITHUB_API" | jq -r '.tag_name')
        DOWNLOAD_URL=$(curl -s "$GITHUB_API" | jq -r ".assets[] | select(.name == \"$ARCHIVE_NAME\") | .browser_download_url")
    else
        # Fallback to grep/sed
        LATEST_VERSION=$(curl -s "$GITHUB_API" | grep '"tag_name"' | head -1 | cut -d'"' -f4)
        DOWNLOAD_URL=$(curl -s "$GITHUB_API" | grep "browser_download_url.*$ARCHIVE_NAME" | cut -d'"' -f4)
    fi

    if [ -z "$DOWNLOAD_URL" ]; then
        print_error "Could not find download URL for $ARCHIVE_NAME"
        print_info "Available assets:"
        curl -s "$GITHUB_API" | grep "name.*sdd-" | head -10
        exit 1
    fi
}

# Install binary
install_binary() {
    local temp_dir=$(mktemp -d)
    local archive_path="$temp_dir/$ARCHIVE_NAME"

    print_info "Downloading Ultimate SDD Framework $LATEST_VERSION for $OS/$ARCH..."

    # Download archive
    download "$DOWNLOAD_URL" "$archive_path"

    # Extract archive
    print_info "Extracting..."
    cd "$temp_dir"
    if [ "$ARCHIVE_EXT" = "tar.gz" ]; then
        tar -xzf "$ARCHIVE_NAME"
    else
        unzip "$ARCHIVE_NAME"
    fi

    # Find the binary
    local binary_path=""
    if [ -f "sdd" ]; then
        binary_path="sdd"
    elif [ -f "sdd.exe" ]; then
        binary_path="sdd.exe"
    else
        print_error "Could not find binary in archive"
        ls -la
        exit 1
    fi

    # Make binary executable
    chmod +x "$binary_path"

    # Determine install location
    local install_dir=""
    if [ -w "/usr/local/bin" ] || [ -w "/usr/local" ]; then
        install_dir="/usr/local/bin"
    elif [ -w "/usr/bin" ]; then
        install_dir="/usr/bin"
    else
        install_dir="$HOME/.local/bin"
        mkdir -p "$install_dir"
    fi

    # Install binary
    print_info "Installing to $install_dir..."
    mv "$binary_path" "$install_dir/"

    # Clean up
    cd /
    rm -rf "$temp_dir"

    print_success "Ultimate SDD Framework $LATEST_VERSION installed successfully!"
    print_info "Binary location: $install_dir/sdd"
    print_info "Run 'sdd --help' to get started"
}

# Verify installation
verify_installation() {
    if command_exists sdd; then
        print_success "Installation verified!"
        sdd --version 2>/dev/null || sdd --help | head -3
    else
        print_warning "Binary not found in PATH. You may need to:"
        echo "  export PATH=\"$install_dir:\$PATH\""
        echo "  # Or add $install_dir to your shell profile"
    fi
}

# Main installation process
main() {
    print_info "Ultimate SDD Framework Installer"
    print_info "================================="

    # Check for required tools
    if ! command_exists curl; then
        print_error "curl is required but not installed. Please install curl and try again."
        exit 1
    fi

    # Detect platform
    detect_platform
    print_info "Detected platform: $OS/$ARCH"

    # Get latest release
    get_latest_release
    print_info "Latest version: $LATEST_VERSION"

    # Install
    install_binary

    # Verify
    verify_installation

    print_success "Installation complete! 🎉"
    echo ""
    print_info "Next steps:"
    echo "  1. Run: sdd mcp add openai-main --provider openai --model gpt-4"
    echo "  2. Run: sdd init \"Your Project\""
    echo "  3. Start developing with: sdd specify \"your feature\""
    echo ""
    print_info "Documentation: https://github.com/$REPO#readme"
}

# Run main function
main "$@"