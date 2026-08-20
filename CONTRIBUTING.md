# Contributing

- Keep v0.1 within the frozen scope in `PRODUCT_SPEC.md`.
- Do not add a runtime dependency, service, telemetry, network call, plugin
  system or configuration mutation without approved scope change.
- Never commit secrets, raw environment dumps, personal reports or build files.

Before a pull request, run:

```powershell
gofmt -w .
go test ./...
go build ./...
```

Add deterministic fixtures and tests before a new behavior. For WSL-affecting
changes, record controlled WSL validation in the pull request. Issues must use
safe reproduction steps and redacted artifacts only.
