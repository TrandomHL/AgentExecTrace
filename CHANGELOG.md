# Changelog

This project follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and Semantic Versioning.

## [Unreleased]

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
  home-path redactions; resolver reports basic provenance and honors POSIX
  execute permissions.
