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

## Install

After v0.1 is published, install its tagged source:

```powershell
go install github.com/TrandomHL/AgentExecTrace@v0.1.0
```

`go install` builds from source and identifies itself as `0.1.0`; use a release
asset when you need the release-workflow build provenance.

## Commands

| Command | Purpose |
|---|---|
| `snapshot [--output file]` | Emit OS/WSL evidence, CWD, path namespace, PATH and PATHEXT metadata. It never dumps environment values. |
| `resolve [--output file] <name>` | Explain candidate existence and executability in PATH/PATHEXT order, including lightweight provenance. |
| `probe [--max-bytes n] [--output file] [-- <command> [args...]]` | With no command, run a deterministic same-binary self-probe; otherwise capture the supplied argv, bounded stdout/stderr, UTF-8 status and exit result. |
| `diff [--output file] <left.json> <right.json>` | Compare snapshot, resolve, or probe JSON with explicit semantic findings and priorities. |
| `report --redact [--output file] <input>` | Produce an Issue-ready Markdown report with secret and home-path redaction summary. |

## Examples

```powershell
# Capture a context.
.\agentexectrace.exe snapshot --output windows.json

# Explain one executable, including PATH/PATHEXT order.
.\agentexectrace.exe resolve git --output git.json

# Exercise argv, UTF-8, stdout, stderr and a fixed exit code without a shell.
.\agentexectrace.exe probe

# Keep custom process probing when needed.
.\agentexectrace.exe probe -- cmd.exe /c "git --version"

# Compare the same kind of JSON evidence from two contexts.
.\agentexectrace.exe diff windows.json wsl.json

# Prepare a shareable Markdown copy; inspect it before posting.
.\agentexectrace.exe report --redact --output report.md windows.json
```

## Troubleshooting examples

1. When an agent starts in `C:\` but the terminal starts in the project, run
   `snapshot` in both contexts and `diff` the outputs. This matches [Codex issue
   #20858](https://github.com/openai/codex/issues/20858).
2. When `git` works in PowerShell but not an agent, run `resolve git` in both
   contexts and compare PATHEXT/candidate order. See [Codex issue
   #15174](https://github.com/openai/codex/issues/15174).
3. When a WSL-launched agent appears to use Windows/MINGW tools or UNC paths,
   run `snapshot` in WSL and PowerShell, then compare the two files. A finding
   such as `execution_namespace_changed`, `cwd_changed`, or
   `path_namespace_changed` identifies the boundary. See [Claude Code issue
   #19653](https://github.com/anthropics/claude-code/issues/19653).

## Privacy

`snapshot` reports PATH and PATHEXT as diagnostic evidence, but never exports
all environment variables. Before sharing an artifact, use `report --redact`.
It removes bare/prefixed credential assignments (including `token`, `password`,
`secret`, `key`, and `authorization`), URL credentials, bearer tokens, `sk-`
keys, PEM private-key blocks, JWT-shaped values and home/profile path prefixes.
It also counts redactions by category. Review every report before sharing:
redaction is defense-in-depth, not a guarantee.

## Scope and contributing

v0.1 has no GUI, VS Code extension, MCP, LLM/API, telemetry, cloud backend,
auto-fix, config rewrite, vendor-specific repair, deep ETW/kernel tracing, or
full Linux/macOS doctor. See [PRODUCT_SPEC.md](PRODUCT_SPEC.md),
[REQUIREMENTS.md](REQUIREMENTS.md), and [CONTRIBUTING.md](CONTRIBUTING.md).
