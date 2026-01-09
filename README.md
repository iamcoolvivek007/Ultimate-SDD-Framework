# 🚀 Ultimate SDD Framework

**System over Snippets - AI-Powered Development**

The Ultimate SDD Framework implements the "System over Snippets" philosophy: building AI systems where every conversation has clear purpose, enforced through modular rules, context reset planning, and continuous system evolution.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## ✨ What Makes This Different

Traditional AI coding assistants give you unfiltered, unstructured suggestions. The Ultimate SDD Framework enforces **rigorous development discipline**:

- **Structured Gating**: Features progress through mandatory phases (Specify → Plan → Task → Execute)
- **Expert Personas**: Specialized AI agents for each development role
- **Context Awareness**: LSP integration understands your existing codebase
- **Terminal-Native**: Works entirely in your terminal with rich TUI interfaces
- **API Key Support**: Connect to OpenAI, Anthropic, Google Gemini, Ollama, and Azure

## 🏗️ Architecture Overview

```
Ultimate SDD Framework
├── 🎯 Workflow Layer (Spec Kit)     → Gates: Specify → Plan → Task → Execute
├── 🧠 Intelligence Layer (BMAD)     → Roles: PM, Architect, Developer, QA
├── 💻 Execution Layer (OpenCode)    → CLI + TUI + LSP + MCP
└── 📊 State Management              → .sdd/ folder + YAML persistence
```

## 🚀 Quick Start

### 1. Install

#### Option A: Download Pre-built Binary (Recommended)

```bash
# Download the latest release for your platform
# Visit: https://github.com/iamcoolvivek007/Ultimate-SDD-Framework/releases

# Linux/macOS
curl -L https://github.com/iamcoolvivek007/Ultimate-SDD-Framework/releases/latest/download/sdd-linux-amd64.tar.gz | tar xz
sudo mv sdd /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/iamcoolvivek007/Ultimate-SDD-Framework/releases/latest/download/sdd-windows-amd64.tar.gz" -OutFile "sdd.tar.gz"
tar -xzf sdd.tar.gz
# Add to PATH or use directly
```

#### Option B: Build from Source

```bash
# Clone repository
git clone https://github.com/iamcoolvivek007/Ultimate-SDD-Framework.git
cd Ultimate-SDD-Framework

# Build for current platform
make build

# Or build for all platforms
make build-all

# Install system-wide
sudo make install
```

### 2. Configure AI Provider

```bash
# Add OpenAI provider
sdd mcp add my-openai --provider openai --model gpt-4

# Or Anthropic
sdd mcp add my-claude --provider anthropic --model claude-3-sonnet-20240229

# Test connection
sdd mcp test my-openai
```

### 3. Initialize Project

```bash
# Create .agents/ directory with persona definitions
mkdir -p .agents

# Initialize SDD project
sdd init "My Awesome Project"
```

### 4. Develop Features

```bash
# 1. Specify requirements
sdd specify "I want a user authentication system"

# 2. Plan architecture
sdd plan

# 3. Approve plan (required before proceeding)
sdd approve

# 4. Break down into tasks
sdd task

# 5. Generate implementation guide
sdd execute

# 6. Quality review
sdd review

# 7. Final approval
sdd approve
```

## 🎯 Core Concepts

### The SDD Trinity

1. **Workflow (Spec Kit)**: Enforces development sequence
   - **Specify**: Gather requirements with PM agent
   - **Plan**: Design architecture with Architect agent
   - **Task**: Break down work with Developer agent
   - **Execute**: Guide implementation with Developer agent
   - **Review**: Quality assurance with QA agent

2. **Intelligence (BMAD)**: Expert AI personas
   - **Product Manager**: Requirements & edge cases
   - **System Architect**: Design & technology choices
   - **Software Developer**: Clean code & TDD
   - **Quality Assurance**: Testing & validation

3. **Execution (OpenCode)**: Terminal-native development
   - **CLI**: Command-line interface for all operations
   - **TUI**: Rich terminal UI for interactive workflows
   - **LSP**: Codebase context awareness
   - **MCP**: Multi-provider AI model support

