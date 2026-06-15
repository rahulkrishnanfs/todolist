# Code Review Update

Refresh `docs/codereview.md` so it reflects the current state of the codebase, and present every finding in a **table view** that maps directly onto a JIRA / GitHub issue.

> **Scope (hard rule):** The ONLY file this skill may create or modify is `docs/codereview.md`. It may *read* any source, infrastructure, or other file to perform the review, but it must never write to anything else — not source code, not infrastructure, and not any other file under `docs/` (e.g. `docs/note.txt`, `docs/https_status_codes.md`). If a change seems to require editing another file, stop and report it instead.

Review the code through four lenses on every pass:

- **Software architect** — boundaries and layering (ports/adapters), missing service layer, dependency direction, `context.Context` propagation, configuration, statelessness/scalability, API/REST design.
- **Senior expert programmer** — correctness bugs, error handling, concurrency/locking, idiomatic Go, naming, allocations, dead code, logging quality, test coverage.
- **Hacker / attacker** — mass assignment, IDOR/ownership, injection, unbounded input and resource-exhaustion DoS, information disclosure via error messages, race/TOCTOU, path traversal.
- **Security expert** — authentication/authorization, input validation, secrets handling, server timeouts (slowloris), TLS/headers/CORS/rate limiting, dependency and container hardening, logging of sensitive data.

## Output format (table view)

The document must be issue-tracker friendly. Use two layers of tables:

1. **Findings index** — one master table near the top for triage, with one row per finding:

   ```
   | ID | Severity | Status | Lens | Issue | Location |
   | --- | --- | --- | --- | --- | --- |
   | 1.6 | Critical | Open | Hacker | `Update` does not reconcile path id with body id | `memorystore/in_memory_todo.go:45-55` |
   ```

2. **Per-finding detail table** — under each item heading, one vertical `Field | Detail` table holding exactly the fields a developer needs to file a ticket:

   ```
   ### 1.6 — `Update` does not reconcile path id with body id

   | Field | Detail |
   | --- | --- |
   | **Severity** | Critical |
   | **Status** | Open |
   | **Lens** | Hacker / Correctness |
   | **Location** | `memorystore/in_memory_todo.go:45-55` |
   | **Issue** | One-paragraph description of the defect/risk. |
   | **Impact** | What breaks / who is affected. |
   | **How to reproduce** | Concrete, copy-pasteable steps (e.g. a `curl` call + expected vs actual). For non-reproducible/design items, describe the trigger or write `n/a (design)`. For resolved items write `n/a (resolved)`. |
   | **Suggested fix** | Short prose of the fix; if there's a code sample, end with "See sample below." |
   ```

   Immediately after the table, when useful:
   - a fenced ``start:end:path`` reference block showing the **current** offending code, and/or
   - a plain fenced `go` block with the **sample fix**.

Rules for the tables:

- Keep table cells single-line. Never put multi-line code fences *inside* a table cell — put code in a fenced block right after the table and reference it from the `Suggested fix` cell.
- Every finding must appear in both the index table and as its own detail table; keep the two in sync (same ID, severity, status, title).
- `Severity` and `Status` cells use the bare word (`Critical`, `Open`, …), not the bracketed `[Critical]` form, since they now live in a column.
- For resolved items, set `How to reproduce` to `n/a (resolved)` and use `Suggested fix` to describe what was done.

## Workflow

Copy this checklist and track progress:

```
- [ ] 1. Load the existing review (read docs/codereview.md: sections, legends, item IDs + statuses)
- [ ] 2. Inventory and read the source (all .go files + Dockerfile, .github/, go.mod)
- [ ] 3. Verify each existing finding against current code; set status + fix citations
- [ ] 4. Apply the four lenses; add new findings with severity, status, location, reproduce, sample fix
- [ ] 5. Render everything as tables: refresh the Findings index + each per-finding detail table
- [ ] 6. Update the re-review note (today's date), the Closed list, and the Priority Checklist
- [ ] 7. Report a short summary (closed / new / still open)
```

### Step 1 — Load the existing review

Read `docs/codereview.md` fully. Note the section structure (1–11), the existing item numbers, their current status, the legends, and the table layout above. Reuse them; do not invent a new format.

### Step 2 — Inventory and read the source

Enumerate and read the actual code so every judgement is grounded:

- All Go files (e.g. `cmd/`, `controller/`, `memorystore/`, `model/`).
- Infrastructure: `Dockerfile`, `.dockerignore`, `.github/workflows/`, `go.mod`/`go.sum`.

Base every status and finding on what the code says now — open the file before deciding.

### Step 3 — Verify each existing finding

For every item already in the doc:

- Open the cited code and decide the status: `Resolved`, `Partial`, or `Open`.
- Update stale citations (paths, line numbers, renamed symbols) in the `Location` cell and any code blocks to match current code.
- Keep resolved items in the document (mark them `Resolved` / closed) — do not delete history.
- If something cannot be verified, leave the status and note the uncertainty rather than guessing.

### Step 4 — Add new findings (four lenses)

Walk the code through each lens above and add any new issues. For every new finding fill the full detail table:

- An id continuing the section numbering (e.g. next `1.x`, `5.x`).
- `Severity` (`Critical|High|Medium|Low`) and `Status` (usually `Open`).
- `Lens` (e.g. "Hacker", "Architect").
- `Location` — real `file:line`.
- `Issue` / `Impact` — what's wrong and why it matters (state the lens).
- `How to reproduce` — concrete steps.
- `Suggested fix` — prose + a short Go sample below the table.

Do not fabricate issues: every finding must map to real code you have read and cited.

### Step 5 — Render the tables

- Rebuild the **Findings index** table so it lists every item (resolved + open) with its current severity/status.
- Ensure each finding's detail table is present and consistent with the index row.

### Step 6 — Update summaries

- Refresh the **Re-review note** at the top with today's date and a short bullet list of what changed since the last pass.
- Move newly fixed items into the **Closed** checklist.
- Re-order the **Priority Checklist** so the most important open items (security/correctness first) lead.

### Step 7 — Report

Summarize for the user: which items were closed, which are new, and the top remaining risks.

## Conventions

Reuse the document's existing legends verbatim:

- **Severity:** `Critical` data loss/incorrect behavior/security exposure · `High` important correctness/design · `Medium` quality/maintainability · `Low` polish.
- **Status:** `Resolved` fixed · `Partial` partly addressed · `Open` outstanding.

Heading format for an item: `### <id> — <short title>`, immediately followed by its `Field | Detail` table.

Cite existing code with the doc's reference fence (line range + path), placed *after* the detail table, for example:

```text
```45:55:memorystore/in_memory_todo.go
... cited code ...
```
```

Propose new/changed code with a plain fenced `go` block (no line numbers), also after the table.

## Guardrails

- **Write exactly one file: `docs/codereview.md`.** Never edit source code, infrastructure, or any other file (including other files under `docs/`). Reading other files is allowed and expected.
- Ground every status and finding in code you actually read; cite real `file:line`.
- Preserve the document structure and legends; append new items, don't renumber existing ones.
- Mark fixed items `Resolved` rather than removing them.
- Keep table cells single-line; put code samples in fenced blocks after the table.
- Keep examples concise and idiomatic Go.
