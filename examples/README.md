# Controlled diagnostic experiments

These are synthetic, isolated reproductions, **not user incidents or adoption
evidence**. They demonstrate when AgentExecTrace adds actionable evidence and
where its conclusions stop. No external service, account, Git installation or
installed AgentExecTrace is needed. Go 1.26+ is required.

From the repository root, run:

```text
go test -count=1 -v -run "^TestControlledExamples$" .
```

The test builds the production CLI, creates temporary executable fixtures, and
invokes real child processes. Each child receives its own PATH/PATHEXT and CWD;
the parent environment, registry, shell profiles and agent configuration are
unchanged. Fixtures and raw output are temporary and removed by the test
framework. Only assertions and concise outcome summaries are printed on success.
The harness has process timeouts; these do not add a timeout to the production
`probe` command. Never probe an untrusted command: the child command can have
side effects even though the diagnostic CLI does not rewrite configuration.

## What to expect

| Experiment | Observed symptom | Diagnostic findings | Restored/control check |
|---|---|---|---|
| `cwd` | The same fixture reads a relative file in the project but exits 3 in another directory. | `cwd_changed` | Restoring the child CWD restores the read and gives no snapshot differences. |
| `pathext` (Windows only) | A fixture executable launches by name with `.EXE` present, but fails with `.CPL` only. | `pathext_changed`, `command_missing` | The absolute executable still works under the broken PATHEXT; restored PATHEXT restores resolver selection. |
| `path_order` | The same command name prints `tool-A` or `tool-B`, depending on PATH order. | `path_order_changed`, `command_target_changed` | Restoring PATH restores both selected executable and observed output. |

The suite also executes the literal Windows `resolve` examples extracted from
README.md and AI_AGENT_SETUP.md, then checks the full snapshot-file -> diff-file
-> redacted-report workflow. A broken example, missing finding, wrong executable,
wrong process result or failed control makes the command exit nonzero.

Windows runs all three experiments. Linux (including WSL) runs CWD and PATH
ordering and explicitly skips Windows PATHEXT semantics. A skip is not a pass for
Windows behavior. This suite is included in `go test ./...`; use `-v -count=1`
above for a visible uncached demonstration.

## Compare with manual investigation

Without the tool, a developer can inspect each context using `Get-Location`,
`Get-Command -All`, PATH/PATHEXT inspection and explicit executable invocation.
AgentExecTrace does not replace those facts. Its contribution is to capture the
two sides in a common format, name the relevant difference, and produce a
redacted report for another person to inspect. The experiment asserts both the
named difference and a real process symptom instead of merely showing JSON.

The controlled variable is known in these experiments, and the restoration
checks establish its effect on these fixtures. In an unknown real incident,
a reported difference alone does not prove causation. The resolver describes
PATH/PATHEXT candidates, not every shell alias/function or policy decision.
These experiments do not establish adoption, time saved, token savings, full
sandbox equivalence, or a fix for any linked upstream issue.

## Using this as public evidence

Link the source (`controlled_examples_test.go`), the exact commit tested, and the
corresponding CI run. State the OS and which cases passed or skipped. Do not
present controlled reproductions as independent user feedback. No raw local
snapshot or probe output is needed for this demonstration.

## Recorded local validation — 2026-09-17

| Environment | Controlled examples | Full tests / vet / build |
|---|---|---|
| Windows amd64, Go 1.26.7 | CWD, PATHEXT, PATH order, documented commands and report workflow: PASS | PASS |
| Ubuntu under WSL, Go 1.26.0 | CWD, PATH order, documented commands and report workflow: PASS; Windows PATHEXT: SKIP | PASS |

These are local execution results for this change. A new hosted GitHub Actions
run has not been performed for it. No external-user validation or performance
measurement is claimed.