### The "Merged Secret"

**You cannot proceed to task breakdown until a human or QA agent approves the architecture plan.** This prevents poorly designed features from being implemented.

## 📋 Commands Overview

### Project Management
```bash
nexus init <name>           # Initialize SDD project
nexus discovery [--deep]    # Brownfield: Map existing codebase
nexus status                # Show project status
nexus approve               # Approve current phase
nexus team <subcommand>     # Team collaboration management
```

### Development Workflow
```bash
# Greenfield (New Projects)
nexus init <name>           # Initialize project
nexus specify <desc>        # PRD-First requirement gathering
nexus plan                  # Context reset architecture planning
nexus approve               # Quality gates (mandatory approvals)
nexus task                  # Atomic task breakdown
nexus execute               # Modular rule task execution
nexus review [pr]           # AI-powered code review & QA validation
nexus evolve "bug report"   # System evolution from bugs

# Advanced Development Features
nexus analyze               # Comprehensive code quality analysis
nexus pair <subcommand>     # Interactive AI pair programming
nexus learn <subcommand>    # Adaptive learning & personalization

# Brownfield (Existing Codebases)
nexus init <name>           # Initialize project
nexus discovery --deep      # Map existing system (MANDATORY FIRST)
nexus specify <desc>        # Legacy-aware requirement gathering
nexus plan                  # Migration & integration planning
nexus approve               # Quality gates with legacy validation
nexus task                  # File-path-specific task breakdown
nexus execute               # Safeguard execution with regression testing
nexus review                # Legacy integration validation
nexus evolve "bug report"   # Brownfield-specific rule evolution
```

### AI Provider Management
```bash
nexus mcp add <name> --provider <provider>    # Add AI provider
nexus mcp remove <name>                       # Remove provider
nexus mcp list                                # List providers
nexus mcp default <name>                      # Set default provider
nexus mcp test [name]                         # Test connection
nexus mcp chat <message>                      # Direct chat with AI
```

## 🧠 System over Snippets Philosophy

The Ultimate SDD Framework implements five meta-skills that prevent "vibe coding" and ensure every AI interaction serves a specific purpose:

### 1. PRD-First Development
**No coding begins without a validated Product Requirements Document.**
- Every feature starts with comprehensive requirements
- Out-of-scope items are explicitly identified
- Business logic is defined before technical work

### 2. Modular Rules Architecture
**Rules are split by concern to save context and prevent drift.**
- `global.md`: Universal standards (TDD, security, linting)
- `frontend.md`: React/TypeScript specific rules
- `backend.md`: Go/Fiber backend patterns
- `api.md`: REST/GraphQL API standards
- Rules loaded on-demand based on task context

### 3. Command-ification
**Repetitive workflows are mapped to CLI commands.**
- `nexus specify`: PRD generation workflow
- `nexus plan`: Context reset planning
- `nexus execute`: Modular rule task execution
- `nexus evolve`: System evolution from bugs

### 4. Context Reset
**Planning and execution happen in separate mental spaces.**
- Planning: Clean slate with only PRD context
- Execution: Task-specific with relevant rules only
- Prevents context overload and confusion

### 5. System Evolution
**Every bug becomes a learning opportunity.**
- `nexus evolve "bug description"`: Analyzes bugs and updates rules
- Root cause analysis prevents similar issues
- System gets smarter with each bug fixed

### The "Merged Secret" Reinforced
**You cannot proceed from Plan to Task without approval.** This enforced gating, combined with modular rules and system evolution, creates a development environment that continuously improves and prevents common issues.

## 🏭 Brownfield Development Support

The Ultimate SDD Framework includes specialized support for **brownfield development** (working with existing codebases), implementing the "Document-First" and "Reverse Engineering" strategies.

### Brownfield Workflow

```
🔍 Discovery → 📋 Specification → 🏗️ PIV Planning → 🛡️ Safeguard Execution → 🔄 Evolution
```

#### 1. Discovery Phase (`nexus discovery --deep`)
**Map existing codebase and establish system context**
- Comprehensive LSP analysis of all files
- Identification of legacy patterns and anti-patterns
- Mapping of integration points and dependencies
- Assessment of technical debt
- Generation of `CONTEXT.md` as source of truth

