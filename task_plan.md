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

Phase 5.1 — v0.1 Release Gate Closure (in progress on
`agent/phase-5-1-release-gate-closure`, based on `8948c2a`).

## Next step

Close the independently identified release blockers without expanding the
five-command public surface. Remote publication and release tagging are not
authorized by this phase.

## Phase 5.1 plan (Superpowers, simplified by Ponytail)

1. Add adversarial regression coverage before changing redaction, resolver, or
   diff behavior.
2. Redact bare credential-like field names in structured JSON and text using
   the existing name-based redactor; do not add a generic secret scanner.
3. Preserve POSIX no-extension candidate provenance as `unknown`; execute-bit
   detection remains the only selection criterion.
4. Emit `command_provenance_changed` when resolved executable paths are the
   same but their recorded provenance differs. Keep probe comparison limited to
   already-recorded result metadata.
5. Correct release asset matching and require expected assets to exist.
6. Run `gofmt`, `go vet ./...`, `go test ./...`, and `go build ./...` against
   one commit; record Windows/Ubuntu CI and a controlled WSL run as evidence.
7. Synchronize product, requirement, test, roadmap, status, README and
   changelog documents only with results actually observed.

### Explicit exclusions

No sixth command, ELF/shebang inspection, shell emulation, generic rules
engine, entropy/DLP scanner, agent integration, GUI, telemetry, cloud service,
auto-fix, plugin, installer, or ARM release matrix.

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
- [x] Go module, test skeleton and deterministic fixtures exist before command code.
- [x] Unit and integration tests pass on the host platform.

### Release gate

- [x] Windows and Linux CI configured and passing.
- [x] WSL validation recorded (controlled Ubuntu manual gate).
- [x] README, CHANGELOG, CONTRIBUTING and issue templates present.
- [x] Release workflow builds archives and checksums.
- [ ] `git status --short` is empty after the final evidence commit.

### Phase 5 conformance gate

- [x] Real GitHub module identity and all internal imports are aligned.
- [x] FR-05 semantic diff has explicit, golden-backed findings for namespace,
  CWD, PATH/PATHEXT, resolver and applicable probe evidence.
- [x] `probe` has a deterministic no-argument self-probe and retains custom
  command probing.
- [x] `report --redact` produces a schema-aware, Issue-ready report with secret
  and home-profile redaction summaries.
- [x] Resolver reports minimal provenance and applies POSIX executable checks.
- [x] Requirements compliance matrix has execution evidence for FR-01 through
  FR-07; CI enforces gofmt, vet, test and build on Windows and Ubuntu.
- [x] README, CHANGELOG and test plan match verified behavior; controlled WSL
  validation is recorded.

## Phases

### Phase 0: Repository and inherited policy

- [x] Inspect `E:\project_Open` before writing.
- [x] Initialize only because the directory was empty.
- [x] Inspect the requested migration source.
- [x] Migrate `AGENTS.md` and `.codex/config.toml` from corrected source with documented project adaptations.
- [x] Record the current source-unavailable constraint.
- **Status:** complete with a deferred migration dependency.

### Phase 1: Plan and formal specification

- [x] Persist Superpowers plan and Ponytail simplification decision.
- [x] Research and record real public issue evidence.
- [x] Land and review PRODUCT_SPEC, REQUIREMENTS, ARCHITECTURE, TEST_PLAN and ROADMAP.
- [x] Pass documentation gate.
- **Status:** complete

### Phase 2: Test-first foundation

- [x] Create Go module, test harness and checked-in fixtures.
- [x] Add redaction, path and process-result unit tests before command implementation.
- **Status:** complete

### Phase 3: v0.1 implementation

- [x] Implement `snapshot`, `resolve`, `probe`, `diff`, and `report --redact`.
- [x] Update evidence-to-requirement map and user documentation as each command lands.
- **Status:** complete

### Phase 4: Verification and release readiness

