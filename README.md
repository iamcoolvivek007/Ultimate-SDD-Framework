# 🚀 Viki - Ultimate SDD Framework

**System over Snippets - AI-Powered Development with 21+ Specialized Agents**

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-3.0.0-purple.svg)](CHANGELOG.md)

The Ultimate SDD Framework implements the "System over Snippets" philosophy: building AI systems where every conversation has clear purpose, enforced through modular rules, context reset planning, and continuous system evolution.

## ✨ What Makes This Different

Traditional AI coding assistants give you unfiltered, unstructured suggestions. Viki enforces **rigorous development discipline** with features from the best frameworks:

- **21+ Specialized Agents**: From PM to Security Analyst, Tech Lead to UX Designer
- **Scale-Adaptive Intelligence**: Level 0-4 adjusts planning depth automatically
- **SQLite Persistence**: Sessions, messages, and file changes tracked in database
- **Structured Gating**: Features progress through mandatory phases (Specify → Plan → Task → Execute)
- **9 AI Tools**: Bash, edit, grep, view, write, patch, fetch, glob, ls
- **Multi-Provider Support**: OpenAI, Anthropic, Google Gemini, Ollama, Azure

### 🆕 New in v3.0 (Latest)

**Integrated best features from BMAD-METHOD, OpenCode, and Spec-Kit:**

| Feature | Source | Description |
|---------|--------|-------------|
| `viki session` | OpenCode | SQLite-persistent chat sessions |
| `viki workflow` | BMAD | Quick/Standard/Enterprise tracks |
| `viki brainstorm` | BMAD | 6 ideation techniques (SCAMPER, Six Hats, etc.) |
| `viki agents` | BMAD | 21+ specialized AI personas |
| `viki constitution` | Spec-Kit | Project governance & principles |
| `viki clarify` | Spec-Kit | Refine specifications with Q&A |
| `viki checklist` | Spec-Kit | Generate quality checklists |

**Also includes:**
- Scale-Adaptive Levels (0-4) for automatic complexity detection
- Custom slash commands (`viki/command` format)
- File change tracking with undo capability
- Workflow engine with step dependencies

### 🔧 v2.0 Features

- **Interactive Chat Mode** (`viki chat`) - Continuous AI conversation
- **Project Templates** (`viki new`) - Go, React, Python, Next.js templates
- **Web Dashboard** (`viki dashboard`) - Browser-based UI
- **Plugin System** (`viki plugin`) - Extend with custom agents
- **Secrets Management** (`viki secrets`) - OS keychain integration
- **Codebase Indexing** (`viki index`) - LSP-like symbol extraction

## 🏗️ Architecture Overview

```
Ultimate SDD Framework v3.0
├── 🎯 Workflow Layer      → Gates: Specify → Plan → Task → Execute
├── 🧠 Intelligence Layer  → 21+ Agents: PM, Architect, Developer, QA, Security...
├── 💻 Execution Layer     → CLI + TUI + LSP + MCP + 9 AI Tools
├── 🗃️ Persistence Layer   → SQLite: Sessions, Messages, File Changes
└── 📊 State Management    → .sdd/ folder + YAML + Database
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

2. **Intelligence (BMAD)**: 21+ Expert AI Personas

   | Category | Agents |
   |----------|--------|
   | **Core** | Product Manager, System Architect, Software Developer, Quality Assurance |
   | **Product** | UX Designer, Scrum Master, Business Analyst |
   | **Engineering** | DevOps, Security Analyst, Tech Lead, Data Architect, API Designer, Frontend Dev, Backend Dev |
   | **Quality** | Test Automation Engineer, Performance Engineer, Code Reviewer |
   | **Operations** | Site Reliability Engineer, Technical Writer |
   | **Creative** | Innovation Catalyst, Debug Specialist |

3. **Execution (OpenCode)**: Terminal-native with AI Tools
   - **CLI**: Command-line interface for all operations
   - **TUI**: Rich terminal UI for interactive workflows
   - **LSP**: Codebase context awareness
   - **MCP**: Multi-provider AI model support
   - **Tools**: bash, edit, grep, view, write, patch, fetch, glob, ls

### The "Merged Secret"

**You cannot proceed to task breakdown until a human or QA agent approves the architecture plan.** This prevents poorly designed features from being implemented.

## 📋 Commands Overview

### Project Management
```bash
viki init <name>           # Initialize SDD project
viki discovery [--deep]    # Brownfield: Map existing codebase
viki status                # Show project status
viki approve               # Approve current phase
viki team <subcommand>     # Team collaboration management
```

### Development Workflow
```bash
# Greenfield (New Projects)
viki init <name>           # Initialize project
viki specify <desc>        # PRD-First requirement gathering
viki plan                  # Context reset architecture planning
viki approve               # Quality gates (mandatory approvals)
viki task                  # Atomic task breakdown
viki execute               # Modular rule task execution
viki review [pr]           # AI-powered code review & QA validation
viki evolve "bug report"   # System evolution from bugs

