# AgentExecTrace Architecture

## Architectural intent

The v0.1 architecture is deliberately small: one Go module, one CLI binary,
JSON files as interchange format, and Go standard library packages. No daemon,
plugin registry, dependency-injection framework, configuration language, or
network client is needed to collect local diagnostic evidence.

## Proposed layout

```text
.
├── main.go                 # flag parsing and subcommand dispatch
├── internal/model          # versioned snapshot/report structs
├── internal/context        # platform, CWD, PATH/PATHEXT and namespace evidence
├── internal/resolve        # deterministic executable candidate resolution
├── internal/probe          # bounded child-process capture
├── internal/diff           # semantic comparison of model values
├── internal/redact         # name/value-pattern redaction
└── internal/testutil       # helper process and deterministic fixtures only
```

Packages are introduced only where they isolate a distinct v0.1 concern. There
are no public interfaces with one implementation; dependencies are passed as
plain values/functions only when tests need controlled input.

## Data model and flow

```text
local context ──► snapshot.json ──► diff ──► human/JSON findings
      │                │
      └─ resolve/probe ─┴─► report --redact ──► shareable report
```

1. `snapshot` observes allowed local facts and serializes a versioned JSON
   snapshot.
2. `resolve` evaluates a supplied command name against an explicit PATH/PATHEXT
   input model. It does not invoke shell-specific command discovery.
3. `probe` starts only the command supplied after `--`; it uses separate stdout
   and stderr pipes, bounded readers and explicit result fields.
4. `diff` decodes two snapshots and compares structured fields, preserving PATH
   order rather than sorting it away.
5. `report --redact` transforms input to a fresh output; it never writes back to
   the source report.

## Platform design

- Platform detection uses `runtime.GOOS`, `runtime.GOARCH`, documented
  environment hints, and filesystem/path syntax. It reports uncertainty rather
  than guessing a distribution or shell.
- On Windows, resolution uses semicolon-delimited PATH and PATHEXT candidate
  order, preserving explicit extensions and `.exe`/`.cmd`/`.bat` distinctions.
- On Linux/WSL, resolution uses colon-delimited PATH and executable permission
  checks. WSL detection is best-effort and reported as evidence.
- Namespace classification is syntactic and does not rewrite a path or mount
  anything.

## Security boundaries

- Environment values are never bulk-exported by default. A future explicit
  allowlist can be considered only after v0.1, not a generic `--env` dump.
- Probe captures a fixed maximum per stream; data past the limit is discarded
  after draining enough to avoid child deadlock and marked `truncated`.
- Redaction runs before a report is saved as shareable. Keys such as `*_TOKEN`,
  `*_SECRET`, `*_PASSWORD`, `*_KEY`, and known credential formats are replaced
  with typed placeholders.
- Commands never write registry keys, environment settings, WSL config, PATH,
  user/project config, or execution policy.

## Error model

Expected diagnostic failures are data, not tool crashes: unresolved executable,
permission failure, non-zero exit, invalid UTF-8 and output truncation produce
structured records. Invalid command-line syntax or unreadable input files return
a clear non-zero CLI exit with an actionable message.