#### 2. Specification Phase (`nexus specify "feature with legacy integration"`)
**Define interactions with existing system**
- Legacy touchpoints identification
- Regression risk assessment
- Integration constraint documentation
- Context validation against discovered patterns

#### 3. PIV Planning (`nexus plan`)
**Context-reset planning with brownfield awareness**
- Clean mental space for architecture design
- Migration and integration strategy planning
- Legacy code refactoring requirements
- Backwards compatibility planning

#### 4. Safeguard Execution (`nexus execute`)
**Protected implementation with legacy validation**
- Modular rule loading based on task context
- Automatic regression testing integration
- Legacy pattern compliance enforcement
- Integration point validation

#### 5. Evolution (`nexus evolve "legacy integration bug"`)
**System learns from brownfield challenges**
- Root cause analysis of legacy integration issues
- Rule updates to prevent similar problems
- Pattern evolution based on real implementation experience
- Continuous improvement of brownfield development practices

### Brownfield Benefits

- **Risk Mitigation**: Clear identification of integration points and regression risks
- **Pattern Compliance**: Automatic enforcement of established legacy patterns
- **Context Awareness**: AI agents understand existing system constraints
- **System Evolution**: Framework learns and improves from each brownfield challenge
- **Auditability**: Complete traceability of changes to legacy systems

### Getting Started with Brownfield

```bash
# 1. Discover the existing system
nexus init "My Legacy Project"
nexus discovery --deep

# 2. Review the generated context
cat .sdd/CONTEXT.md

# 3. Specify features with legacy awareness
nexus specify "Add user export feature without breaking existing streak logic"

# 4. Plan with integration constraints
nexus plan  # Includes migration strategies
nexus approve

# 5. Execute with safeguard validation
nexus task
nexus execute  # Validates against legacy patterns

# 6. Review with regression testing
nexus review

# 7. Evolve from any integration issues
nexus evolve "Export feature broke existing habit tracking"
```

## 🧠 Advanced AI Features

### Code Quality Analysis (`nexus analyze`)
**Comprehensive automated code quality assessment:**
- **Code Metrics**: Lines of code, complexity, maintainability scores
- **Security Scanning**: Vulnerability detection and hardening recommendations
- **Performance Analysis**: Bottleneck identification and optimization suggestions
- **Test Coverage**: Gap analysis and testing strategy recommendations
- **Quality Scoring**: Overall A-F grading with detailed breakdowns

```bash
nexus analyze  # Full codebase analysis
# Generates .sdd/analysis_report.md with detailed findings
```

### AI-Powered Code Review (`nexus review [pr-number]`)
**Intelligent automated code review with AI analysis:**
- **Pattern Recognition**: Identifies anti-patterns and best practice violations
- **Security Analysis**: Automated vulnerability and weakness detection
- **Performance Insights**: Optimization opportunities and bottleneck warnings
- **Maintainability Assessment**: Code complexity and refactoring suggestions
- **Quality Scoring**: Automated approval/rejection recommendations

```bash
nexus review 123              # Review PR #123
nexus review --deep           # Comprehensive analysis
# Generates .sdd/review_report.md with detailed feedback
```

### Interactive Pair Programming (`nexus pair`)
**Real-time AI-assisted development sessions:**
- **Context-Aware Suggestions**: Intelligent code completion and refactoring
- **Best Practice Enforcement**: Automatic application of learned patterns
- **Learning Integration**: Session insights improve future suggestions
- **Multiple Modes**: Code completion, refactoring, testing, explanation
- **Session Analytics**: Productivity tracking and pattern recognition

```bash
nexus pair start developer "api development"  # Start session
nexus pair suggest --file api.go --line 42 --type refactor
nexus pair action --id sugg_123 --action accepted
nexus pair report  # View session insights
nexus pair end     # Complete session
```

