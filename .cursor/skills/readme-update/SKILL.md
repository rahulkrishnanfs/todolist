---
name: readme-update
description: Re-scan the entire source tree and update README.md in place — this is the ONLY file the skill writes; it reads all other files but never modifies them. Keeps the overview, project structure, features, API endpoints, architecture overview, and C4 diagrams in sync with the code; documents the Cursor skills under .cursor/skills as developer tooling; and keeps the project approachable for newcomers. Use when the user invokes this skill or asks to refresh, regenerate, or update the README.
disable-model-invocation: true
---

# README Update

Refresh `README.md` so it accurately reflects the current codebase and is easy for a newcomer to follow.

> **Scope (hard rule):** The ONLY file this skill may modify is `README.md`. It may *read* any source, infrastructure, skill, or doc file to gather information, but must never write to anything else. If a change seems to require editing another file, stop and report it instead.

## Workflow

Copy this checklist and track progress:

```
- [ ] 1. Read the current README.md
- [ ] 2. Inventory and read the source (all .go files + Dockerfile, .github/, go.mod) and .cursor/skills/*/SKILL.md
- [ ] 3. Reconcile every README section with the code (see "Sections to keep in sync")
- [ ] 4. Verify and update the C4 diagrams
- [ ] 5. Update the Developer Tooling (Cursor Skills) section from .cursor/skills
- [ ] 6. Newcomer check: can someone clone, understand, run, and call the API from the README alone?
- [ ] 7. Report a short summary of what changed
```

### Step 1 — Read the current README

Read `README.md` fully so edits are surgical and the existing tone, badges, and structure are preserved.

### Step 2 — Inventory and read the source

Ground every statement in the real code:

- All Go files (`cmd/`, `controller/`, `memorystore/`, `model/`).
- Routing: `cmd/*routes.go` for exact methods and paths.
- Infrastructure: `Dockerfile`, `.github/workflows/`, `go.mod` (Go version, module name).
- Skills: every `.cursor/skills/*/SKILL.md` (for the Developer Tooling section).

### Step 3 — Sections to keep in sync

Check each section against the code and fix any drift:

- **Overview / intro** — what the service is and how it's served.
- **Features** — capabilities that actually exist (logging, error model, status codes, etc.).
- **Project Structure** — the file tree and per-file descriptions match the real layout and symbol names.
- **Getting Started** — Go version (`go.mod`), run command, and any prerequisites are correct.
- **API Endpoints** — methods and paths exactly match `cmd/*routes.go` (verbs, `/api/v1/...`, `{id}`); the example `curl`s use real routes and JSON field names from `model/`.
- **Architecture Overview** — ports/adapters description and the text flow diagram are accurate.
- **Roadmap** — reflects what is still outstanding.

### Step 4 — Verify the C4 diagrams (always)

Always re-check the Mermaid C4 diagrams against the code and update them:

- **Level 1 (System Context):** actors and the system.
- **Level 2 (Container):** runtime pieces (HTTP server + port, store, logging sink, etc.) and their relationships.
- **Level 3 (Component):** packages and components (`mux`/routes, controllers, repository interfaces/ports, adapters, domain models, logger) and the edges between them — names must match the real types and packages.

Keep the diagrams valid Mermaid (quote labels containing `[ ]` or `( )`), and keep route labels (e.g. `/api/v1/todos`) current.

### Step 5 — Developer Tooling (Cursor Skills)

Maintain a section that lists the project's Cursor skills so contributors can use them. For each `.cursor/skills/<name>/SKILL.md`, read its `name` and `description` and produce a row:

```markdown
## Developer Tooling (Cursor Skills)

Skills live in `.cursor/skills/` and are invoked in the Cursor agent.

| Skill | What it does | How to use |
| --- | --- | --- |
| `commit-with-issue` | <one-line summary from the skill> | `/commit-with-issue` |
| `code-review-update` | <one-line summary from the skill> | `/code-review-update` |
| `readme-update` | <one-line summary from the skill> | `/readme-update` |
```

Add new skills, remove ones that no longer exist, and refresh summaries from each skill's description.

### Step 6 — Newcomer check

Before finishing, confirm a first-time reader can, from the README alone: understand the project's purpose and architecture, install prerequisites, run it, and make a successful API call. Fix anything missing or confusing. Prefer clear prose, accurate examples, and short code/diagram blocks over verbosity.

### Step 7 — Report

Summarize what changed (e.g. "updated API table, refreshed Level 3 C4, added readme-update to tooling").

## Conventions

- Preserve existing badges, section order, and overall style; make targeted edits.
- Use fenced ```mermaid``` blocks for C4 diagrams and ```text```/```bash```/```go``` blocks elsewhere as appropriate.
- Match JSON field names and routes to the code exactly (don't paraphrase tags or paths).
- Do not document Cursor as a commit co-author; attribute work to the developer.

## Guardrails

- **Write exactly one file: `README.md`.** Never edit source code, infrastructure, skills, or any other file. Reading other files is allowed and expected.
- Ground every statement in code you actually read; don't invent endpoints, fields, or features.
- Always re-verify the C4 diagrams on each run, even if other sections look unchanged.
- If there is nothing to change, say so instead of making edits for their own sake.
