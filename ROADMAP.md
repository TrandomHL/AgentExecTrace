# AgentExecTrace Roadmap

## Phase 0 — Project foundation (current)

- Persist the plan, findings and progress files.
- Publish formal product, requirements, architecture, test and roadmap docs.
- Record source-policy migration status without inventing missing source files.
- **Gate:** documentation/plan gate in `task_plan.md` passes before code.

## Phase 1 — Test-first core

- Initialize the Go module and deterministic test helpers/fixtures.
- Define versioned snapshot/report structures and redaction test corpus.
- **Gate:** test skeleton proves the required failure modes before command code.

## Phase 2 — v0.1 commands

- Build `snapshot`, then `resolve`, then `probe`, `diff`, and `report --redact`.
- Keep one focused commit per coherent feature and update requirement evidence.
- **Gate:** functional requirements FR-01 through FR-07 verified locally.

## Phase 3 — Open-source readiness

- Add README with at least three real troubleshooting examples:
  1. terminal and agent have different CWD;
  2. broken PATHEXT causes `git` name resolution to fail;
  3. WSL/Windows namespace or stdout/stderr mismatch.
- Add CHANGELOG, CONTRIBUTING, issue templates, Windows/Linux CI, release
  workflow and checksums.
- **Gate:** release criteria in `PRODUCT_SPEC.md` pass; repository worktree is
  clean after the release-readiness commit.

## Phase 4 — 6–8 week adoption review

- Collect only organic, publicly verifiable or consented troubleshooting use.
- Evaluate the adoption gate in `PRODUCT_SPEC.md`; do not create fake issues,
  stars, reports or adoption metrics.
- Reprioritize only evidenced needs; keep vendor-specific repair and automated
  remediation outside v0.1 unless separately approved.

## Deferred decisions

- Project name may change without changing the v0.1 boundary.
- Choose MIT or Apache-2.0 before the first public release.
- Migrate source policy files only after `E:\project\_AD9361` becomes available.
