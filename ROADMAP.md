# AgentExecTrace Roadmap

## Phases 0–4 — Foundation through local conformance

Completed on the v0.1 baseline: scope and evidence mapping, test-first Go CLI
implementation, the five frozen commands, repository documentation, Windows
and Linux CI configuration, and a tag-triggered release workflow. Completion
of these implementation phases does not itself satisfy a release gate.

## Phase 5.1 — v0.1 Release Gate Closure (complete)

- Close only the independently identified blockers: POSIX candidate existence
  versus executability facts and `go install` version evidence.
- Add focused adversarial regression tests; do not add a secret scanner,
  executable-format detector, shell emulator, rules engine, sixth command, or
  new platform/release matrix.
- Run formatting, vet, tests and build on one commit, then collect actual
  Windows CI, Ubuntu CI, and controlled WSL evidence.
- Synchronize the product, requirement, test, roadmap, README, changelog and
  status records with that evidence. Do not create a release tag or release.
- **Gate:** independent v0.1 Release Gate review after all required evidence is
  recorded. This phase does not self-declare release readiness.

## Phase 6 — Dogfood and external validation (after release)

- Collect only organic, reproducible troubleshooting cases from supported
  Windows/WSL agent environments.
- Assess the adoption gate in `PRODUCT_SPEC.md` after 6–8 weeks using actual
  evidence; absence of evidence leaves the gate open.
- Consider only evidence-backed v0.1.1 changes. GUI, MCP, LLM/API integration,
  telemetry, backends, auto-fix, plugins, installers, full platform doctors,
  and ARM matrix expansion remain out of scope unless separately approved.
