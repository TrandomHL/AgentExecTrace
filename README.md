# AgentExecTrace

AgentExecTrace is a vendor-neutral CLI for comparing the execution context seen
by a developer and an AI coding agent. v0.1 targets Windows 11 and WSL2 and is
diagnostic-only: it never rewrites PATH, registry, WSL, execution policy or
agent configuration.

## Build

Go 1.26 or newer is required to build from source.

```powershell
go build -o agentexectrace.exe .
.\agentexectrace.exe snapshot
.\agentexectrace.exe resolve git
.\agentexectrace.exe probe -- cmd.exe /c "git --version"
```

The executable has no Node, Python, service, database or network runtime
dependency.

## Commands

| Command | Purpose |
|---|---|
| `snapshot [--output file]` | Emit OS/WSL evidence, CWD, path namespace, PATH and PATHEXT metadata. It never dumps environment values. |
| `resolve [--output file] <name>` | Explain executable candidates in PATH/PATHEXT order. Candidate output is bounded and marked when truncated. |
| `probe [--max-bytes n] [--output file] -- <command> [args...]` | Capture argv, bounded stdout/stderr, UTF-8 status and exit result. |
| `diff [--output file] <left.json> <right.json>` | Compare two snapshots semantically. |
| `report --redact [--output file] <input>` | Write or print a new redacted report. |

## Troubleshooting examples

1. When an agent starts in `C:\` but the terminal starts in the project, run
   `snapshot` in both contexts and `diff` the outputs. This matches [Codex issue
   #20858](https://github.com/openai/codex/issues/20858).
2. When `git` works in PowerShell but not an agent, run `resolve git` in both
   contexts and compare PATHEXT/candidate order. See [Codex issue
   #15174](https://github.com/openai/codex/issues/15174).
3. When a WSL-launched agent appears to use Windows/MINGW tools or UNC paths,
   compare platform, CWD and path namespace. See [Claude Code issue
   #19653](https://github.com/anthropics/claude-code/issues/19653).

## Privacy

`snapshot` reports PATH and PATHEXT as diagnostic evidence, but never exports
all environment variables. Before sharing an artifact, use `report --redact`.
It removes common credential assignments, bearer tokens, `sk-` keys, PEM
private-key blocks and JWT-shaped values. Review every report before sharing:
redaction is defense-in-depth, not a guarantee.

## Scope and contributing

v0.1 has no GUI, VS Code extension, MCP, LLM/API, telemetry, cloud backend,
auto-fix, config rewrite, vendor-specific repair, deep ETW/kernel tracing, or
full Linux/macOS doctor. See [PRODUCT_SPEC.md](PRODUCT_SPEC.md),
[REQUIREMENTS.md](REQUIREMENTS.md), and [CONTRIBUTING.md](CONTRIBUTING.md).
