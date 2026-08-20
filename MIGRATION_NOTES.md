# Inherited Policy Migration Record

## Requested source

| Source | Target | Status |
|---|---|---|
| `E:\project\_AD9361\AGENTS.md` | `E:\project_Open\AGENTS.md` | Deferred: source path absent on 2026-08-20 |
| `E:\project\_AD9361\.codex\config.toml` | `E:\project_Open\.codex\config.toml` | Deferred: source path absent on 2026-08-20 |

## Reason for no target files

The request requires retaining every general rule from the source. The source
directory and both requested files were absent when checked. Creating a guessed
`AGENTS.md` or `config.toml` would falsely claim that general clauses were
preserved and could silently weaken a user policy. Therefore no replacement,
truncated copy, or project-local approximation has been created.

## Safe migration procedure once source is available

1. Copy each source file without deleting or rewording any general clause.
2. Compare source and target; retain source order unless a syntax constraint
   requires a documented relocation.
3. Append only these AgentExecTrace-specific constraints: Go-only v0.1,
   standard-library-first, read-only diagnostics, privacy-by-default, no secret
   output, and no machine/configuration mutation.
4. Update this record with a line-by-line change summary and rationale.
5. Validate TOML syntax and inspect `git diff`; commit the migration separately.
