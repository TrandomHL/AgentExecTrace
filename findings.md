# AgentExecTrace — Findings

## Requirements captured

- Product: AgentExecTrace (working name), vendor-neutral diagnostic CLI for AI
  coding-agent execution-context mismatches.
- Primary platform: Windows 11 and WSL2; Go; one executable; standard library
  preferred; no service, database, network dependency, Node, or Python runtime.
- v0.1 commands: `snapshot`, `resolve`, `probe`, `diff`, and `report --redact`.
- Required signals: Windows/WSL identity, CWD, PATH/PATHEXT, executable
  resolution, path namespace, argv, stdout, stderr, exit code and UTF-8.
- Privacy: no default environment dump; shareable reports redact common secrets.
- No mutation of PATH, registry, WSL settings, execution policy or user config.

## Repository discovery

- `E:\project_Open` existed as an empty directory and was initialized as a new
  Git repository on `main` after inspection.
- The originally used `E:\project\_Open` is a separate path. It contains only
  an empty Git initialization caused by the corrected-path typo; it remains
  untouched after the correction.
- Requested migration source `E:\project\_AD9361` was absent when inspected.
  Consequently neither requested source file can be copied or modified yet.
- No Go toolchain is available on PATH or at the common Program Files,
  LocalAppData or Scoop locations. The bundled workspace runtime supplies Git,
  Node and Python, but not Go. This blocks compiling or executing the required
  Go test-first foundation until an approved Go toolchain is available.

## Public issue evidence to verify and cite in requirements

- OpenAI Codex #20858: Windows Desktop agent command runner may use `C:\` while
  the integrated terminal uses the project directory. Supports context snapshot
  and CWD comparison.
- OpenAI Codex #19171: Windows PATH/Path environment construction issue.
  Supports explicitly observing PATH and case-insensitive environment behavior.
- Anthropic Claude Code #19653: native WSL invocation reported commands running
  in a Windows/MINGW environment with UNC/native path mismatch. Supports
  platform and path-namespace diagnostics.
- Anthropic Claude Code #26430: Git Bash command execution succeeded but output
  pipe capture returned empty output and failure. Supports argv/stdout/stderr/
  exit-code probing.
- OpenAI Codex #20858 and Claude Code #19653 together support semantic context
  comparison rather than raw-text-only diffs.

## Technical decisions

| Decision | Rationale |
|---|---|
| Snapshot JSON schema is explicit and versioned | Enables stable comparison without a database or service. |
| Probe output is bounded | Prevents accidental secret disclosure and oversized reports. |
| Redaction is name- and value-pattern based | Defends common secrets even when input data is user-provided. |
| WSL is detected, not reconfigured | Diagnoses namespace evidence without unsafe mutation. |

## Migration record

No migration has been performed. When `E:\project\_AD9361` is supplied or
becomes available, copy `AGENTS.md` and `.codex\config.toml` first, retain all
general clauses, then append only project-specific clauses. Record every added,
changed, or incompatible clause and its reason in `MIGRATION_NOTES.md`.
