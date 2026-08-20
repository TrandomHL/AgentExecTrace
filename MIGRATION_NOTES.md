# Inherited Policy Migration Record

## Requested source

| Source | Target | Status |
|---|---|---|
| `E:\project_AD9361\AGENTS.md` | `E:\project_Open\AGENTS.md` | Migrated on 2026-08-20 |
| `E:\project_AD9361\.codex\config.toml` | `E:\project_Open\.codex\config.toml` | Migrated on 2026-08-20 |

## Superseded initial lookup

The initially supplied source path used `E:\project\_AD9361`, which was absent.
No replacement was created. The user corrected it to `E:\project_AD9361`, where
both files existed and were read before this migration.

## Migration changes and rationale

1. Retained the source's technical-lead role, evidence-before-review rule,
   Terra/Sol escalation discipline, bounded review policy and worktree hygiene.
2. Replaced FPGA/RTL/RF/hardware-specific risk language with the equivalent
   AgentExecTrace risks: privacy, secret exposure, mutation safety, CLI/data
   compatibility and Windows/WSL behavior.
3. Adapted the authoritative-document order to existing AgentExecTrace files.
4. Preserved the source TOML model, reasoning-effort and review-model defaults;
   only the project header and path provenance changed.
5. Recorded the user-authorized initial bootstrap exception to the source's
   normal linked-worktree rule; future concurrent/risky work follows it.