### Adaptive Learning System (`nexus learn`)
**Continuous improvement through development pattern recognition:**
- **Interaction Recording**: Track successful patterns and learn from failures
- **Personalized Suggestions**: Context-aware recommendations based on history
- **Rule Evolution**: Automatic suggestions for framework rule improvements
- **Pattern Recognition**: Identify and promote successful coding approaches
- **Learning Analytics**: Development habit insights and productivity metrics

```bash
nexus learn record --type refactoring --context api --action extract-method --success
nexus learn suggest database api  # Get personalized recommendations
nexus learn report               # View learning insights
nexus learn evolve               # Suggest rule improvements
```

## 👥 Team Collaboration Features

### Team Management (`nexus team`)
**Collaborative development environment:**
- **Team Structure**: Member management with roles and skills
- **Shared Standards**: Team-wide coding rules and best practices
- **Knowledge Base**: Centralized documentation and pattern library
- **Code Patterns**: Reusable solutions and architectural patterns
- **Decision Log**: Important architectural and design decisions

```bash
nexus team init --name "Backend Team" --description "API development team"
nexus team member add --name "Alice" --role senior --skills "go,api,testing"
nexus team rule add --category coding_standards --title "Use meaningful names"
nexus team knowledge add --title "API Design Patterns" --category best_practices
nexus team pattern add --name "Repository Pattern" --language go --code "..."
nexus team search "error handling"  # Search team knowledge
nexus team report                  # Comprehensive team overview
```

## 🤖 AI Providers

### Supported Providers

| Provider | Status | Models |
|----------|--------|--------|
| OpenAI | ✅ | GPT-4, GPT-4-Turbo, GPT-3.5-Turbo |
| Anthropic | ✅ | Claude 3 Opus/Sonnet/Haiku |
| Google Gemini | ✅ | Gemini Pro, Gemini 1.5 Pro |
| Ollama | ✅ | Local models (Llama, Mistral, etc.) |
| Azure OpenAI | ✅ | GPT-4, GPT-3.5-Turbo |

### Configuration Examples

```bash
# OpenAI
sdd mcp add openai-prod --provider openai --model gpt-4
# Enter API key when prompted

# Anthropic
sdd mcp add claude-dev --provider anthropic --model claude-3-sonnet-20240229

# Google Gemini
sdd mcp add gemini-test --provider google --model gemini-pro

# Ollama (local)
sdd mcp add ollama-local --provider ollama --model llama2

# Azure OpenAI
sdd mcp add azure-prod --provider azure --model gpt-4 \
  --base-url https://your-resource.openai.azure.com/
```

## 🎨 Agent Personas

The framework includes four specialized AI personas in `.agents/`:

### Product Manager (`pm.md`)
- Requirements analysis and edge case identification
- Business logic and user experience focus
- Acceptance criteria definition

### System Architect (`architect.md`)
- Technology stack recommendations
- System design and component architecture
- Scalability and performance considerations

### Software Developer (`developer.md`)
- Clean code principles and best practices
- Test-Driven Development (TDD) guidance
- Implementation patterns and patterns

### Quality Assurance (`qa.md`)
- Testing strategy and coverage analysis
- Security and performance validation
- Code quality and standards enforcement

## 📁 Project Structure

```
your-project/
├── .sdd/                    # SDD state and configuration
│   ├── state.yaml          # Project state and phase tracking
│   ├── spec.md             # Feature specifications
│   ├── plan.md             # Architecture plans
│   ├── tasks.md            # Task breakdowns
│   ├── implementation.md   # Implementation guides
│   ├── review.md           # QA reviews
│   └── mcp.json            # AI provider configurations
├── .agents/                # AI persona definitions
│   ├── pm.md               # Product Manager
│   ├── architect.md        # System Architect
│   ├── developer.md        # Software Developer
│   └── qa.md               # Quality Assurance
└── [your code]             # Your application code
```

## 🔧 Advanced Usage

### Environment Variables

```bash
# Set API key via environment
export SDD_API_KEY=your-api-key-here
sdd mcp add my-provider --provider openai --model gpt-4

# Custom MCP config location
export SDD_CONFIG_DIR=/path/to/config
```

### Custom Agent Personas

Edit `.agents/` files to customize agent behavior:

