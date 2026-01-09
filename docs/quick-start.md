# 🚀 Quick Start Guide

Get up and running with the Ultimate SDD Framework in 10 minutes.

## Prerequisites

- Go 1.21+ installed
- Git repository initialized
- Basic terminal knowledge

## Step 1: Install SDD

```bash
# Clone the framework
git clone <repository-url>
cd ultimate-sdd-framework

# Build the CLI
go build -o sdd ./cmd/sdd

# Add to PATH (optional)
sudo mv sdd /usr/local/bin/
```

## Step 2: Configure AI Provider

Choose your preferred AI provider:

### OpenAI (Recommended)
```bash
sdd mcp add openai-main --provider openai --model gpt-4
# Enter your OpenAI API key when prompted
```

### Anthropic Claude
```bash
sdd mcp add claude-main --provider anthropic --model claude-3-sonnet-20240229
# Enter your Anthropic API key
```

### Test Connection
```bash
sdd mcp test openai-main
```

Expected output:
```
Testing connection to OpenAI...
✅ Connection successful!
```

## Step 3: Initialize Project

```bash
# Create a new project
sdd init "User Authentication System"

# Verify setup
sdd status
```

Expected output:
```
🚀 Ultimate SDD Framework Status
=================================
Project: User Authentication System
Current Phase: init
Last Updated: [timestamp]

Phase Status:
  init: ✅ APPROVED
  specify: ○ PENDING
  plan: ○ PENDING
  task: ○ PENDING
  execute: ○ PENDING
  review: ○ PENDING
  complete: ○ PENDING

Next Steps:
  Run: sdd specify "your feature description"
```

## Step 4: Develop Your First Feature

### Specify Requirements
```bash
sdd specify "I want users to be able to register, login, and logout with email and password"
```

This creates `spec.md` with detailed requirements.

### Design Architecture
```bash
sdd plan
```

This creates `plan.md` with the technical architecture.

### Approve the Plan (Required!)
```bash
sdd approve
```

**Important:** You cannot proceed without this approval!

### Break Down Tasks
```bash
sdd task
```

This creates `tasks.md` with specific implementation tasks.

### Generate Implementation Guide
```bash
sdd execute
```

This creates `implementation.md` with development guidance.

### Quality Review
```bash
sdd review
```

This creates `review.md` with quality assessment.

### Final Approval
```bash
sdd approve
```

Congratulations! Your feature is now complete.

## Step 5: Check Your Work

```bash
sdd status
```

Expected output:
```
🚀 Ultimate SDD Framework Status
=================================
Project: User Authentication System
Current Phase: complete

Phase Status:
  init: ✅ APPROVED
  specify: ✅ APPROVED
  plan: ✅ APPROVED
  task: ✅ APPROVED
  execute: ✅ APPROVED
  review: ✅ APPROVED
  complete: ✅ APPROVED
```

## File Structure Created

```
your-project/
├── .sdd/
│   ├── state.yaml          # Project state tracking
│   ├── spec.md             # Feature specifications
│   ├── plan.md             # Architecture plans
│   ├── tasks.md            # Task breakdowns
│   ├── implementation.md   # Development guides
│   ├── review.md           # Quality reports
│   └── mcp.json            # AI provider configs
└── .agents/                # AI personas (auto-created)
    ├── pm.md               # Product Manager
    ├── architect.md        # System Architect
    ├── developer.md        # Software Developer
    └── qa.md               # Quality Assurance
```

## Next Steps

### Start a New Feature
```bash
sdd specify "I want user profile management with avatar uploads"
sdd plan
sdd approve
# ... continue workflow
```

### Customize Agents
```bash
# Edit agent personas for your team
vim .agents/pm.md
vim .agents/architect.md
```

### Add More Providers
```bash
# Add backup provider
sdd mcp add claude-backup --provider anthropic --model claude-3-haiku-20240307

# Test all providers
sdd mcp test
sdd mcp test claude-backup
```

## Troubleshooting

### "Project not initialized"
```bash
sdd init "Your Project Name"
```

### "No AI providers configured"
```bash
sdd mcp add my-provider --provider openai --model gpt-4
```

### "Cannot proceed to next phase"
```bash
# Check current status
sdd status

# If stuck on plan approval
sdd approve
```

### "Connection failed"
```bash
# Check API key
sdd mcp test your-provider

# Re-add with correct key
sdd mcp remove your-provider
sdd mcp add your-provider --provider openai --model gpt-4
```

## Advanced Quick Start

### Environment Variables
```bash
export SDD_API_KEY=your-openai-api-key
sdd mcp add openai-env --provider openai --model gpt-4
```

### Custom Agent Setup
```bash
# Copy default agents
cp -r .agents .agents.backup

# Customize for your tech stack
echo "Specialize in Go, PostgreSQL, and React" >> .agents/architect.md
```

### CI/CD Integration
```yaml
# .github/workflows/sdd.yml
name: SDD Quality Check
on: [pull_request]

jobs:
  sdd-validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Validate SDD
        run: |
          ./sdd status
          ./sdd mcp test
```

## Learning Resources

- [Full Documentation](../README.md) - Complete framework guide
- [Agent Configuration](agents.md) - Customize AI personas
- [MCP Integration](mcp.md) - AI provider management
- [Workflow Reference](workflow.md) - Detailed phase guide

---

**🎉 You're now ready to build software with AI-powered rigor!**

The Ultimate SDD Framework ensures every feature is properly specified, designed, implemented, and validated. No more "vibe coding" - every line of code serves a clear purpose in a well-architected system.