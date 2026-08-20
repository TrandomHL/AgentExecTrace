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

## Phase 5 baseline — 2026-08-20

- The only working repository is `E:\project_Open`. `E:\project\_Open` is a
  prior empty typo-path repository and must remain untouched. The policy source
  is `E:\project_AD9361`; its migrated policy/configuration is already present
  in the working repository.
- Local `main` and `origin/main` both point to `b4c13f0`; GitHub Actions run
  `32346400206` completed successfully for that commit on Windows and Ubuntu.
- Current source does not yet conform to Phase 5: `go.mod` uses the old module
  identity; `diff` uses raw `reflect.DeepEqual`; `probe` rejects no arguments;
  `report --redact` is text-only; resolver lacks provenance and POSIX execute
  checks; CI lacks gofmt/vet; no FR-01…FR-07 evidence matrix exists.
- The existing release workflow creates archives/checksums on tag pushes. No
  tag, GitHub Release, issue, PR or further remote write is authorized by
  Phase 5.

### Phase 5 validation evidence

- Windows local validation passed gofmt, vet, test, build and all five CLI
  commands. `go list -m` reports `github.com/TrandomHL/AgentExecTrace`, used as
  the local module-identity equivalent to a tagged `go install` validation.
- Ubuntu WSL has no Go toolchain, so it cannot supply Linux test/vet evidence.
  A Windows cross-built Linux amd64 binary ran all five commands in that WSL and
  reported `os=linux`, `is_wsl=true`, CWD `/mnt/e/project_Open`, and namespace
  `wsl-mount`. The temporary binary was removed after the run.

## Phase 5 CI correction — 2026-08-20

- CI run `32350238908` failed for two verified causes: Linux lost output from a
  fast self-probe because `cmd.Wait` closed `StdoutPipe`/`StderrPipe` before the
  capture goroutines completed; Windows checked out Go files as CRLF and made
  `gofmt -l` report every file.
- `probe.Run` now drains both pipes before waiting for the child, and
  `.gitattributes` fixes Go source checkout line endings to LF.
- With Go 1.26.0 newly available in Ubuntu WSL, direct WSL `go test ./...`,
  gofmt, vet, build and all five CLI commands passed.

## Phase 5.1 release-gate baseline — 2026-08-20

- Independent release-gate follow-up starts from clean `main`/`origin/main`
  commit `8948c2a1fbf2739c04fb213c338b62f3dd4aeb2f` in isolated worktree branch
  `agent/phase-5-1-release-gate-closure`.
- Confirmed blockers: bare `token`, `password`, `secret`, `key`, and
  `authorization` fields evade parts of redaction; POSIX no-extension files
  are overclaimed as native executables; resolver diff ignores same-target
  provenance changes; release asset glob is not multiline and has no unmatched
  asset failure gate.
- The policy source `E:\project_AD9361` is on dirty `master` and is read-only
  for this task. Its migrated general policy is already present in the target.
