# AgentExecTrace — Persistent Delivery Plan

## Goal

Deliver AgentExecTrace v0.1: a vendor-neutral, privacy-first Go CLI that gathers
and compares execution-context evidence for AI coding agents on Windows 11 and
WSL2, without changing the user's machine or configuration.

## Planning lineage

- **Superpowers plan:** establish repository policy, freeze the scope and
  evidence mapping, create tests and fixtures before implementation, then build
  command-by-command with verification gates.
- **Ponytail review:** use one Go executable and the standard library; keep a
  small package layout and JSON files. Do not add a server, database, Node or
  Python runtime, plugins, DI, rules engine, telemetry, auto-fix, or platform
  configuration mutation.

## Current phase

Phase 2 — test-first foundation (in progress).

## Next step

Create deterministic Go test fixtures and failing tests for the frozen v0.1
behaviors before adding command implementation.

## Gates

### Planning/documentation gate (must pass before code)

- [x] Correct repository inspected and initialized without replacing user files.
- [x] Persistent planning files created in the repository root.
- [x] Five formal Phase 0 documents landed and cross-reviewed.
- [x] Each v0.1 command mapped to a real public issue.
- [x] Migration status recorded; no invented replacement for unavailable source rules.
- [x] Scope, privacy, non-goals, v0.1 acceptance gate and 6–8 week adoption gate explicit.

### Implementation gate

- [x] Documentation gate passed.
- [ ] Go module, test skeleton and deterministic fixtures exist before command code.
- [ ] Unit and integration tests pass on the host platform.

### Release gate

- [ ] Windows and Linux CI configured and passing.
- [ ] WSL validation recorded (CI or controlled manual run).
- [ ] README, CHANGELOG, CONTRIBUTING and issue templates present.
- [ ] Release workflow builds archives and checksums.
- [ ] `git status --short` is empty after the phase commit.

## Phases

### Phase 0: Repository and inherited policy

- [x] Inspect `E:\project_Open` before writing.
- [x] Initialize only because the directory was empty.
- [x] Inspect the requested migration source.
- [ ] Copy `AGENTS.md` and `.codex/config.toml` verbatim-plus-minimal-project-additions when the source becomes available.
- [x] Record the current source-unavailable constraint.
- **Status:** complete with a deferred migration dependency.

### Phase 1: Plan and formal specification

- [x] Persist Superpowers plan and Ponytail simplification decision.
- [x] Research and record real public issue evidence.
- [x] Land and review PRODUCT_SPEC, REQUIREMENTS, ARCHITECTURE, TEST_PLAN and ROADMAP.
- [x] Pass documentation gate.
- **Status:** complete

### Phase 2: Test-first foundation

- [ ] Create Go module, test harness and checked-in fixtures.
- [ ] Add redaction, path and process-result unit tests before command implementation.
- **Status:** in_progress

### Phase 3: v0.1 implementation

- [ ] Implement `snapshot`, `resolve`, `probe`, `diff`, and `report --redact`.
- [ ] Update evidence-to-requirement map and user documentation as each command lands.
- **Status:** pending

### Phase 4: Verification and release readiness

- [ ] Complete unit, integration, Windows, Linux and WSL gates.
- [ ] Add release automation, repository contribution files and troubleshooting examples.
- [ ] Commit each coherent completed phase; leave the worktree clean.
- **Status:** pending

## Decisions

| Decision | Rationale |
|---|---|
| Go single executable | User requirement; simple portable distribution with no runtime dependency. |
| Go standard library only for v0.1 | Meets no-network/no-runtime goal and avoids supply-chain and maintenance cost. |
| JSON snapshots and reports | Portable, inspectable, diffable, and supported by Go stdlib. |
| Opt-in environment collection | Privacy by default; default output never dumps all environment variables. |
| No source-file substitutes | The requested inherited policies cannot be faithfully preserved while their source directory is absent. |

## Errors and constraints

| Item | Attempt | Resolution |
|---|---:|---|
| Initial path typo created empty `E:\project\_Open` repository | 1 | Stopped immediately when corrected; no product files or commits were created there. |
| `E:\project\_AD9361` absent | 1 | Defer migration; record the exact paths and do not fabricate policy content. |
