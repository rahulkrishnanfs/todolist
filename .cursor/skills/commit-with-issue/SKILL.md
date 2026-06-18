---
name: commit-with-issue
description: Create an industry-standard Conventional Commit tied to a GitHub issue and a matching feature branch. Use when the user wants to commit changes, asks for help writing a commit message, or invokes this skill. Always prompts for the GitHub issue number, reuses the issue's existing branch or creates one named feature/#[issue]-[branchname], and writes a typed commit (feat, fix, docs, etc.) whose subject ends with the issue number in parentheses (e.g. (#42)).
disable-model-invocation: true
---

# Commit with Issue

Create a Conventional Commit for the current changes, on a dedicated feature branch, always linked to a GitHub issue.

## Workflow

Copy this checklist and track progress:

```
- [ ] 1. Inspect the changes (git status + diff)
- [ ] 2. Get the GitHub issue number (REQUIRED — ask if not provided)
- [ ] 3. Pick the commit type and a short subject
- [ ] 4. Reuse the issue's branch if it exists, else create feature/#[issue]-[branchname]
- [ ] 5. Stage and commit using the Conventional Commit format
- [ ] 6. Show the result (git status + git log -1)
```

### Step 1 — Inspect the changes

Run `git status` and `git diff` (and `git diff --cached`) to understand what changed. Use this to choose the right type, scope, and subject. Do not commit files that may contain secrets (`.env`, credentials, keys) — warn the user instead.

### Step 2 — Get the issue number (required)

The commit MUST reference a GitHub issue. If the user did not provide one, ask:

> "What is the GitHub issue number this commit relates to?"

Do not proceed past this step without a number. Accept either `42` or `#42`. Always write the issue with a leading `#` in both the branch name and at the end of the commit subject (e.g. `#42`).

### Step 3 — Choose the commit type

Pick the single type that best describes the change:

| Type | Use for |
|------|---------|
| `feat` | A new feature or capability |
| `fix` | A bug fix |
| `docs` | Documentation only (README, comments, docs/) |
| `style` | Formatting/whitespace, no logic change |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `perf` | A performance improvement |
| `test` | Adding or fixing tests |
| `build` | Build system, Dockerfile, or dependencies |
| `ci` | CI configuration (workflows, pipelines) |
| `chore` | Maintenance, tooling, or housekeeping |
| `revert` | Reverting a previous commit |

Add an optional `(scope)` when one component is clearly affected (e.g. `feat(auth)`, `fix(store)`). If the change spans many areas, omit the scope.

### Step 4 — Reuse or create the feature branch

Each issue gets **one** branch — never open a second branch for the same issue. Before creating anything, check whether a branch for this issue already exists:

```bash
git branch --list "feature/#[issue]-*"
```

- If a matching branch exists (or you are already on it), **reuse it**: `git switch <existing-branch>`. Commit the new change there.
- Otherwise, create one named `feature/#[issue number]-[branchname]`:
  - Keep the `#` prefix on the issue number.
  - `[branchname]` is a kebab-case slug derived from the commit subject: lowercase, spaces → hyphens, drop punctuation, keep it short (3–5 words).
  - Example: issue `42`, subject "add JWT authentication" → `feature/#42-add-jwt-authentication`.

  ```bash
  git switch -c feature/#42-add-jwt-authentication
  ```

### Step 5 — Stage and commit

Stage the relevant files (`git add <paths>`), then commit using a HEREDOC so the body is formatted correctly:

```bash
git commit -m "$(cat <<'EOF'
type(scope): short imperative subject (#42)

Optional body explaining the WHY of the change, wrapped at ~72 chars.
EOF
)"
```

Commit message rules:
- Append the issue number at the end of the subject line in parentheses: `(#42)`.
- Subject ≤ ~72 chars including the trailing `(#42)`, imperative mood ("add", not "added"/"adds"), no trailing period after the issue reference.
- Keep the subject focused on a single logical change.

### Step 6 — Verify

Run `git status` and `git log -1` to confirm the branch and commit landed as intended, then report the branch name and commit subject back to the user.

## Examples

**Example 1** — new feature, issue 42

Branch: `feature/#42-add-jwt-authentication`

```
feat(auth): add JWT-based login endpoint (#42)

Add a /login route that issues signed JWTs and middleware that
validates them on protected routes.
```

**Example 2** — bug fix, issue 17

Branch: `feature/#17-fix-empty-store-500`

```
fix(store): return empty list when store is empty (#17)

GetAll returned an error for an empty store, which surfaced as a 500.
Return an empty slice so the API responds 200 with [].
```

**Example 3** — docs only, issue 23

Branch: `feature/#23-update-readme-endpoints`

```
docs(readme): correct REST endpoints and verbs (#23)
```

## Guardrails

- Never update git config, and never use `--no-verify` or force-push unless the user explicitly asks.
- Never create a commit without an issue number.
- Only commit when the user has asked to commit; this skill performs the commit as its purpose.
- If there are no changes to commit, say so instead of creating an empty commit.
