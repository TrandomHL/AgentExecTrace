# AgentExecTrace for AI agents

AgentExecTrace helps an AI coding agent investigate execution-context
mismatches. It is useful when a command works in the user's terminal but
fails, resolves to a different executable, starts in the wrong directory, or
behaves differently inside an agent, PowerShell, Git Bash, or WSL.

It is a diagnostic tool, not an environment repair tool. An agent should
collect evidence before changing PATH, permissions, shell configuration,
execution policy, WSL settings, or reinstalling a dependency.

## Copyable instruction snippet

Add the following to `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`, or an equivalent
agent-instructions file after installing AgentExecTrace:

```markdown
## Execution-context diagnostics

When a command works in the user's terminal but fails or behaves differently
in the AI agent, first suspect an execution-context mismatch. Before changing
PATH, permissions, shell settings, execution policy, WSL configuration, or
reinstalling tools, use AgentExecTrace to compare the contexts:

1. Run `snapshot` in both contexts, then compare the JSON with `diff`.
2. Use `resolve <command>` when executable lookup, PATH, or PATHEXT is suspect.
3. Use `probe -- <command> [args...]` when lookup succeeds but process behavior
   or stdout/stderr/exit status differs.
4. Treat output as diagnostic evidence, not automatic proof of root cause.
5. Before sharing output externally, run `report --redact` and inspect the
   result manually.
```

This snippet is vendor-neutral and does not authorize an agent to modify the
user's machine or project configuration.

## Minimal first-run prompt

After installing the executable, an agent can be given this prompt:

```text
Read the AgentExecTrace documentation at
https://github.com/TrandomHL/AgentExecTrace and use its diagnostic workflow
for this execution mismatch. First collect comparable snapshot evidence in
the terminal and agent contexts. Then use diff, resolve, or probe as the
evidence requires. Do not change PATH, permissions, shell configuration, WSL,
execution policy, or project files until the context difference is understood.
```

## Short diagnostic workflow

Use the same binary and comparable commands in each context:

```powershell
# In each context
.\agentexectrace.exe snapshot --output context.json

# In the context where the command differs
.\agentexectrace.exe resolve git --output resolve.json
.\agentexectrace.exe probe -- git --version
```

Then compare matching evidence files:

```powershell
.\agentexectrace.exe diff terminal.json agent.json
```

Interpret differences in CWD, OS/WSL identity, path namespace, PATH/PATHEXT,
candidate order, argv, stdout, stderr, exit status, and UTF-8 handling. A
matching resolver result does not prove that the environments are equivalent;
it only narrows the investigation.

## Sharing evidence safely

Before attaching or posting an artifact, create a redacted report:

```powershell
.\agentexectrace.exe report --redact --output report.md context.json
```

Redaction is defense-in-depth, not a guarantee. Inspect the generated report
for credentials, private paths, proprietary arguments, and other sensitive
data before sharing it.

## Boundaries

AgentExecTrace v0.1 does not provide MCP, plugins, an API, telemetry, cloud
collection, automatic fixes, or automatic edits to agent instruction files.
Its role is to produce bounded local evidence that a user or agent can review.