```yaml
---
role: Custom Developer
expertise: React, TypeScript, Testing
personality: Pragmatic, efficient, detail-oriented
tone: Technical, helpful, solution-focused
---

# Custom Developer Agent

Specialized in React/TypeScript development with focus on:
- Component architecture and state management
- TypeScript best practices
- Testing strategies for React applications
```

### Integration with CI/CD

```yaml
# .github/workflows/sdd.yml
name: SDD Quality Gates
on: [pull_request]

jobs:
  sdd-validation:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup SDD
        run: |
          # Install SDD framework
          # Configure AI providers
          # Run SDD validation
```

## 🎯 Best Practices

### 1. Start Small
Begin with simple features to learn the framework patterns.

### 2. Customize Agents
Adapt the agent personas to your team's preferences and technologies.

### 3. Regular Reviews
Use the approval workflow to ensure quality at each phase.

### 4. Context Awareness
Keep your codebase clean - the LSP integration uses your existing code as context.

### 5. Multiple Providers
Configure multiple AI providers for different use cases (e.g., Claude for planning, GPT-4 for implementation).

## 🤝 Contributing

The Ultimate SDD Framework is designed to be extensible:

- **Add New Providers**: Extend `internal/mcp/client.go`
- **Custom Workflows**: Modify phase transitions in `internal/gates/`
- **New Agent Types**: Add personas to `.agents/` directory
- **TUI Enhancements**: Improve Bubble Tea interfaces

## 📚 Documentation

- [Quick Start Guide](docs/quick-start.md)
- [Agent Configuration](docs/agents.md)
- [MCP Integration](docs/mcp.md)
- [Workflow Reference](docs/workflow.md)
- [API Reference](docs/api.md)

## 🔄 Migration from Existing Tools

### From GitHub Spec Kit
```bash
# Your existing spec templates work unchanged
# SDD provides the enforcement layer
```

### From BMAD METHOD
```bash
# Agent definitions are compatible
# SDD adds workflow management and execution
```

### From OpenCode/Crush
```bash
# Terminal-native approach maintained
# SDD adds structured development workflow
```

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Inspired by:
- [GitHub Spec Kit](https://github.com/github/spec-kit) - Structured specifications
- [BMAD METHOD](https://github.com/bmad-sim) - Expert agent personas
- [OpenCode/Crush](https://github.com/opencode-cc) - Terminal-native development

## 📦 Releases & Downloads

### Pre-built Binaries

The Ultimate SDD Framework provides pre-built binaries for all major platforms:

- **Linux**: `amd64`, `arm64`
- **macOS**: `amd64` (Intel), `arm64` (Apple Silicon)
- **Windows**: `amd64`, `arm64`

Download from: [GitHub Releases](https://github.com/iamcoolvivek007/Ultimate-SDD-Framework/releases)

### Release Channels

- **Latest**: Most recent stable release
- **Pre-release**: Beta versions with new features
- **Nightly**: Development builds (unstable)

### Verification

All releases include SHA256 checksums for verification:

```bash
# Verify download integrity
sha256sum -c sdd-linux-amd64.sha256
```

### Build from Source

For development or custom builds:

```bash
# Quick build for current platform
make build

# Cross-platform builds
make build-all

# Create release archives
make release
```

## 🤝 Contributing

### Development Setup

```bash
# Clone repository
git clone https://github.com/iamcoolvivek007/Ultimate-SDD-Framework.git
cd Ultimate-SDD-Framework

# Install dependencies
make deps

# Run tests
make test

# Start development
make dev
```

### Building Releases

The project uses automated CI/CD for releases:

1. **Push tagged commits** (`git tag v1.0.0 && git push --tags`)
2. **GitHub Actions** automatically builds for all platforms
3. **Release artifacts** are uploaded to GitHub Releases
4. **Checksums** are generated for verification

### Release Process

```bash
# Create and push a new tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will:
# 1. Build binaries for all platforms
# 2. Generate checksums
# 3. Create GitHub release
# 4. Upload artifacts
```

---

**The Ultimate SDD Framework** - Because great software deserves great process.