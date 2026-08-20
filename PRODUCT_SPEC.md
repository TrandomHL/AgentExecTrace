# AgentExecTrace Product Specification

**Status:** Phase 5.1 — v0.1 Release Gate Closure (scope frozen)
**Target:** Windows 11 + WSL2
**Implementation:** Go, single executable, standard library first
**Audience:** developers and maintainers diagnosing AI coding-agent command
execution differences

## Product statement

AgentExecTrace is a vendor-neutral command-line diagnostic tool that collects,
compares, and safely shares evidence about execution contexts. It answers:

> Why does a command succeed in one developer/agent context but fail or behave
> differently in another?

It is not a repair utility or an agent-vendor health check. It produces evidence
that a user can attach to an issue or use to reproduce a mismatch.

## v0.1 user outcomes

| Command | User outcome |
|---|---|
| `snapshot` | Capture a bounded, structured description of the current context. |
| `resolve <name>` | Explain how an executable name resolves under the current PATH/PATHEXT and namespace. |
| `probe [-- <command> [args...]]` | With no command, run a deterministic self-probe; otherwise capture argv, launch result, bounded stdout/stderr, exit code and UTF-8 status. |
| `diff <left.json> <right.json>` | Report meaningful context changes rather than an unstructured text diff. |
| `report --redact <input>` | Produce a shareable copy with common secrets removed. |

## Supported context evidence

- OS and architecture; Windows versus WSL2 identification.
- CWD and its classified path namespace: Windows drive, UNC, WSL Linux path,
  mounted Windows path, or unknown.
- shell/process identity where discoverable without extra dependencies.
- PATH segment metadata and PATHEXT values, without copying all environment
  variables by default.
- executable candidates, selected path and reason for selection.
- probe argv, captured output bytes/UTF-8 validity, truncation state and exit
  result.

The self-probe verifies OS/process-layer transport facts only. It cannot claim
to observe what an AI model ultimately received; users may compare its recorded
facts with agent output separately.

## Privacy and safety contract

- Default `snapshot` does **not** dump all environment variables.
- Output is local only. v0.1 contains no telemetry, cloud service, database, or
  network dependency.
- `report --redact` is mandatory before a report is labelled shareable. It
  redacts values associated with bare and prefixed credential-like field names
  (`token`, `password`, `secret`, `key`, `authorization`, and `*_TOKEN`-style
  names), URL credentials, and recognized token/private-key patterns; it also
  reports the number and category of redactions without displaying the values.
- Probe output is size-bounded; reports identify truncation rather than silently
  omitting data.
- Every command is read-only with respect to user machine configuration.

## Explicit non-goals for v0.1

No GUI; VS Code extension; MCP; LLM/API; telemetry; cloud backend; auto-fix;
configuration rewrite; execution-policy modification; PATH, registry or WSL
configuration mutation; deep ETW/kernel tracing; vendor-specific repair; or
full Linux/macOS doctor.

## Acceptance gate: v0.1

v0.1 is acceptable only when all five commands work with documented, bounded
JSON/text outputs; Windows and WSL path detection and Windows PATH/PATHEXT
resolution are tested; probe preserves distinct stdout/stderr and exit-code
facts; semantic diff highlights useful differences; report redaction passes
secret fixtures; and unit/integration plus Windows/Linux CI pass. A controlled
WSL verification is required if CI cannot supply WSL.

## Adoption gate: 6–8 weeks after release

The project may claim early validation only with evidence, not fabricated
metrics: at least three independently reproducible troubleshooting reports,
maintainer-confirmed usefulness in their diagnosis, and no unresolved
high-severity secret-exposure defect. Issues, stars and usage are never
manufactured; absence of sufficient organic evidence means the gate remains
open.
