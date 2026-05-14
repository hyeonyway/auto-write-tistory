# Agent Instructions

Follow `GROUND_RULES.md` first. This file adds agent-specific operating rules.

## Scope Control

- Read the relevant design document before changing implementation:
    - `devlog_studio_overall_design.md`
    - `docs/phase1/initial_design.md`
    - `docs/phase1/known_issues.md`
- Keep edits scoped to the user's current request.
- Do not introduce unrelated refactors.
- Do not create GitHub issues unless the user explicitly asks.

## Phase Documentation

- Keep phase documents under `docs/phaseN/`.
- Each phase must keep these two baseline documents:
    - `initial_design.md`: first design, scope, architecture, APIs, data model, and completion criteria.
    - `known_issues.md`: confirmed problems, root cause, resolution direction, and current status.
- When a known issue is fixed, update `known_issues.md` with the implementation status instead of deleting the history.

## Git Discipline

- Work on a feature branch.
- Do not commit directly to `main`.
- Rebase the feature branch on `main` before merge.
- Resolve conflicts on the feature branch.
- Use squash merge into `main`.
- Do not run destructive git commands such as `git reset --hard` or `git checkout -- <file>` unless explicitly asked.

## Formatting And Checks

- Use 4-space indentation.
- Respect `.editorconfig`.
- Use `.prettierrc` for frontend-supported files.
- Run pre-commit checks before claiming work is ready when possible.
- For frontend work, run:

```bash
cd frontend
npm run format:check
npm run build
```

## Secrets

- Never write real credentials into tracked files.
- Never log Kakao/Tistory ID or password.
- Never store Kakao/Tistory ID or password in browser storage or PostgreSQL.
- Use placeholder values in docs and examples.
