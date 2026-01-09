# 🤖 MCP Integration Guide

The Model Context Protocol (MCP) enables the Ultimate SDD Framework to connect with multiple AI providers through a unified interface. This guide covers configuring and managing AI providers.

## Supported Providers

| Provider | Status | API Key Required | Base URL Configurable |
|----------|--------|------------------|----------------------|
| OpenAI | ✅ | Yes | No |
| Anthropic | ✅ | Yes | No |
| Google Gemini | ✅ | Yes | No |
| Ollama | ✅ | No | Yes |
| Azure OpenAI | ✅ | Yes | Yes |

## Quick Setup

### OpenAI (Most Popular)
```bash
sdd mcp add openai-prod --provider openai --model gpt-4
# Enter your OpenAI API key when prompted
```

### Anthropic Claude
```bash
sdd mcp add claude-dev --provider anthropic --model claude-3-sonnet-20240229
# Enter your Anthropic API key
```

### Local Ollama
```bash
# First, install Ollama and pull a model
ollama pull llama2

# Then configure SDD
sdd mcp add ollama-local --provider ollama --model llama2
```

## Provider Configuration

### Adding Providers

```bash
sdd mcp add <name> --provider <provider> [options]
```

**Options:**
- `--provider`: AI provider (openai, anthropic, google, ollama, azure)
- `--model`: Model name (provider-specific)
- `--base-url`: Custom API endpoint (for Azure, custom deployments)
- `--default`: Set as default provider

### Examples

#### OpenAI GPT-4
```bash
sdd mcp add gpt4-prod --provider openai --model gpt-4 --default
```

#### Anthropic Claude 3
```bash
sdd mcp add claude-large --provider anthropic --model claude-3-opus-20240229
sdd mcp add claude-fast --provider anthropic --model claude-3-haiku-20240307
```

#### Google Gemini
```bash
sdd mcp add gemini-pro --provider google --model gemini-pro
sdd mcp add gemini-experimental --provider google --model gemini-1.5-pro-latest
```

#### Azure OpenAI
```bash
sdd mcp add azure-gpt4 --provider azure --model gpt-4 \
  --base-url https://your-resource.openai.azure.com/
```

#### Ollama Local Models
```bash
sdd mcp add codellama --provider ollama --model codellama
sdd mcp add mistral --provider ollama --model mistral
```

## Managing Providers

### List Providers
```bash
sdd mcp list
```

Output:
```
🤖 Configured AI Providers
==========================
my-openai (default) ✅ Enabled
  Provider: OpenAI
  Model: gpt-4

claude-dev ✅ Enabled
  Provider: Anthropic
  Model: claude-3-sonnet-20240229

ollama-local ✅ Enabled
  Provider: Ollama (Local)
  Model: llama2
```

### Set Default Provider
```bash
sdd mcp default claude-dev
```

### Test Provider Connection
```bash
# Test default provider
sdd mcp test

# Test specific provider
sdd mcp test claude-dev
```

### Remove Provider
```bash
sdd mcp remove old-provider
```

## Environment Variables

You can set API keys via environment variables:

```bash
export SDD_API_KEY=your-openai-api-key
sdd mcp add openai-env --provider openai --model gpt-4
```

Or for multiple providers:
```bash
export OPENAI_API_KEY=sk-...
export ANTHROPIC_API_KEY=sk-ant-...
export GOOGLE_API_KEY=...
```

## Direct Chat

Test providers with direct chat:

```bash
# Chat with default provider
sdd mcp chat "Explain the difference between REST and GraphQL"

# Chat with specific provider
sdd mcp chat "Write a Go function to reverse a string" --provider codellama

# Advanced options
sdd mcp chat "Design a user authentication system" \
  --provider claude-dev \
  --temperature 0.8 \
  --max-tokens 2000
```

## Configuration File

MCP configurations are stored in `.sdd/mcp.json`:

```json
{
  "providers": {
    "openai-prod": {
      "provider": "openai",
      "api_key": "sk-...",
      "model": "gpt-4",
      "enabled": true
    },
    "claude-dev": {
      "provider": "anthropic",
      "api_key": "sk-ant-...",
      "model": "claude-3-sonnet-20240229",
      "enabled": true
    }
  },
  "default_provider": "openai-prod"
}
```

## Best Practices

### 1. Multiple Providers for Different Tasks