# Advanced Development Features
viki analyze               # Comprehensive code quality analysis
viki pair <subcommand>     # Interactive AI pair programming
viki learn <subcommand>    # Adaptive learning & personalization

# Brownfield (Existing Codebases)
viki init <name>           # Initialize project
viki discovery --deep      # Map existing system (MANDATORY FIRST)
viki specify <desc>        # Legacy-aware requirement gathering
viki plan                  # Migration & integration planning
viki approve               # Quality gates with legacy validation
viki task                  # File-path-specific task breakdown
viki execute               # Safeguard execution with regression testing
viki review                # Legacy integration validation
viki evolve "bug report"   # Brownfield-specific rule evolution
```

### AI Provider Management
```bash
viki mcp add <name> --provider <provider>    # Add AI provider
viki mcp remove <name>                       # Remove provider
viki mcp list                                # List providers
viki mcp default <name>                      # Set default provider
viki mcp test [name]                         # Test connection
viki mcp chat <message>                      # Direct chat with AI
```

### 🆕 v3.0 Commands
```bash
# Session Management (from OpenCode)
viki session list              # List all chat sessions
viki session switch <id>       # Switch to different session
viki session new "title"       # Create new session
viki session export <id>       # Export to markdown
viki session delete <id>       # Delete session

# Workflow Engine (from BMAD)
viki workflow init             # Analyze project & recommend track
viki workflow status           # Show current progress
viki workflow next             # Execute next step
viki workflow list             # List available tracks

# Brainstorming (from BMAD)
viki brainstorm "topic"        # Classic brainstorm
viki brainstorm --technique reverse "topic"    # Reverse brainstorm
viki brainstorm --technique six_hats "topic"   # Six Thinking Hats
viki brainstorm --technique scamper "topic"    # SCAMPER method
viki brainstorm --technique party_mode "topic" # Multi-agent discussion
viki brainstorm --list         # List all techniques

# Agent Selection (from BMAD)
viki agents                    # List all 21+ agents with details

# Governance (from Spec-Kit)
viki constitution "principles" # Create project constitution
viki constitution --view       # View constitution
viki constitution --amend "change" # Amend constitution

# Specification (from Spec-Kit)
viki clarify                   # Generate clarification questions
viki checklist                 # Generate quality checklists
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
- `viki specify`: PRD generation workflow
- `viki plan`: Context reset planning
- `viki execute`: Modular rule task execution
- `viki evolve`: System evolution from bugs

### 4. Context Reset
**Planning and execution happen in separate mental spaces.**
- Planning: Clean slate with only PRD context
- Execution: Task-specific with relevant rules only
- Prevents context overload and confusion

### 5. System Evolution
**Every bug becomes a learning opportunity.**
- `viki evolve "bug description"`: Analyzes bugs and updates rules
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

#### 1. Discovery Phase (`viki discovery --deep`)
**Map existing codebase and establish system context**
- Comprehensive LSP analysis of all files
- Identification of legacy patterns and anti-patterns
- Mapping of integration points and dependencies
- Assessment of technical debt
- Generation of `CONTEXT.md` as source of truth

#### 2. Specification Phase (`viki specify "feature with legacy integration"`)
**Define interactions with existing system**
- Legacy touchpoints identification
- Regression risk assessment
- Integration constraint documentation
- Context validation against discovered patterns

