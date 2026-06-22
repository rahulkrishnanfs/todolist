---
name: codereview-to-issues
description: Read docs/codereview.md and open GitHub issues for every finding whose Status is Open, mapping defects to the bug_report template and design/enhancement gaps to the feature_request template. Use when the user wants to file/sync code-review findings as GitHub issues, or invokes this skill.
disable-model-invocation: true
---

# Code Review → GitHub Issues

Turn the **Open** findings in `docs/codereview.md` into GitHub issues — one issue per finding — using the repository's own issue templates. Defects become **bug reports**; missing-capability / design items become **feature requests**.

Issues are created with the GitHub MCP (`user-github`) tools: `search_issues` (dedup) and `issue_write` with `method: "create"`. Create them **one by one**, reporting each as it lands.

## Workflow

Copy this checklist and track progress:

```
- [ ] 1. Resolve the repo (owner/name) from the git remote
- [ ] 2. Parse docs/codereview.md and select findings with Status = Open
- [ ] 3. Classify each Open finding as bug or feature
- [ ] 4. Dedup: skip findings that already have an open issue
- [ ] 5. Preview the plan and WAIT for the user's go-ahead
- [ ] 6. Create the issues one by one, reporting each
- [ ] 7. Report a summary (created / skipped)
```

### Step 1 — Resolve the repository

Run `git remote get-url origin` and parse `owner` / `repo` (e.g. `https://github.com/rahulkrishnanfs/todolist.git` → owner `rahulkrishnanfs`, repo `todolist`). Use these for every MCP call.

### Step 2 — Parse and select findings

Read `docs/codereview.md`. It has a **Findings index** table and one per-finding detail table each headed `### <id> — <title>`. For each finding read the detail table's `Severity`, `Status`, `Lens`, `Location`, `Issue`, `Impact`, `How to reproduce`, `Suggested fix`, and any code/`go` blocks under it.

**Select only findings whose `Status` is `Open`.** Skip `Resolved` and `Partial`.

### Step 3 — Classify bug vs feature

Decide the issue type per finding:

- **Bug report** — a defect with observable wrong behavior. Signal: `How to reproduce` is concrete (a `curl`, a command, a runtime trigger) and the `Lens` is correctness/security oriented (`Programmer`, `Hacker`, `Security`, `QA`, `Observability`). Examples: `1.6`, `1.7`, `1.8`, `5.7`, `5.9`.
- **Feature request** — a missing capability or design/architecture improvement. Signal: `How to reproduce` is `n/a (design)` (or the item is "add/introduce X"). `Lens` is usually `Architect`, `API design`, or an `Ops` enhancement. Examples: `2.1`, `2.5`, `5.3`, `6.2`, `6.5`, `7.3`.

When genuinely ambiguous, default to **bug** if there is any reproducible wrong behavior, otherwise **feature**.

### Step 4 — Dedup against existing issues

For each selected finding, search for an existing open issue carrying its marker before creating:

- Marker = the finding id embedded in the title as `[codereview <id>]` (e.g. `[codereview 1.6]`).
- Use `search_issues` with a query like: `repo:<owner>/<repo> is:issue is:open in:title "[codereview 1.6]"`.
- If a match exists, **skip** that finding (note it as "already filed: #N").
- Also skip if you spot an obvious existing open issue with the same title/topic even without the marker.

### Step 5 — Preview and wait

Before creating anything, present a table of what will be created and **wait for explicit confirmation**:

```
| Finding | Type | Severity | Title | Location |
| --- | --- | --- | --- | --- |
| 1.6 | bug | Critical | Update does not reconcile path id with body id | pkg/memorystore/in_memory_todo.go:45-55 |
```

Also list anything skipped (resolved/partial/already filed). Do not proceed to Step 6 until the user says go.

### Step 6 — Create issues one by one

For each planned finding, call `issue_write` with `method: "create"`, `owner`, `repo`, a `title`, the templated `body` (below), and `labels`. Do them **sequentially**, and after each one report the created number/URL (e.g. `created #27`). If a call fails (e.g. a label doesn't exist), retry that one issue without `labels` and warn — don't abort the batch.

- **Title:** `[codereview <id>] <short title>` — e.g. `[codereview 1.6] Update does not reconcile path id with body id`.
- **Labels:** bug → `["bug"]`; feature → `["enhancement"]`.

#### Bug body (repo `bug_report.md` format)

Fill this template from the finding (keep the bold headings verbatim):

```markdown
**Describe the bug**
<Issue> <Impact>

**To Reproduce**
Steps to reproduce the behavior:
<How to reproduce — as numbered steps or the exact curl/command from the finding>

**Expected behavior**
<What should happen instead, derived from Suggested fix / Impact>

**Screenshots**
N/A (backend service).

**Desktop (please complete the following information):**
 - OS: N/A (server-side)
 - Browser N/A
 - Version N/A

**Smartphone (please complete the following information):**
 - Device: N/A
 - OS: N/A
 - Browser N/A
 - Version N/A

**Additional context**
- Severity: <Severity>
- Lens: <Lens>
- Location: `<Location>`
- Suggested fix: <Suggested fix (prose); include the code sample if short>
- Source: docs/codereview.md finding <id>
```

#### Feature body (repo `feature_request.md` format)

```markdown
**Is your feature request related to a problem? Please describe.**
<Issue> <Impact>

**Describe the solution you'd like**
<Suggested fix (prose); include the code sample if short>

**Describe alternatives you've considered**
<Any alternatives noted in the finding, else: None documented.>

**Additional context**
- Severity: <Severity>
- Lens: <Lens>
- Location: `<Location>`
- Source: docs/codereview.md finding <id>
```

### Step 7 — Report

Summarize: which issues were created (id → #number), and which findings were skipped and why (resolved/partial, or already filed as #N).

## Conventions

- Always keep the `[codereview <id>]` marker in titles — it is the dedup key for re-runs.
- One issue per finding; never batch multiple findings into one issue.
- Keep the template bold headings exactly as written so they match the repo's templates.
- Use the real `owner`/`repo` from the git remote, never hardcode.

## Guardrails

- This skill **creates GitHub issues** (a write). Only run it when the user explicitly invokes it, and only after the Step 5 preview is  confirmed.
- Never modify `docs/codereview.md` or any source file — this skill only reads the review and writes GitHub issues.
- Only file findings whose `Status` is `Open`.
- If `docs/codereview.md` has no Open findings, say so instead of creating anything.
- Don't create duplicates: honor Step 4 dedup on every run.

