# Changelog

This project follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and Semantic Versioning.

## [Unreleased]

### Fixed

- Corrected `resolve --output file name` examples in the README and agent setup.
- Made the first-use workflow create both snapshots, use stable executable paths
  and filenames, compare them, and produce a separate redacted report.
- Clarified release/source installation command names and snapshot privacy wording.

### Added

- Three executable controlled reproductions: CWD-dependent file access, Windows
  PATHEXT lookup failure, and PATH precedence. Real CLI/process checks include
  restoration controls, documentation-command regressions and the report workflow.
- A reproducibility guide explaining the comparison with manual checks and the
  limits of these synthetic examples; no adoption or performance claims.

## [0.1.0] - 2026-08-23

### Added

- v0.1 baseline: snapshot, executable resolution, process probe, semantic
  snapshot diff and redacted report commands.
- Windows/WSL path classification, PATH/PATHEXT evidence, bounded output and
  common-secret redaction tests.

### Changed

- Corrected the Go module path to `github.com/TrandomHL/AgentExecTrace`.
- `diff` now reports explicit semantic findings for context, PATH/PATHEXT,
  resolver and probe evidence instead of raw slice changes.
- `probe` supports a deterministic no-argument self-probe while retaining
  custom command probing.
- `report --redact` renders shareable Markdown and summarizes secret, token and
  home-path redactions; it now covers bare credential fields, URL credentials,
  and adversarial nested/text cases.
- POSIX execute permission remains the selection check, but a no-extension
  candidate now reports provenance as `unknown`; `diff` reports a provenance
  change when the selected executable remains the same.
- Resolver candidates now report existence separately from executability, so a
  non-executable POSIX file is not reported as absent.
- The `go install` instructions pin the intended v0.1.0 source version and
  distinguish source-build from release-workflow provenance.
- Release assets use a multiline file list and fail when an expected file is
  absent.
