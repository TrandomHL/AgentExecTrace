# AgentExecTrace v0.1 Requirements

## Requirement conventions

`MUST` is required for v0.1. Public issues are problem evidence, not a promise
that a given vendor will adopt or endorse this tool.

## Evidence mapping

| Evidence ID | Real public issue | Requirement supported |
|---|---|---|
| E-01 | [OpenAI Codex #20858](https://github.com/openai/codex/issues/20858) — agent CWD differs from the integrated terminal on Windows | Snapshot CWD and semantic context comparison |
| E-02 | [OpenAI Codex #15174](https://github.com/openai/codex/issues/15174) — truncated PATHEXT prevents bare executable resolution | PATH/PATHEXT collection and resolver explanation |
| E-03 | [OpenAI Codex #19171](https://github.com/openai/codex/issues/19171) — PATH/Path environment construction issue | Case-insensitive Windows environment handling |
| E-04 | [Claude Code #19653](https://github.com/anthropics/claude-code/issues/19653) — WSL invocation executes in Windows/MINGW and uses UNC paths | Windows/WSL identity and namespace classification |
| E-05 | [Claude Code #26430](https://github.com/anthropics/claude-code/issues/26430) — command runs but stdout/stderr pipe capture is broken | Probe preserves stdout, stderr and exit status independently |
| E-06 | [OpenAI Codex #34233](https://github.com/openai/codex/issues/34233) — credentials exposed in stored tool output | Opt-in environment collection and shareable redaction |

## Functional requirements

| ID | Requirement | Evidence | Verification |
|---|---|---|---|
| FR-01 | `snapshot` MUST emit versioned structured context evidence without a default full environment dump. | E-01, E-03, E-04, E-06 | Unit schema and integration fixture tests |
| FR-02 | `snapshot` MUST classify Windows drive, UNC, WSL Linux, mounted Windows, and unknown paths. | E-01, E-04 | Table-driven classifier tests |
| FR-03 | `resolve` MUST inspect PATH and, on Windows, PATHEXT, then report candidate order, selected candidate and unresolved reason. | E-02, E-03 | Deterministic temporary-directory fixtures |
| FR-04 | `probe` MUST record provided argv; separately capture bounded stdout and stderr; record exit code/launch error; and mark invalid UTF-8 rather than corrupting it. | E-05 | Helper-process integration tests |
| FR-05 | `diff` MUST compare snapshot fields semantically, identify added/removed/reordered PATH entries and report namespace, CWD, executable-resolution and probe-result changes. | E-01, E-04 | Golden JSON-pair tests |
| FR-06 | `report --redact` MUST redact bare and prefixed common secret-bearing names, URL credentials, and recognized token/private-key patterns from structured data and text. | E-06 | Adversarial secret fixture tests with no plaintext expected output |
| FR-07 | Commands MUST be read-only and MUST NOT modify PATH, registry, WSL settings, execution policy, user configuration or project configuration. | Product safety boundary | Integration test plus code review checklist |

## Non-functional requirements

| ID | Requirement |
|---|---|
| NFR-01 | Build as one Go executable with no Node/Python runtime, service, database or required network access. |
| NFR-02 | Prefer Go standard library; any dependency needs a documented v0.1 necessity review. |
| NFR-03 | Outputs are deterministic where machine state permits and include schema/tool versions. |
| NFR-04 | Capture sizes are bounded and truncation is explicit. |
| NFR-05 | Linux CI checks portable behavior; Windows CI checks Windows semantics; WSL is validated in controlled CI or a recorded manual gate. |

## v0.1 compliance matrix — Phase 5 baseline

| Requirement | Implementation | Test | Status | Evidence |
|---|---|---|---|---|
| FR-01 | Versioned `model.Snapshot`; no environment dump | `TestSnapshotDoesNotDumpEnvironmentValues` | PASS | local gate and GitHub CI run `32364080361` |
| FR-02 | `pathinfo.Classify` | `TestClassify` | PASS | local gate and GitHub CI run `32364080361` |
| FR-03 | PATH/PATHEXT resolver, Windows case handling and provenance | `internal/resolve` tests | PASS | local gate and GitHub CI run `32364080361` |
| FR-04 | Bounded custom probe and same-binary self-probe | `TestProbeCommandCapturesResult`, `TestProbeWithoutArgumentsRunsSelfProbe` | PASS | local gate and GitHub CI run `32364080361` |
| FR-05 | Explicit snapshot/PATH/PATHEXT/resolve/probe comparisons | `internal/diff` semantic and golden tests | PASS | local gate and GitHub CI run `32364080361` |
| FR-06 | Structured JSON/text report redaction, summaries and home paths | `internal/redact` and report tests | PASS | local gate and GitHub CI run `32364080361` |
| FR-07 | Read-only command implementation; no configuration mutation | code review plus command integration tests | PASS | no mutation APIs or writes outside explicit `--output` paths inspected 2026-08-20 |

Windows and Ubuntu CI passed for the Phase 5.1 implementation commit in
GitHub Actions run `32364080361`; the final evidence commit also passed in run
`32364353630`. Controlled Ubuntu WSL validation passed all five public
commands. Any later behavior-changing commit must obtain its own validation
evidence before independent review.

## Out of scope

The v0.1 non-goals in `PRODUCT_SPEC.md` are binding requirements. In particular,
AgentExecTrace does not diagnose every Linux/macOS condition or repair an
installation for any particular coding agent.
