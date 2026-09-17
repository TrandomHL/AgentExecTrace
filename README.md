# AgentExecTrace

AgentExecTrace is a vendor-neutral CLI for comparing the execution context seen
by a developer and an AI coding agent. v0.1 targets Windows 11 and WSL2 and is
diagnostic-only: it never rewrites PATH, registry, WSL, execution policy or
agent configuration.

## Quick start

### Windows: download a release

Download the Windows amd64 ZIP and `checksums.txt` from the
[Releases page](https://github.com/TrandomHL/AgentExecTrace/releases). Compare
`Get-FileHash .\agentexectrace-windows-amd64.zip -Algorithm SHA256` with the
matching checksum, then extract the ZIP. Rename the extracted executable once
so that all Windows examples below use the same name:

```powershell
Rename-Item .\agentexectrace-windows-amd64.exe agentexectrace.exe
.\agentexectrace.exe snapshot
.\agentexectrace.exe resolve git
```

To use the command from any directory, put the extracted executable in a
directory already on `PATH`, or add its directory to `PATH` using your normal
Windows settings. AgentExecTrace does not change `PATH` for you.

### Linux or WSL2: download a release

Download the Linux amd64 tarball from the [Releases
page](https://github.com/TrandomHL/AgentExecTrace/releases), extract the
`agentexectrace-linux-amd64` binary, and run:

```bash
chmod +x ./agentexectrace-linux-amd64
./agentexectrace-linux-amd64 snapshot
./agentexectrace-linux-amd64 resolve git
```

### The smallest useful comparison

For two Windows contexts, use the same executable and a shared output directory.
In your **developer terminal**, substitute the executable path below. Keep your
current project working directory; do not switch to the download folder:

```powershell
$tool = (Resolve-Path '<absolute path to agentexectrace.exe>').Path
$evidence = Join-Path (Get-Location) 'aet-evidence'
New-Item -ItemType Directory -Path $evidence -Force | Out-Null
& $tool snapshot --output (Join-Path $evidence 'terminal.json')
$tool
$evidence
```

Give the two printed absolute paths to the agent. In the **agent's command
runner**, substitute those paths below. Do not change its working directory:
that difference may be the evidence you are looking for.

```powershell
$tool = '<absolute path printed above to agentexectrace.exe>'
$evidence = '<absolute path printed above to aet-evidence>'
& $tool snapshot --output (Join-Path $evidence 'agent.json')
& $tool diff --output (Join-Path $evidence 'diff.json') (Join-Path $evidence 'terminal.json') (Join-Path $evidence 'agent.json')
& $tool report --redact --output (Join-Path $evidence 'report.md') (Join-Path $evidence 'diff.json')
```

Read `report.md` before sharing it. `[]` or `null` in the raw diff means no
reported differences for the compared fields, not proof that every aspect of
the environments matches. Raw JSON stays local. If a sandbox cannot access the
shared directory, use an approved accessible location; do not disable isolation.

#### Windows terminal versus a native WSL agent

First capture `terminal.json` using the Windows steps above. In the **WSL
agent runner**, use the Linux executable from the same release. Substitute its
absolute Linux path and the Windows evidence-directory path printed above:

```bash
tool='/absolute/Linux/path/to/agentexectrace-linux-amd64'
evidence="$(wslpath -u '<Windows absolute path to aet-evidence>')"
"$tool" snapshot --output "$evidence/wsl.json"
"$tool" diff --output "$evidence/windows-wsl-diff.json" "$evidence/terminal.json" "$evidence/wsl.json"
"$tool" report --redact --output "$evidence/windows-wsl-report.md" "$evidence/windows-wsl-diff.json"
```

Keep the WSL runner's working directory unchanged. Both contexts must be allowed
to access the shared folder. Running a Windows `.exe` through WSL interop measures
a Windows process, not the native Linux context. Windows/WSL platform and path
namespace differences are expected; they are not by themselves a fault.
For two native Linux contexts, use the same Linux binary, separate snapshot
filenames and a shared absolute Linux evidence path (no `wslpath` conversion).

Use `resolve git` when executable lookup is suspicious, and `probe` when the
command resolves but process behavior differs.

## AI agent setup

For a copyable instruction snippet and a short first-run prompt, see
[AI_AGENT_SETUP.md](AI_AGENT_SETUP.md). The snippet can be placed in
`AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, or an equivalent agent-instructions
file.

## Reproduce the diagnostic value

Run three [controlled experiments](examples/README.md) from source:

```text
go test -count=1 -v -run "^TestControlledExamples$" .
```

They check CWD-dependent file access, Windows PATHEXT lookup failure, and PATH
order selecting different executables. Each verifies the symptom, the diagnostic
finding, and a restored-context control. These are synthetic experiments, not
external adoption evidence. Windows runs all three; Linux skips PATHEXT.

## Interactive code graph

Explore the latest committed code-structure snapshot in the [interactive
AgentExecTrace graph](https://trandomhl.github.io/AgentExecTrace/). The raw
[graph JSON](graphify-out/graph.json) and [graph report](graphify-out/GRAPH_REPORT.md)
are also stored in the repository. The Pages snapshot is regenerated and
deployed when updated graph files reach `main`.

## Build from source

Go 1.26 or newer is required to build from source.

```powershell
go build -o agentexectrace.exe .
.\agentexectrace.exe snapshot
.\agentexectrace.exe resolve git
.\agentexectrace.exe probe -- cmd.exe /c "git --version"
```

The executable has no Node, Python, service, database or network runtime
dependency.

## Install with Go

Go 1.26 or newer is required:

```powershell
go install github.com/TrandomHL/AgentExecTrace@v0.1.0
```

`go install` builds from source and identifies itself as `0.1.0`; use a release
asset when you need the release-workflow build provenance.

Go installs `AgentExecTrace.exe` on Windows and `AgentExecTrace` on Linux
into its normal binary directory. Ensure that directory is on `PATH`. When
using this installation, replace `./agentexectrace-linux-amd64` or
`.\agentexectrace.exe` in examples with the installed command name.

## Upgrade and uninstall

For a release install, download the new archive and replace the old extracted
binary after checking its checksum. For a Go install, run the installation
command again with the desired version or `@latest`. To uninstall, remove the
binary or Go-installed executable; no service or configuration is created.

## Commands

| Command | Purpose |
|---|---|
| `snapshot [--output file]` | Emit OS/WSL evidence, CWD, path namespace, PATH and PATHEXT metadata. It does not dump the full environment. |
| `resolve [--output file] <name>` | Explain candidate existence and executability in PATH/PATHEXT order, including lightweight provenance. |
| `probe [--max-bytes n] [--output file] [-- <command> [args...]]` | With no command, run a deterministic same-binary self-probe; otherwise capture the supplied argv, bounded stdout/stderr, UTF-8 status and exit result. |
| `diff [--output file] <left.json> <right.json>` | Compare snapshot, resolve, or probe JSON with explicit semantic findings and priorities. |
| `report --redact [--output file] <input>` | Produce an Issue-ready Markdown report with secret and home-path redaction summary. |

## Examples

```powershell
# Capture a context.
.\agentexectrace.exe snapshot --output windows.json

# Explain one executable, including PATH/PATHEXT order.
.\agentexectrace.exe resolve --output git.json git

# Exercise argv, UTF-8, stdout, stderr and a fixed exit code without a shell.
.\agentexectrace.exe probe

# Keep custom process probing when needed.
.\agentexectrace.exe probe -- cmd.exe /c "git --version"

# Compare the snapshots captured in the first-use workflow.
.\agentexectrace.exe diff .\aet-evidence\terminal.json .\aet-evidence\agent.json

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
