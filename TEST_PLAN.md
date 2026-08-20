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
| Resolver | selected candidate, no candidate, non-executable candidate, `.exe`/`.cmd`/`.bat` ordering |
| Probe | argv preservation, stdout versus stderr separation, exit 0/non-zero, launch failure, invalid UTF-8, truncation |
| Diff | CWD/namespace change, added/removed/reordered PATH, resolver/probe semantic changes, stable no-difference result |
| Redaction | name patterns, bearer/JWT-like/token/private-key values, nested JSON/text, no plaintext in output |

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
