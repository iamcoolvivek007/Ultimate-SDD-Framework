#!/bin/bash

# Ultimate SDD Framework Build Script
# Cross-compiles binaries for multiple platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="ultimate-sdd-framework"
BINARY_NAME="sdd"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR="build"
SOURCE_DIR="cmd/sdd"

# Supported platforms
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

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

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Using Go version: $GO_VERSION"
}

# Clean build directory
clean() {
    print_info "Cleaning build directory..."
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
}

# Build for specific platform
build_platform() {
    local platform=$1
    local os=$(echo $platform | cut -d'/' -f1)
    local arch=$(echo $platform | cut -d'/' -f2)

    local output_name="$BINARY_NAME"

    # Add .exe extension for Windows
    if [ "$os" = "windows" ]; then
        output_name="$BINARY_NAME.exe"
    fi

    local output_path="$BUILD_DIR/$BINARY_NAME-$VERSION-$os-$arch"

    # Add .exe extension to output path for Windows
    if [ "$os" = "windows" ]; then
        output_path="$output_path.exe"
    fi

    print_info "Building for $os/$arch..."

    # Set environment variables for cross-compilation
    export GOOS=$os
    export GOARCH=$arch
    export CGO_ENABLED=0  # Disable CGO for static binaries

    # Build with version information
    go build \
        -ldflags "-X main.version=$VERSION -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -s -w" \
        -o "$output_path" \
        "./$SOURCE_DIR"

    # Calculate file size
    local size=$(du -h "$output_path" | cut -f1)
    print_success "Built $output_path ($size)"

    # Generate checksum
    if command -v sha256sum &> /dev/null; then
        sha256sum "$output_path" > "$output_path.sha256"
        print_info "Generated checksum: $output_path.sha256"
    fi
}

# Build all platforms
build_all() {
    print_info "Building Ultimate SDD Framework v$VERSION for all platforms..."

    for platform in "${PLATFORMS[@]}"; do
        build_platform "$platform"
    done

    print_success "All builds completed!"
    list_artifacts
}

# Build only for current platform
build_current() {
    local current_os=$(go env GOOS)
    local current_arch=$(go env GOARCH)
    local platform="$current_os/$current_arch"

    print_info "Building for current platform: $platform"
    build_platform "$platform"
}

# List build artifacts
list_artifacts() {
    print_info "Build artifacts:"
    ls -la "$BUILD_DIR"/
}

# Create release archive
create_archive() {
    local archive_name="$PROJECT_NAME-$VERSION.tar.gz"

    print_info "Creating release archive: $archive_name"

    # Copy license and readme if they exist
    cp README.md "$BUILD_DIR/" 2>/dev/null || true
    cp LICENSE "$BUILD_DIR/" 2>/dev/null || true

    # Create archive
    cd "$BUILD_DIR"
    tar -czf "../$archive_name" *
    cd ..

    local size=$(du -h "$archive_name" | cut -f1)
    print_success "Created archive: $archive_name ($size)"
}

# Show usage
usage() {
    cat << EOF
Ultimate SDD Framework Build Script

Usage: $0 [COMMAND]

Commands:
    all         Build for all supported platforms
    current     Build for current platform only
    clean       Clean build directory
    archive     Create release archive (run after build)
    help        Show this help message

Supported platforms:
    linux/amd64, linux/arm64
    darwin/amd64, darwin/arm64 (macOS)
    windows/amd64, windows/arm64

Examples:
    $0 all          # Build for all platforms
    $0 current      # Build for current platform
    $0 clean        # Clean build artifacts
    $0 all && $0 archive  # Build everything and create archive

EOF
}

# Main script logic
main() {
    check_go

    case "${1:-all}" in
        "all")
            clean
            build_all
            ;;
        "current")
            build_current
            ;;
        "clean")
            clean
            ;;
        "archive")
            create_archive
            ;;
        "help"|"-h"|"--help")
            usage
            exit 0
            ;;
        *)
            print_error "Unknown command: $1"
            usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"