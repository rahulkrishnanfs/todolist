# Contributing to Todo List

Thanks for your interest in contributing! This project is a small Go TODO API
built around clean / hexagonal architecture (ports and adapters). These
guidelines explain how to set up your environment, the conventions we follow,
and how to propose changes.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Coding Guidelines](#coding-guidelines)
- [Architecture Guidelines](#architecture-guidelines)
- [Commit Messages](#commit-messages)
- [Cursor Skills](#using-the-commit-with-issue-cursor-skill)
- [Pull Requests](#pull-requests)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

Be respectful and constructive. Assume good intent, keep discussions focused on
the code and the problem at hand, and help keep this a welcoming project for
everyone.

## Getting Started

Requirements: **Go 1.22+**.

```bash
# Clone and enter the repo
git clone <your-fork-url>
cd todolist

# Download dependencies
go mod download

# Run the HTTP server (listens on :8080)
go run ./cmd
```

Verify it works:

```bash
curl localhost:8080/api/v1/todos
```

See the [README](README.md) for the full list of API endpoints and example
requests.

## Project Structure

```text
todolist/
├── cmd/            # Entrypoint: wires stores, controllers, routes; starts server
├── controller/     # HTTP handlers (depend on repository interfaces)
├── memorystore/    # In-memory adapter implementations
├── model/          # Domain entities + repository interfaces (ports)
└── docs/           # Supporting notes and references
```

Before adding code, find the layer it belongs to. Keep concerns separated:
HTTP/JSON handling stays in `controller/`, persistence in adapters such as
`memorystore/`, and persistence-independent types and interfaces in `model/`.

## Development Workflow

1. Fork the repository and create a feature branch from `main`:

```bash
git checkout -b feat/short-description
```

2. Make your change in small, focused commits.
3. Format, vet, and build before pushing (see below).
4. Open a pull request describing what changed and why.

### Before You Commit

Run these from the repository root and make sure they pass:

```bash
gofmt -l .          # should print nothing; run `gofmt -w .` to fix
go vet ./...        # static checks
go build ./...      # everything compiles
go test ./...       # run tests (once present)
```

## Coding Guidelines

- **Formatting:** All Go code must be `gofmt`-formatted. CI-style reviews will
  reject unformatted code.
- **Naming:** Follow standard Go conventions (exported identifiers documented
  with a leading comment, mixedCaps, short receiver names).
- **Errors:** Define reusable sentinel errors in `model/model.go` (for example
  `ErrObjectNotFound`) and compare with `errors.Is`. Return errors up the stack;
  let controllers translate them into HTTP status codes.
- **Logging:** Use the structured `log/slog` logger that is already threaded
  through controllers. Prefer `logger.LogAttrs(ctx, level, "message", slog.String(...))`
  with key/value attributes over interpolated strings.
- **HTTP status codes:** Map outcomes to correct codes (for example `201` on
  create, `204` on update/delete, `400` on bad input, `404` on not found,
  `500` on unexpected failures). See [docs/https_status_codes.md](docs/https_status_codes.md).
- **Comments:** Comment intent and non-obvious decisions, not what the code
  plainly does.
- **Dependencies:** This project leans on the standard library. Discuss in an
  issue before introducing a new third-party dependency.

## Architecture Guidelines

This project follows ports and adapters. Please preserve the boundaries:

- Controllers depend on repository **interfaces** (`ToDoRepository`,
  `CategoryRepository`), never on a concrete store.
- New storage backends (SQL, file, etc.) should be added as **adapters** that
  implement the existing repository interfaces, without changing controllers or
  domain models.
- Domain models in `model/` must stay persistence-independent (no database or
  transport details).

```text
HTTP route -> Controller -> Repository interface (port) -> Adapter -> Domain model
```

When adding a new endpoint, wire it through this same flow and register the
route in the `cmd/` route files.

## Commit Messages

Write clear, imperative commit messages that explain the "why":

```text
Add file-based repository adapter

Implements ToDoRepository against the local filesystem so data
survives restarts. No changes to controllers or domain models.
```

A short type prefix is encouraged (`feat:`, `fix:`, `docs:`, `refactor:`,
`test:`, `chore:`).

### Using the `commit-with-issue` Cursor skill

If you work in [Cursor](https://cursor.com), this repo provides a skill that
automates the convention above. Run `/commit-with-issue` in the agent and it
will:

- ask for the GitHub issue number (required) and reference it as `#[issue]`,
- create a branch named `feature/#[issue]-[branchname]`,
- write a Conventional Commit (`type(scope): subject`) with a `Closes #[issue]`
  footer.

The skill lives in `.cursor/skills/commit-with-issue/SKILL.md`. Commits are
attributed to your own git identity, so set it once before committing:

```bash
git config user.name "Your Name"
git config user.email "you@example.com"
```

See the README's [Developer Tooling](README.md#developer-tooling-cursor-skills)
section for the full list of available skills.

## Pull Requests

Before requesting review, make sure your PR:

- [ ] Is scoped to a single logical change.
- [ ] Passes `gofmt -l .`, `go vet ./...`, and `go build ./...`.
- [ ] Includes tests for new behavior where practical.
- [ ] Updates the [README](README.md) or docs if behavior or endpoints change.
- [ ] Has a description covering what changed and why.

A maintainer will review and may request changes. Keep the discussion in the PR
thread so context stays in one place.

## Reporting Issues

When opening an issue, include:

- What you expected to happen and what actually happened.
- Steps to reproduce (include the exact request, for example the `curl` command).
- Relevant logs (the app emits structured JSON logs via `slog`).
- Your Go version (`go version`) and OS.

Thanks for contributing!