- [x] Complete unit/integration tests, Windows build, Linux cross-build and controlled WSL gate.
- [x] Add release automation, repository contribution files and troubleshooting examples.
- [x] Commit the verified implementation; leave the worktree clean.
- **Status:** complete locally; remote CI is pending repository publication.

### Phase 5: v0.1 conformance closure

- Correct the module identity, then replace raw-field diffing with the small,
  explicit comparisons required by FR-05.
- Add a hidden same-binary probe helper so no-argument `probe` is deterministic
  on Windows and WSL while preserving `probe -- <command> [args...]`.
- Make redacted reports schema-aware and Markdown-friendly; add only the
  redactions and summary needed for safe issue sharing.
- Add reliable resolver provenance and POSIX execute-permission behavior; avoid
  speculative alias, reparse-point or shell discovery logic.
- Add test fixtures before or with each behavior, align docs and CI, execute
  local/controlled WSL gates, and record only actual evidence.
- **Status:** complete; superseded by Phase 5.1 independent release-gate closure.

### Phase 5.1: v0.1 Release Gate Closure

- [x] Redact bare credential fields and corresponding structured/text cases.
- [x] Keep POSIX no-extension provenance as `unknown`.
- [x] Report same-target resolver provenance changes semantically.
- [x] Require all release assets through multiline action globs.
- [x] Pass local gofmt, vet, test and build, GitHub Windows/Ubuntu CI, and
  controlled Ubuntu WSL validation for the implementation commit.
- [ ] Complete independent v0.1 Release Gate review; do not create tag/release
  in this phase.
- **Status:** implementation and evidence complete; awaiting independent review.

## Phase 5 implementation plan (Superpowers)

1. Freeze the five public commands and record the baseline/remote CI evidence.
2. Change module identity and imports, then validate the module builds locally.
3. Extend the existing data model only where snapshot/report evidence requires
   it; implement explicit semantic diff comparisons rather than a rule engine.
4. Implement self-probe through a hidden helper mode and test process facts.
5. Add structured redacted-report rendering, lightweight provenance, and POSIX
   executable checks with deterministic fixtures.
6. Bring documentation, compliance evidence and CI checks into line, then run
   all prescribed local and controlled platform validation before review.

## Phase 5 Ponytail review

- Keep the existing single binary, package layout and standard library.
- Do not add a sixth command, helper executable, shell dependency, generic rule
  engine, plugin system, high-entropy scanner, alias parser or WSL CI service.
- Prefer small comparison/rendering functions and checked-in fixtures; retain
  uncertain provenance as `unknown`.

## Decisions

| Decision | Rationale |
|---|---|
| Go single executable | User requirement; simple portable distribution with no runtime dependency. |
| Go standard library only for v0.1 | Meets no-network/no-runtime goal and avoids supply-chain and maintenance cost. |
| JSON snapshots and reports | Portable, inspectable, diffable, and supported by Go stdlib. |
| Opt-in environment collection | Privacy by default; default output never dumps all environment variables. |
| No source-file substitutes | No substitute was created while the initial path was absent; corrected-source migration is now complete and documented. |
| MIT license | A permissive target was allowed by the original Phase 0 draft; MIT is the simplest v0.1 choice. |

## Errors and constraints

| Item | Attempt | Resolution |
|---|---:|---|
| Initial path typo created empty `E:\project\_Open` repository | 1 | Stopped immediately when corrected; no product files or commits were created there. |
| `E:\project\_AD9361` absent | 1 | Defer migration; record the exact paths and do not fabricate policy content. |
| Go unavailable on PATH | 1 | Inspect local standard installation locations; request permission only if installation is required. |
| Local Go-location check had an invalid PowerShell pipeline | 1 | Retry with a grouped `foreach` expression; no system or repository state changed. |
| Go installer completed download output but no `go.exe` existed at the standard install path | 1 | Verify package-manager installation state and standard locations before choosing a different installer path. |
| GitHub template/workflow parent directories did not yet exist | 1 | Created only the required empty parent directories, then added repository metadata. |
| GitHub template/workflow parent directories did not yet exist | 1 | Create empty `.github` parent directories, then apply repository metadata files. |
