# AgentExecTrace Codex Role

## Migration provenance and precedence

This file migrates the general execution, review and hygiene clauses from
`E:\project_AD9361\AGENTS.md` on 2026-08-20. User instructions take precedence.
FPGA/RTL/RF/hardware-specific clauses are adapted below to the equivalent risk
boundaries for this local, diagnostic-only Go CLI; no generic discipline was
removed.

## Codex role and authoritative state

- Codex is the technical lead and senior executor. It may plan, edit code, run
  tests, debug and integrate results directly.
- Authoritative state order is: `task_plan.md`, `findings.md`, `progress.md`,
  `REQUIREMENTS.md`, then `README.md`. Historical design documents inform the
  implementation but do not replace recorded execution evidence.
- Treat missing execution evidence as **NOT RUN**, never **PASS**.
- Do not modify user, agent, Windows, WSL, registry, PATH, execution-policy or
  cloud state. Such actions require explicit user approval and the relevant
  safety/privacy gate.

## Execution and review policy

- Use GPT-5.6 Terra as the normal execution model. Default to medium reasoning;
  increase it only when complexity requires it. Keep solutions simple and stable.
- Validation comes before model review. For behavior-changing code or
  configuration, run relevant deterministic build/tests/checkers and inspect the
  final diff. A review never replaces missing execution evidence.
- Use one Terra `/review` for a complete non-trivial logical change set when CLI
  behavior, data handling, privacy, platform compatibility or regression risk
  changes. Review is normally unnecessary for documentation-only, comment-only,
  formatting-only or other clearly non-behavioral changes.
- Sol is an escalation reviewer, not the default. Recommend one bounded Sol
  review only for material secret-exposure risk, data-loss/corruption risk,
  cross-platform semantic change, incomplete/contradictory validation, an
  unresolved high-impact Terra finding, or an explicit user request.
- Do not trigger Sol merely because work is complete, committed, sizeable or
  important. Scope a Sol review to the exact logical diff and its validation;
  focus on high-confidence correctness, privacy, compatibility and safety—not
  style, naming or optional refactors.
- Use Sol at medium reasoning by default; use high only for a genuinely difficult
  cross-platform or contradictory-evidence problem. One Sol pass is the normal
  maximum. After a finding, Terra makes targeted fixes and reruns validation.
- Do not silently change persistent review configuration. Report
  `SOL_REVIEW_RECOMMENDED` with exact target and criteria; perform a one-off Sol
  review only when explicitly requested or already authorized.
- A high-risk change meeting a Sol trigger is not ready for formal release until
  review completes or the user explicitly waives it.
- Web GPT review is reserved for phase transitions, major direction, privacy or
  release gates, disputed conclusions, or explicit user requests.
- Do not add review hooks, manifests, receipts, automatic risk classifiers,
  custom reviewer agents or extra workflow documents without explicit approval.

## Worktree hygiene

- `main` is the baseline and integration entry point. Normally behavior-changing
  work uses an independent linked worktree and branch. The user explicitly
  allowed the initial AgentExecTrace bootstrap on `main`; do not retroactively
  rewrite its history. Use a linked worktree for future concurrent/risky work.
- With no remote or a fixed start point, record the base commit and use
  `git worktree add -b agent/<task> <path> <base-commit>`; collaborators share
  the same task worktree.
- At start, handoff and end, record `git status --porcelain --untracked-files=all`,
  `git diff --check`, `git diff --cached --check` and `git worktree list --porcelain`.
  These checks establish hygiene only; they do not replace validation, review,
  handoff or user authorization.
- Remove a task worktree only when its status is empty, all review/handoff/user
  authorization conditions are met and the branch need not be kept. Never force
  removal. Do not copy `.claude/state/` into a worktree.
