#!/bin/bash

# 🚀 Ultimate SDD Framework - Publishing Script
# This script helps automate the publishing process

set -e

echo "🚀 Ultimate SDD Framework - Publishing Script"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO_URL="https://github.com/iamcoolvivek007/Ultimate-SDD-Framework.git"
VERSION="v1.0.0"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Pre-flight checks
print_step "Running pre-flight checks..."

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -f "cmd/sdd/main.go" ]; then
    print_error "Not in the Ultimate SDD Framework directory. Please run from the project root."
    exit 1
fi

# Check if git repository is clean
if [ -n "$(git status --porcelain)" ]; then
    print_error "Git repository has uncommitted changes. Please commit or stash them first."
    git status
    exit 1
fi

# Check if version tag exists
if git tag -l | grep -q "^${VERSION}$"; then
    print_status "Version tag ${VERSION} already exists"
else
    print_error "Version tag ${VERSION} does not exist. Please create it first."
    exit 1
fi

print_status "Pre-flight checks passed!"

# Build test
print_step "Testing build..."
if ! go build -o nexus ./cmd/sdd; then
    print_error "Build failed! Please fix build issues before publishing."
    exit 1
fi
print_status "Build successful"

# Test version command
print_step "Testing version command..."
if ! ./nexus version; then
    print_error "Version command failed!"
    exit 1
fi
print_status "Version command works"

# Check remote
print_step "Checking git remote..."
if ! git remote get-url origin | grep -q "github.com/iamcoolvivek007/Ultimate-SDD-Framework"; then
    print_warning "Remote URL doesn't match expected repository"
    echo "Current remote: $(git remote get-url origin)"
    echo "Expected: ${REPO_URL}"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Authentication check
print_step "Checking GitHub authentication..."

# Try to access GitHub
if command -v gh &> /dev/null; then
    if gh auth status &> /dev/null; then
        print_status "GitHub CLI authentication detected"
        USE_GH=true
    else
        print_warning "GitHub CLI installed but not authenticated"
        USE_GH=false
    fi
else
    print_warning "GitHub CLI not installed"
    USE_GH=false
fi

# Try HTTPS authentication
if [ "$USE_GH" = false ]; then
    echo "Testing HTTPS authentication..."
    if git ls-remote "${REPO_URL}" &> /dev/null; then
        print_status "HTTPS authentication appears to work"
    else
        print_warning "HTTPS authentication may not be configured"
        echo "Please ensure you have:"
        echo "1. A GitHub Personal Access Token, OR"
        echo "2. SSH keys configured, OR"
        echo "3. GitHub CLI installed and authenticated"
        echo
        read -p "Continue anyway? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
fi

# Ready to publish?
echo
print_step "Ready to publish Ultimate SDD Framework ${VERSION}?"
echo
echo "This will:"
echo "  1. Push commits to GitHub"
echo "  2. Push version tag to GitHub"
echo "  3. Create a GitHub release (if using GitHub CLI)"
echo
read -p "Proceed with publishing? (y/N): " -n 1 -r
echo

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_status "Publishing cancelled"
    exit 0
fi

# Push commits
print_step "Pushing commits to GitHub..."
if ! git push origin main; then
    print_error "Failed to push commits"
    exit 1
fi
print_status "Commits pushed successfully"

# Push tag
print_step "Pushing version tag to GitHub..."
if ! git push origin "${VERSION}"; then
    print_error "Failed to push version tag"
    exit 1
fi
print_status "Version tag pushed successfully"

# Create release if using GitHub CLI
if [ "$USE_GH" = true ]; then
    print_step "Creating GitHub release..."

    # Create temporary release notes file
    RELEASE_NOTES="/tmp/sdd-release-notes.md"
    cat > "$RELEASE_NOTES" << 'EOF'
# 🎉 Ultimate SDD Framework v1.0.0 Release!

## 🚀 What's New

### 🧠 Advanced AI Features
- **AI-Powered Code Review** - Automated PR analysis with security & performance insights
- **Interactive Pair Programming** - Real-time AI assistance with context awareness
- **Adaptive Learning System** - Framework learns and improves from your patterns
- **Multi-Provider AI Support** - OpenAI, Anthropic, Google, Ollama, Azure integration

### 🏭 Enterprise Brownfield Support
- **Complete Legacy Discovery** - Comprehensive mapping of existing codebases
- **CONTEXT.md Generation** - Source of truth for system understanding
- **Safe Modification** - Regression prevention for existing systems
- **Integration Safeguards** - Protected changes to legacy code

### 👥 Team Collaboration Platform
- **Team Member Management** - Roles, skills, and responsibilities
- **Shared Coding Standards** - Team-wide rules and best practices
- **Knowledge Base** - Centralized documentation and patterns
- **Decision Logging** - Architectural decision tracking

### 🎯 Quality Assurance Excellence
- **Automated Security Scanning** - Vulnerability detection
- **Performance Analysis** - Bottleneck identification
- **Code Complexity Assessment** - Maintainability scoring
- **Test Coverage Analysis** - Gap identification and strategies

## 📋 Installation

### Pre-built Binaries (Recommended)
```bash
curl -L https://github.com/iamcoolvivek007/Ultimate-SDD-Framework/releases/download/v1.0.0/install.sh | bash
```

### From Source
```bash
git clone https://github.com/iamcoolvivek007/Ultimate-SDD-Framework.git
cd Ultimate-SDD-Framework
go build -o nexus ./cmd/sdd
```

## 🚀 Quick Start
```bash
nexus init "My Project"
nexus discovery --deep    # For existing codebases
nexus specify "Add feature"
nexus analyze             # Quality assessment
```

---
**Built with ❤️ using Go and Charm - The future of AI-assisted development!**
EOF

    if gh release create "${VERSION}" \
        --title "Ultimate SDD Framework ${VERSION} - The Most Advanced AI-Powered Development Platform" \
        --notes-file "$RELEASE_NOTES" \
        --latest; then
        print_status "GitHub release created successfully!"
        print_status "Visit: https://github.com/iamcoolvivek007/Ultimate-SDD-Framework/releases/tag/${VERSION}"
    else
        print_error "Failed to create GitHub release"
        print_warning "You can create it manually at: https://github.com/iamcoolvivek007/Ultimate-SDD-Framework/releases/new"
    fi

    # Clean up
    rm -f "$RELEASE_NOTES"
else
    print_warning "GitHub CLI not available - skipping automatic release creation"
    print_status "Please create the release manually:"
    print_status "  1. Go to: https://github.com/iamcoolvivek007/Ultimate-SDD-Framework/releases/new"
    print_status "  2. Tag: ${VERSION}"
    print_status "  3. Title: Ultimate SDD Framework ${VERSION} - The Most Advanced AI-Powered Development Platform"
    print_status "  4. Use release notes from PUBLISH_GUIDE.md"
fi

echo
print_status "🎉 Publishing complete!"
echo
echo "Next steps:"
echo "  1. Monitor the GitHub repository for any issues"
echo "  2. Consider announcing on social media, forums, or blogs"
echo "  3. Gather user feedback for v1.1.0"
echo
print_status "The Ultimate SDD Framework is now live! 🚀✨"