#### 3. PIV Planning (`viki plan`)
**Context-reset planning with brownfield awareness**
- Clean mental space for architecture design
- Migration and integration strategy planning
- Legacy code refactoring requirements
- Backwards compatibility planning

#### 4. Safeguard Execution (`viki execute`)
**Protected implementation with legacy validation**
- Modular rule loading based on task context
- Automatic regression testing integration
- Legacy pattern compliance enforcement
- Integration point validation

#### 5. Evolution (`viki evolve "legacy integration bug"`)
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
viki init "My Legacy Project"
viki discovery --deep

# 2. Review the generated context
cat .sdd/CONTEXT.md

# 3. Specify features with legacy awareness
viki specify "Add user export feature without breaking existing streak logic"

# 4. Plan with integration constraints
viki plan  # Includes migration strategies
viki approve

# 5. Execute with safeguard validation
viki task
viki execute  # Validates against legacy patterns

# 6. Review with regression testing
viki review

# 7. Evolve from any integration issues
viki evolve "Export feature broke existing habit tracking"
```

## 🧠 Advanced AI Features

### Code Quality Analysis (`viki analyze`)
**Comprehensive automated code quality assessment:**
- **Code Metrics**: Lines of code, complexity, maintainability scores
- **Security Scanning**: Vulnerability detection and hardening recommendations
- **Performance Analysis**: Bottleneck identification and optimization suggestions
- **Test Coverage**: Gap analysis and testing strategy recommendations
- **Quality Scoring**: Overall A-F grading with detailed breakdowns

```bash
viki analyze  # Full codebase analysis
# Generates .sdd/analysis_report.md with detailed findings
```

### AI-Powered Code Review (`viki review [pr-number]`)
**Intelligent automated code review with AI analysis:**
- **Pattern Recognition**: Identifies anti-patterns and best practice violations
- **Security Analysis**: Automated vulnerability and weakness detection
- **Performance Insights**: Optimization opportunities and bottleneck warnings
- **Maintainability Assessment**: Code complexity and refactoring suggestions
- **Quality Scoring**: Automated approval/rejection recommendations

```bash
viki review 123              # Review PR #123
viki review --deep           # Comprehensive analysis
# Generates .sdd/review_report.md with detailed feedback
```

### Interactive Pair Programming (`viki pair`)
**Real-time AI-assisted development sessions:**
- **Context-Aware Suggestions**: Intelligent code completion and refactoring
- **Best Practice Enforcement**: Automatic application of learned patterns
- **Learning Integration**: Session insights improve future suggestions
- **Multiple Modes**: Code completion, refactoring, testing, explanation
- **Session Analytics**: Productivity tracking and pattern recognition

```bash
viki pair start developer "api development"  # Start session
viki pair suggest --file api.go --line 42 --type refactor
viki pair action --id sugg_123 --action accepted
viki pair report  # View session insights
viki pair end     # Complete session
```

### Adaptive Learning System (`viki learn`)
**Continuous improvement through development pattern recognition:**
- **Interaction Recording**: Track successful patterns and learn from failures
- **Personalized Suggestions**: Context-aware recommendations based on history
- **Rule Evolution**: Automatic suggestions for framework rule improvements
- **Pattern Recognition**: Identify and promote successful coding approaches
- **Learning Analytics**: Development habit insights and productivity metrics

```bash
viki learn record --type refactoring --context api --action extract-method --success
viki learn suggest database api  # Get personalized recommendations
viki learn report               # View learning insights
viki learn evolve               # Suggest rule improvements
```

## 👥 Team Collaboration Features

### Team Management (`viki team`)
**Collaborative development environment:**
- **Team Structure**: Member management with roles and skills
- **Shared Standards**: Team-wide coding rules and best practices
- **Knowledge Base**: Centralized documentation and pattern library
- **Code Patterns**: Reusable solutions and architectural patterns
- **Decision Log**: Important architectural and design decisions

```bash
viki team init --name "Backend Team" --description "API development team"
viki team member add --name "Alice" --role senior --skills "go,api,testing"
viki team rule add --category coding_standards --title "Use meaningful names"
viki team knowledge add --title "API Design Patterns" --category best_practices
viki team pattern add --name "Repository Pattern" --language go --code "..."
viki team search "error handling"  # Search team knowledge
viki team report                  # Comprehensive team overview
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