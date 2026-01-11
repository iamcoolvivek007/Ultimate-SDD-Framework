#!/bin/bash
# Viki Install Script
# https://github.com/iamcoolvivek007/Ultimate-SDD-Framework
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/iamcoolvivek007/Ultimate-SDD-Framework/main/scripts/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

# Config
REPO="iamcoolvivek007/Ultimate-SDD-Framework"
BINARY_NAME="sdd"
INSTALL_DIR="/usr/local/bin"

echo -e "${PURPLE}"
echo "╔══════════════════════════════════════════╗"
echo "║     🤖 Viki - AI Development Assistant   ║"
echo "║         Installation Script              ║"
echo "╚══════════════════════════════════════════╝"
echo -e "${NC}"

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) 
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac

    case "$OS" in
        linux) OS="linux" ;;
        darwin) OS="darwin" ;;
        msys*|mingw*|cygwin*) OS="windows" ;;
        *)
            echo -e "${RED}Unsupported OS: $OS${NC}"
            exit 1
            ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    echo -e "${BLUE}Detected platform: ${PLATFORM}${NC}"
}

# Get latest release version
get_latest_version() {
    echo -e "${BLUE}Fetching latest version...${NC}"
    
    VERSION=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" | \
        grep '"tag_name"' | \
        sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$VERSION" ]; then
        echo -e "${YELLOW}Could not fetch latest version, using 'main' branch${NC}"
        VERSION="main"
    fi
    
    echo -e "${GREEN}Latest version: ${VERSION}${NC}"
}

# Download and install
install_viki() {
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/ultimate-sdd-framework-${VERSION}-sdd-${VERSION}-${PLATFORM}.tar.gz"
    TEMP_DIR=$(mktemp -d)
    
    echo -e "${BLUE}Downloading Viki...${NC}"
    
    # Use -f to fail on HTTP errors (404, etc.) instead of downloading HTML error pages
    if ! curl -fSL "$DOWNLOAD_URL" -o "${TEMP_DIR}/viki.tar.gz" 2>/dev/null; then
        # Try alternative naming convention
        DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/sdd-${VERSION}-${PLATFORM}.tar.gz"
        if ! curl -fSL "$DOWNLOAD_URL" -o "${TEMP_DIR}/viki.tar.gz" 2>/dev/null; then
            echo -e "${YELLOW}Release binary not found, building from source...${NC}"
            rm -rf "$TEMP_DIR"
            install_from_source
            return
        fi
    fi
    
    echo -e "${BLUE}Extracting...${NC}"
    tar -xzf "${TEMP_DIR}/viki.tar.gz" -C "${TEMP_DIR}"
    
    # Find the binary
    BINARY=$(find "${TEMP_DIR}" -name "sdd*" -type f | head -1)
    
    if [ -z "$BINARY" ]; then
        echo -e "${RED}Binary not found in archive${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}Installing to ${INSTALL_DIR}...${NC}"
    
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY" "${INSTALL_DIR}/viki"
        chmod +x "${INSTALL_DIR}/viki"
    else
        sudo mv "$BINARY" "${INSTALL_DIR}/viki"
        sudo chmod +x "${INSTALL_DIR}/viki"
    fi
    
    rm -rf "$TEMP_DIR"
}

# Install from source (fallback)
install_from_source() {
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Go is not installed. Please install Go 1.21+ first.${NC}"
        echo -e "${BLUE}Visit: https://go.dev/dl/${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}Installing from source...${NC}"
    
    TEMP_DIR=$(mktemp -d)
    git clone --depth 1 "https://github.com/${REPO}.git" "$TEMP_DIR"
    cd "$TEMP_DIR"
    
    go build -o viki ./cmd/sdd
    
    if [ -w "$INSTALL_DIR" ]; then
        mv viki "${INSTALL_DIR}/viki"
    else
        sudo mv viki "${INSTALL_DIR}/viki"
    fi
    
    cd - > /dev/null
    rm -rf "$TEMP_DIR"
}

# Verify installation
verify_install() {
    if command -v viki &> /dev/null; then
        echo ""
        echo -e "${GREEN}╔══════════════════════════════════════════╗${NC}"
        echo -e "${GREEN}║     ✅ Viki installed successfully!      ║${NC}"
        echo -e "${GREEN}╚══════════════════════════════════════════╝${NC}"
        echo ""
        echo -e "${BLUE}Get started:${NC}"
        echo "  viki --help        # Show all commands"
        echo "  viki init \"name\"   # Start a new project"
        echo "  viki dashboard     # Open web UI"
        echo "  viki chat          # Chat with AI"
        echo ""
        echo -e "${BLUE}Set up AI provider:${NC}"
        echo "  viki mcp add openai --provider openai --model gpt-4"
        echo ""
    else
        echo -e "${RED}Installation failed. Please try manual installation.${NC}"
        exit 1
    fi
}

# Main
main() {
    detect_platform
    get_latest_version
    install_viki
    verify_install
}

main
