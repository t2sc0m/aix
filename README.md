# aix — AI eXchange CLI

Delegate tasks to AI CLI tools with context file injection.

```
aix "Review this code for bugs" -c expert-prompt.md -f target.go
```

## Why

MCP-based AI delegation has overhead — server setup, protocol negotiation, context window consumption from tool definitions. `aix` replaces this with direct subprocess calls. Expert prompts are injected via context files (`-c` flag), not protocol layers.

## Install

```bash
go install github.com/t2sc0m/aix@latest
```

Requires: [Codex CLI](https://github.com/openai/codex) installed and authenticated.

## Usage

```bash
# Basic prompt
aix "Explain this error: connection refused"

# With expert context file (replaces MCP developer-instructions)
aix "Review this plan" -c prompts/plan-reviewer.md -f plan.md

# Attach multiple files
aix "Find security issues" -c prompts/security-analyst.md -f auth.go -f handler.go

# Stdin pipe
echo "Why is this slow?" | aix -f slow-query.sql

# Model override
aix "Analyze this" -m o3

# Sandbox mode
aix "Fix this bug" -s workspace-write -f broken.go

# Check status
aix status
```

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--context` | `-c` | | Context file (expert prompt injection) |
| `--file` | `-f` | | Attach file(s), repeatable |
| `--model` | `-m` | | Model override |
| `--sandbox` | `-s` | `read-only` | `read-only`, `workspace-write`, `danger-full-access` |
| `--cwd` | | | Working directory for codex |
| `--raw` | | `false` | Passthrough stdout/stderr |
| `--timeout` | `-t` | `300` | Timeout in seconds |

## Config

Optional `~/.config/aix/config.yaml`:

```yaml
timeout: 600
sandbox: read-only
adapters:
  codex:
    enabled: true
    model: o3
```

Precedence: CLI flags > config.yaml > defaults.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Codex not installed |
| 3 | Auth failure |
| 4 | Codex exec failure |
| 5 | Timeout |

## Architecture

```
cmd/        Cobra CLI commands (ask, status)
adapter/    AI CLI backend interface + Codex implementation
prompt/     Prompt assembly with size validation
runner/     Subprocess abstraction (enables mock testing)
config/     YAML config loader with merge precedence
```

## License

MIT
