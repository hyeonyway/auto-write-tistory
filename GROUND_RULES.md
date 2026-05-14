# DevLog Studio Ground Rules

## Purpose

This document defines the working rules for humans and coding agents working on DevLog Studio.
Follow it whether the implementation is done manually, with Codex, with Claude, or with another assistant.

## Branching And Merge Policy

- Do not commit directly to `main`.
- Create a feature branch for each unit of work.
- Keep commits focused and reviewable while working on the branch.
- Before merging, rebase the feature branch on top of the latest `main`.
- Resolve conflicts during the rebase on the feature branch.
- After verification passes, merge into `main` with a squash merge.
- The final squash commit message should describe the completed feature or fix, not every intermediate attempt.
- Issue creation is optional and not required for this project.

Recommended flow:

```bash
git switch main
git pull --rebase
git switch -c feature/<short-topic>

# work and commit locally

git fetch origin
git rebase origin/main
# resolve conflicts if needed
# run verification
# squash merge through the chosen Git UI or hosting provider
```

## Commit Rules

- Commit only files related to the current task.
- Do not include secrets, `.env`, local DB data, screenshots with credentials, or generated dependency folders.
- Prefer small commits while developing, then squash merge to `main`.
- Use clear commit messages:
    - `feat: add algorithm preview API`
    - `fix: handle tistory login failure`
    - `refactor: split tistory publisher adapter`
    - `docs: document selenium publish constraints`
    - `chore: add pre-commit hooks`
- Do not rewrite or revert another contributor's work unless explicitly asked.

## Formatting

- Use spaces, not tabs.
- Indentation is 4 spaces across the repository unless a tool-specific format requires otherwise.
- Keep line endings as LF.
- Keep files UTF-8 encoded.
- Run formatting before committing.

The source of truth for editor behavior is `.editorconfig`.
The source of truth for frontend formatting is `.prettierrc`.

## Pre-Commit

Install pre-commit once:

```bash
pre-commit install
```

Run all hooks manually:

```bash
pre-commit run --all-files
```

The hooks check common file hygiene and run Prettier for frontend-supported file types.
Install frontend dependencies before running the Prettier hook:

```bash
cd frontend
npm install
```

## Verification Before Merge

Before squash merging to `main`, run the relevant checks for the files changed.

Minimum checks for frontend changes:

```bash
cd frontend
npm run format:check
npm run build
```

Minimum checks for docs/config-only changes:

```bash
pre-commit run --all-files
```

If a check cannot be run, document why in the final handoff or pull request description.

## Tistory Automation Safety

- Tistory/Kakao login credentials must not be committed.
- The app receives Kakao login ID/PW from the web publish modal at publish time only.
- The frontend must not store credentials in localStorage, sessionStorage, IndexedDB, or files.
- The backend must not store credentials in PostgreSQL, logs, API responses, or screenshots.
- Selenium publishing is expected to work only on a local environment where Kakao additional verification can be completed.
