# AgentExecTrace Test Plan

## Test-first rule

Before implementation code, create fixtures and tests for the schema, path
classification, resolution order, process capture, semantic diff and redaction.
Every behavior below has a deterministic fixture; tests never require secrets,
network access, registry changes, or WSL configuration changes.

## Unit tests

| Area | Required cases |
|---|---|
| Model | schema version encoding/decoding, omitted optional fields, forward-safe unknown fields |
| Path namespace | Windows drive, UNC, POSIX, `/mnt/<drive>`, relative, malformed/unknown |
| PATH/PATHEXT | separator handling, empty entry, order preservation, case-insensitive Windows comparison, explicit extension |
| Resolver | selected candidate, no candidate, non-executable POSIX candidate, provenance, `.exe`/`.cmd`/`.bat` ordering |
| Probe | deterministic self-probe argv/UTF-8/stdout/stderr/exit 7; custom argv preservation, launch failure, invalid UTF-8 and truncation |
| Diff | execution namespace, CWD/path namespace, added/removed/reordered PATH, PATHEXT, resolver/probe semantic changes and golden JSON pairs |
| Redaction | bare and prefixed credential fields in nested JSON and stdout text, URL credentials, bearer/JWT-like/token/private-key values, Windows/Linux home paths, structured JSON/Markdown output and no plaintext leakage |

## Integration tests

- Invoke each CLI subcommand from a temporary directory with fixture inputs.
- Use a Go helper process to emit controlled stdout, stderr, binary bytes and
  known exit codes; do not rely on shell quoting to prove argv behavior.
- Resolve only temporary executable fixtures; do not depend on the host's Git,
  Node, Python or shell installation.
- Confirm report redaction writes a new sanitized file and retains the original.
- Confirm no command changes a designated fixture directory except its explicit
  output destination.

## Platform matrix

| Environment | Required verification |
|---|---|
| Windows GitHub Actions | Build, unit/integration tests, Windows PATH/PATHEXT and drive/UNC fixtures |
| Linux GitHub Actions | Build, unit/integration tests, POSIX fixtures |
| WSL2 | Controlled manual gate or suitable CI: `snapshot`, `resolve`, `probe`, `diff`, and `report --redact` on a known fixture. Record OS/distro/tool version and result. |

## Security test data

Fixtures use obviously fake values such as `sk_test_NOT_A_REAL_SECRET` and
synthetic JWT-shaped strings. Expected outputs contain redaction placeholders
only. No live environment variable, `.env`, local credential or production log
is ever read by the suite.

## Definition of done

All unit and integration tests pass; Windows and Linux CI are green; WSL gate is
recorded; output and error paths are covered; and an inspection of fixtures and
goldens confirms no test artifact stores unredacted secret-like values.
