# 🚀 Ultimate SDD Framework

**Spec-Driven Development with AI Agents**

The Ultimate SDD Framework merges structured gating, expert AI personas, and terminal-native execution to enforce rigorous development practices. No more "vibe coding" - every feature follows a proven sequence with AI assistance.

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

```bash
# Clone and build
git clone <repository-url>
cd ultimate-sdd-framework
go build -o sdd ./cmd/sdd

# Or download pre-built binary
# (coming soon)
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
sdd init <name>           # Initialize SDD project
sdd status                # Show project status
sdd approve               # Approve current phase
```

### Development Workflow
```bash
sdd specify <description> # Generate feature specifications
sdd plan                  # Create architecture plan
sdd task                  # Break down into tasks
sdd execute               # Generate implementation guide
sdd review                # Quality assurance review
```

### AI Provider Management
```bash
sdd mcp add <name> --provider <provider>    # Add AI provider
sdd mcp remove <name>                       # Remove provider
sdd mcp list                                # List providers
sdd mcp default <name>                      # Set default provider
sdd mcp test [name]                         # Test connection
sdd mcp chat <message>                      # Direct chat with AI
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

---

**The Ultimate SDD Framework** - Because great software deserves great process.