```bash
# Claude for complex reasoning (planning, architecture)
sdd mcp add claude-planner --provider anthropic --model claude-3-opus-20240229

# GPT-4 for implementation details
sdd mcp add gpt4-coder --provider openai --model gpt-4

# Local models for cost-effective tasks
sdd mcp add ollama-fast --provider ollama --model llama2
```

### 2. Cost Optimization

```bash
# Use faster/cheaper models for routine tasks
sdd mcp add gpt3-fast --provider openai --model gpt-3.5-turbo

# Reserve expensive models for complex reasoning
sdd mcp add claude-expensive --provider anthropic --model claude-3-opus-20240229
```

### 3. Reliability

```bash
# Configure backup providers
sdd mcp add openai-backup --provider openai --model gpt-3.5-turbo
sdd mcp add claude-backup --provider anthropic --model claude-2

# Test all providers regularly
sdd mcp test openai-backup
sdd mcp test claude-backup
```

## Troubleshooting

### Connection Issues

**OpenAI API Key Invalid:**
```bash
# Check your API key
curl -H "Authorization: Bearer YOUR_KEY" https://api.openai.com/v1/models

# Re-add provider with correct key
sdd mcp remove openai-prod
sdd mcp add openai-prod --provider openai --model gpt-4
```

**Anthropic Rate Limited:**
```bash
# Anthropic has stricter rate limits than OpenAI
# Consider using GPT-4 as backup or implement retry logic
```

**Ollama Not Running:**
```bash
# Start Ollama service
ollama serve

# Pull required model
ollama pull llama2

# Test connection
curl http://localhost:11434/api/tags
```

### Model Availability

**Model Not Found:**
```bash
# Check available models for provider
sdd mcp chat "What models do you support?" --provider openai

# Update to available model
sdd mcp add gpt4-updated --provider openai --model gpt-4-turbo
```

### Configuration Issues

**Multiple Defaults:**
```bash
# Only one provider can be default
sdd mcp list
sdd mcp default correct-provider
```

**Provider Disabled:**
```bash
# Check provider status
sdd mcp list

# Re-enable if needed (edit .sdd/mcp.json)
```

## Advanced Configuration

### Custom Base URLs

For corporate firewalls or custom deployments:

```bash
# OpenAI proxy
sdd mcp add openai-proxy --provider openai --model gpt-4 \
  --base-url https://your-proxy.company.com/openai/

# Self-hosted Ollama
sdd mcp add ollama-remote --provider ollama --model llama2 \
  --base-url http://gpu-server.company.com:11434/
```

### API Key Rotation

```bash
# Update API key
sdd mcp remove old-openai
sdd mcp add new-openai --provider openai --model gpt-4

# Update default
sdd mcp default new-openai

# Test new configuration
sdd mcp test
```

## Security Considerations

### API Key Management

- **Never commit API keys** to version control
- Use environment variables for CI/CD
- Rotate keys regularly
- Monitor usage and costs

### Local vs Cloud Models

**Cloud Providers (OpenAI, Anthropic, Google):**
- ✅ Latest models and features
- ✅ High reliability and performance
- ❌ API costs and rate limits
- ❌ Data sent to third parties

**Local Models (Ollama):**
- ✅ No API costs
- ✅ Data stays local
- ✅ Custom fine-tuning possible
- ❌ Slower inference
- ❌ Limited model selection

### Provider Selection Guide

| Use Case | Recommended Provider |
|----------|---------------------|
| Complex reasoning | Claude 3 Opus |
| Code generation | GPT-4, CodeLlama |
| Fast responses | GPT-3.5-Turbo, Claude 3 Haiku |
| Cost-effective | Local Ollama models |
| Enterprise | Azure OpenAI |
| Research | Multiple providers |

## Integration with CI/CD

```yaml
# .github/workflows/sdd-check.yml
name: SDD Provider Health Check
on:
  schedule:
    - cron: '0 */6 * * *'  # Every 6 hours

jobs:
  health-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Test AI Providers
        run: |
          ./sdd mcp test openai-prod || echo "OpenAI failed"
          ./sdd mcp test claude-dev || echo "Claude failed"
          ./sdd mcp test ollama-local || echo "Ollama failed"
```

---

The MCP integration makes the Ultimate SDD Framework provider-agnostic, allowing you to choose the best AI model for each task while maintaining a consistent interface.