[![SonarCloud analysis](https://github.com/rahulkrishnanfs/todolist/actions/workflows/sonarcloud.yml/badge.svg)](https://github.com/rahulkrishnanfs/todolist/actions/workflows/sonarcloud.yml) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=bugs)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=coverage)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist)

[![Quality gate](https://sonarcloud.io/api/project_badges/quality_gate?project=rahulkrishnanfs_todolist)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![SonarQube Cloud](https://sonarcloud.io/images/project_badges/sonarcloud-light.svg)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist)
 
# Todo List

A small TODO list application written in Go, structured around clean / hexagonal
architecture (ports and adapters). Domain models and persistence are decoupled
through repository interfaces, so the storage backend can be swapped without
touching business logic.

The project ships with an in-memory store and a REST API (`net/http`).
Write/read routes for todos and categories are protected by **RS256 JWT
authentication**; clients obtain a token via a user signup/login flow. The same
server also exposes a **Model Context Protocol (MCP)** endpoint at `/mcp`, so
MCP clients (e.g. Cursor) can drive the service through tools — currently a
`create_user` tool. Runtime settings (port, keystore, TLS cert/key) are loaded
from `config/properties.toml`, and the JWT signing keys are loaded from a
PKCS#12 keystore. The listen port defaults to `:8080`.

> **Transport note:** the server listens over **HTTPS/TLS**
> (`server.ListenAndServeTLS` in `cmd/main.go`), using the cert/key referenced
> by `config/properties.toml` (`server_cert` / `server_key`). The shipped cert
> is signed by a local custom CA with a `todolist.ai` SAN, so clients must
> either trust that CA (see
> [`docs/2.1_trusting_custom_ca.md`](docs/2.1_trusting_custom_ca.md)) or skip
> verification (`curl -k`). The `/mcp` endpoint (the `create_user` tool) is
> **not yet authenticated** (see the Roadmap).

## Features

- Domain models for `Todo`, `Category`, and `User`.
- Repository interfaces (`TodoRepository`, `CategoryRepository`,
  `UserRepository`) acting as ports.
- In-memory adapters (`TodoMap`, `CategoryMap`, `UserMap`).
- Controllers (`TodoController`, `CategoryController`, `UserController`) that
  depend on the abstractions, not concrete storage.
- JWT auth (`auth.Authenticator`): RS256 token generation on login and
  `AuthorizeRequest` middleware on protected routes.
- MCP integration (`pkg/mcptools`): an MCP server (`modelcontextprotocol/go-sdk`)
  mounted on the same `ServeMux` at `/mcp` via `NewStreamableHTTPHandler`,
  exposing a `create_user` tool backed by `UserController.CreateUserTool`.
  (DNS-rebinding protection is currently disabled via
  `DisableLocalhostProtection: true` for local development.)
- HTTPS server via `http.Server.ListenAndServeTLS`, using the PEM cert/key paths
  from config (the plain-HTTP `ListenAndServe` call is kept commented out).
- PKCS#12 keystore loading (`utils.Secret`) into an RSA key pair for JWT signing.
- TOML configuration (`utils.Config`) for port, keystore, and TLS cert/key
  settings.
- Structured JSON logging via `log/slog`, created in `main` and injected into
  the controllers.
- Sentinel domain errors (`ErrObjectNotFound`, `ErrObjectAlreadyExists`,
  `ErrStoreEmpty`) defined in the model; the stores return `ErrObjectNotFound` /
  `ErrObjectAlreadyExists`, and `GetAll` returns an empty list (not an error)
  when the store is empty.
- Meaningful HTTP status codes for writes (`201 Created` on create,
  `204 No Content` on update/delete).

## Project Structure

```text
todolist/
├── cmd/
│   └── main.go                     # Entrypoint: loads config + keystore, wires stores/controllers/auth/routes + MCP tools, starts HTTPS server
├── pkg/
│   ├── controller/
│   │   ├── todo_controller.go      # TodoController (depends on TodoRepository)
│   │   ├── category_controller.go  # CategoryController (depends on CategoryRepository)
│   │   └── user_controller.go      # UserController: signup + login (issues JWT) + CreateUserTool (MCP)
│   ├── mcptools/
│   │   └── user_mcp.go             # SetTools: registers the create_user MCP tool + mounts /mcp (StreamableHTTPHandler)
│   ├── router/
│   │   ├── todo_routes.go          # SetTodoRoutes: /api/v1/todos handlers (JWT-protected)
│   │   ├── category_routes.go      # SetCategoryRoutes: /api/v1/categories handlers (JWT-protected)
│   │   └── user_routes.go          # SetUserRoutes: /api/v1/users signup + login (public)
│   ├── memorystore/
│   │   ├── in_memory_todo.go       # In-memory adapter: TodoMap
│   │   ├── in_memory_category.go   # In-memory adapter: CategoryMap
│   │   └── in_memory_user.go       # In-memory adapter: UserMap
│   ├── auth/
│   │   └── auth.go                 # Authenticator: RS256 JWT + AuthorizeRequest middleware
│   ├── utils/
│   │   ├── config.go               # Config: loads config/properties.toml (port, keystore, TLS cert/key)
│   │   └── secrets.go              # Secret: loads PKCS#12 keystore -> RSA key pair (JWT signing)
│   └── model/
│       └── model.go                # Domain entities (Todo, Category, User) + repository ports
├── config/
│   └── properties.toml             # Service config: port, keystore path + password, TLS cert/key paths
├── charts/todolist/                # Helm chart: Deployment, Service, ConfigMap, Secret, Ingress
├── scripts/                        # docker-run.sh, k8s-secrets.sh helpers
├── secrets/                        # PKCS#12 keystore (JWT) + custom-CA/server PEM TLS cert+key (keep real secrets out of git)
├── docs/                           # codereview.md, design/, HTTPS status codes, cert/key generation, trusting the CA, k8s notes
├── magfile.go                      # Mage targets: build, docker login/build/push
├── Dockerfile                      # Multi-stage build -> distroless nonroot image
├── .github/workflows/              # CI: SonarCloud, golangci-lint, AI review, issue labeler, moderator
├── go.mod                          # Module: todolist (Go 1.25)
└── README.md
```

All application packages live under `pkg/` (imported as `todolist/pkg/...`); `cmd/main.go` is the only `package main` and just wires everything together.

## Getting Started

Requirements: **Go 1.25+**. A PKCS#12 keystore (for JWT signing) **and** a TLS
certificate/key pair (PEM, used by `ListenAndServeTLS`) are required at startup,
along with the config file. See
[`docs/2_cert_key_generation.md`](docs/2_cert_key_generation.md) for the OpenSSL
commands that generate the keystore and the CA-signed server cert/key, and
[`docs/2.1_trusting_custom_ca.md`](docs/2.1_trusting_custom_ca.md) for trusting
the custom CA so clients don't need `curl -k`.

1. Ensure `config/properties.toml` points at a valid keystore and TLS cert/key:

```toml
[service]
port = ":8080"
keystore_file_path = "/absolute/path/to/secrets/keystore.p12"
keystore_password = "changeit"
server_cert = "/absolute/path/to/secrets/fullchain.crt"
server_key = "/absolute/path/to/secrets/server.key"
```

2. Run from the repository root (config is read from `./config/properties.toml`):

```bash
go run ./cmd
```

This loads the config and keystore, wires up the in-memory stores, controllers,
auth, REST routes, and the MCP tools, then serves HTTPS on the configured port
(default `:8080`). Todo/category endpoints require a JWT — sign up and log in
first to get one.

### API Endpoints

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| POST | `/api/v1/users/signup` | None | Register a user |
| POST | `/api/v1/users/login` | None | Authenticate; returns a JWT |
| POST | `/api/v1/todos` | Bearer | Create a TODO |
| GET | `/api/v1/todos` | Bearer | List all TODOs |
| GET | `/api/v1/todos/{id}` | Bearer | Get a TODO by id |
| PUT | `/api/v1/todos/{id}` | Bearer | Update a TODO |
| DELETE | `/api/v1/todos/{id}` | Bearer | Delete a TODO by id |
| POST | `/api/v1/categories` | Bearer | Create a category |
| GET | `/api/v1/categories` | Bearer | List all categories |
| GET | `/api/v1/categories/{id}` | Bearer | Get a category by id |
| PUT | `/api/v1/categories/{id}` | Bearer | Update a category |
| DELETE | `/api/v1/categories/{id}` | Bearer | Delete a category by id |
| POST | `/mcp` | None | MCP endpoint (streamable HTTP) exposing the `create_user` tool |

Example flow (the server runs over HTTPS; `-k` skips verification of the
self-signed cert — trust the custom CA to drop it):

```bash
# 1. Sign up
curl -k -X POST https://localhost:8080/api/v1/users/signup \
  -H 'Content-Type: application/json' \
  -d '{"uid":"u1","username":"alice","password":"s3cret","email_address":"alice@example.com"}'

# 2. Log in and capture the token
TOKEN=$(curl -sk -X POST https://localhost:8080/api/v1/users/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"alice","password":"s3cret"}' | jq -r .token)

# 3. Call a protected route with the bearer token
curl -k -X POST https://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"tid":"1","activity":"Write docs","description":"Update README","is_done":false,"category_id":"c1","user_id":"alice"}'

curl -k https://localhost:8080/api/v1/todos -H "Authorization: Bearer $TOKEN"
```

### MCP (Model Context Protocol)

The server registers an MCP server (`name: "todolist"`, `version: 1.0.0`) and
mounts it at `/mcp` using the go-sdk's streamable-HTTP handler. It currently
exposes one tool:

| Tool | Description | Required arguments |
| --- | --- | --- |
| `create_user` | Create a new user account for the todolist service | `uid`, `username`, `password`, `email_address` |

To consume it from an MCP client such as Cursor, add the endpoint to your
`mcp.json` (project-level `.cursor/mcp.json` or global `~/.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "todolist": {
      "type": "http",
      "url": "https://localhost:8080/mcp"
    }
  }
}
```

Because the server uses a self-signed (custom-CA) certificate, the MCP client
must trust that CA — otherwise the connection is rejected. See
[`docs/2.1_trusting_custom_ca.md`](docs/2.1_trusting_custom_ca.md). Start the
server first, then enable `todolist` in Cursor's MCP settings; the `create_user`
tool will become available to the agent.

> The endpoint currently has **no authentication** and DNS-rebinding protection
> is disabled (`DisableLocalhostProtection: true`) for local development — do
> not expose it as-is (see the Roadmap).

## Architecture Overview

The application follows a ports-and-adapters layout:

- An HTTPS layer: `http.Server.ListenAndServeTLS` serves a `ServeMux` (wired by
  the `pkg/router` package) that maps RESTful `/api/v1/*` routes to controller
  methods, using the cert/key from config. (A commented-out plain-HTTP
  `ListenAndServe` path remains for local experiments.)
- An MCP layer: `pkg/mcptools.SetTools` registers the `create_user` tool on an
  `mcp.Server` and mounts it on the same `ServeMux` at `/mcp` via
  `NewStreamableHTTPHandler`, so MCP clients hit the service alongside REST.
- Controllers (HTTP handlers) depend on repository **interfaces**, never on a
  concrete store.
- Repository interfaces (`TodoRepository`, `CategoryRepository`,
  `UserRepository`) are the **ports** defined alongside the domain model.
- The in-memory stores (`TodoMap`, `CategoryMap`, `UserMap`) are one set of
  **adapters** implementing those ports. Other adapters (e.g. SQL, file) could
  be added without changing controllers or domain logic.
- Domain models (`Todo`, `Category`, `User`) are persistence-independent.
- An **auth layer** (`auth.Authenticator`) wraps protected routes with
  `AuthorizeRequest`, which verifies the RS256 JWT using the public key;
  `UserController.Login` mints tokens with the private key.
- **Config** (`utils.Config`) loads `config/properties.toml`; **secrets**
  (`utils.Secret`) decode a PKCS#12 keystore into the RSA key pair injected into
  the `Authenticator`. The TLS server certificate and key (PEM) listed in the
  config are consumed at startup by `ListenAndServeTLS`.
- A `*slog.Logger` (JSON handler writing to stdout) is constructed in
  `cmd/main.go` and injected into the controllers, which emit structured logs
  for each request.

```text
HTTPS  ->  route  ->  JWT middleware  ->  Controller (handler)  ->  Repository interface (port)  ->  In-memory adapter  ->  Domain model
   |                                            |
   |                                            +--> structured logs (slog JSON -> stdout)
   +--> /mcp (StreamableHTTPHandler) -> create_user tool -> UserController.CreateUserTool -> UserRepository
```

## C4 Architecture Diagrams

The diagrams below follow the [C4 model](https://c4model.com/) and are written
in Mermaid, which renders natively on GitHub.

### Level 1 - System Context

```mermaid
flowchart TD
    user["API Client / User<br/>[Person]<br/>Tracks tasks and categories"]
    agent["AI Agent / MCP Client<br/>[Person/System]<br/>e.g. Cursor"]
    system["Todo List API<br/>[Software System]<br/>Manages TODO items and categories over HTTPS"]

    user -->|"Signs up / logs in for a JWT, then<br/>sends authenticated HTTPS/JSON requests"| system
    agent -->|"Calls MCP tools (create_user)<br/>over HTTPS at /mcp"| system
```

### Level 2 - Container

```mermaid
flowchart TD
    user["API Client / User<br/>[Person]"]
    agent["AI Agent / MCP Client<br/>[Person/System]"]

    subgraph todoApp [Todo List API]
        api["HTTPS API Server<br/>[Container: Go net/http TLS on :8080]<br/>Routes, JWT auth, slog logging"]
        mcpsrv["MCP Server<br/>[Component: go-sdk StreamableHTTP at /mcp]<br/>create_user tool"]
        authc["JWT Authenticator<br/>[Component: golang-jwt RS256]<br/>Signs + verifies bearer tokens"]
        store["In-Memory Stores<br/>[Container: Go maps]<br/>TODOs, Categories, Users"]
    end

    cfg["Config File<br/>[config/properties.toml]"]
    keys["RSA Keystore<br/>[secrets/keystore.p12 - PKCS#12]"]
    tls["TLS Cert + Key<br/>[secrets/fullchain.crt + server.key]"]
    logs["Structured Logs<br/>[stdout: JSON via log/slog]"]

    user -->|"HTTPS/JSON + Bearer JWT<br/>/api/v1/*"| api
    agent -->|"MCP over HTTPS<br/>/mcp"| mcpsrv
    api -->|"same ServeMux"| mcpsrv
    api -->|"sign / verify tokens"| authc
    api -->|"reads / writes via repository ports"| store
    mcpsrv -->|"create_user -> Create"| store
    api -->|"loads port + keystore + cert paths"| cfg
    api -->|"serves TLS with"| tls
    authc -->|"loads RSA key pair"| keys
    api -->|"emits structured JSON logs"| logs
```

> The TLS cert/key (`secrets/fullchain.crt` + `secrets/server.key`) are loaded at
> startup by `ListenAndServeTLS`; the cert is signed by a local custom CA with a
> `todolist.ai` SAN.

### Level 3 - Component

```mermaid
flowchart TD
    user["API Client / User<br/>[Person]"]
    agent["AI Agent / MCP Client<br/>[Person/System]"]

    subgraph server [HTTP API Server - cmd + router packages]
        httpsrv["http.Server<br/>[Component: cmd]<br/>ListenAndServeTLS on :8080"]
        mux["ServeMux + Set*Routes<br/>[Component: pkg/router]<br/>Maps /api/v1/* to handlers"]
        mw["AuthorizeRequest<br/>[Middleware: pkg/auth]<br/>Validates JWT on todo/category routes"]
    end

    subgraph mcppkg [MCP - pkg/mcptools package]
        setTools["SetTools<br/>[Component]<br/>Registers tools + mounts /mcp"]
        mcpHandler["/mcp StreamableHTTPHandler<br/>[Component]<br/>create_user tool"]
    end

    subgraph authpkg [Auth - pkg/auth package]
        authr["Authenticator<br/>[Component]<br/>GenerateJWT (RS256) + AuthorizeRequest"]
    end

    subgraph handlers [Controllers - pkg/controller package]
        todoCtrl["TodoController<br/>[Component]"]
        catCtrl["CategoryController<br/>[Component]"]
        userCtrl["UserController<br/>[Component]<br/>signup + login -> JWT; CreateUserTool"]
    end

    subgraph ports [Ports - pkg/model package]
        todoRepo["TodoRepository<br/>[Interface]"]
        catRepo["CategoryRepository<br/>[Interface]"]
        userRepo["UserRepository<br/>[Interface]<br/>Create / Login"]
    end

    subgraph adapters [Adapters - pkg/memorystore package]
        todoMap["TodoMap<br/>[Component]"]
        catMap["CategoryMap<br/>[Component]"]
        userMap["UserMap<br/>[Component]"]
    end

    subgraph cfgpkg [Config + Secrets - pkg/utils package]
        cfg["Config<br/>[Component]<br/>Loads properties.toml"]
        secret["Secret<br/>[Component]<br/>Loads PKCS#12 -> RSA keys"]
    end

    subgraph domain [Domain - pkg/model package]
        todoModel["Todo<br/>[Entity]"]
        catModel["Category<br/>[Entity]"]
        userModel["User<br/>[Entity]"]
        errs["Sentinel errors<br/>[ErrObjectNotFound, ErrObjectAlreadyExists, ErrStoreEmpty]"]
    end

    logger["slog.Logger<br/>[Component]<br/>JSON handler writing to stdout"]

    user -->|"HTTPS/JSON (+ Bearer JWT)"| httpsrv
    agent -->|"MCP over HTTPS /mcp"| mcpHandler
    httpsrv -->|"serves"| mux
    mux --> mw
    mux -->|"mounts /mcp"| mcpHandler
    setTools -->|"registers create_user + handler"| mcpHandler
    mcpHandler -->|"create_user"| userCtrl
    mw -->|"verify token"| authr
    mw -->|"protected routes"| todoCtrl
    mw -->|"protected routes"| catCtrl
    mux -->|"public routes"| userCtrl

    userCtrl -->|"GenerateJWT"| authr
    todoCtrl -->|"depends on"| todoRepo
    catCtrl -->|"depends on"| catRepo
    userCtrl -->|"depends on"| userRepo

    todoMap -.->|"implements"| todoRepo
    catMap -.->|"implements"| catRepo
    userMap -.->|"implements"| userRepo

    secret -->|"RSA keys"| authr
    cfg -->|"port + keystore path"| secret
    cfg -->|"port + cert/key paths"| httpsrv

    todoMap -->|"stores"| todoModel
    catMap -->|"stores"| catModel
    userMap -->|"stores"| userModel
    todoMap -->|"returns"| errs

    todoCtrl -->|"logs via"| logger
    catCtrl -->|"logs via"| logger
    userCtrl -->|"logs via"| logger
```

## Developer Tooling (Cursor Skills)

This repo includes [Cursor](https://cursor.com) Agent Skills under `.cursor/skills/`.
Each skill lives in `.cursor/skills/<name>/SKILL.md` and is invoked by name in the
Cursor agent.

| Skill | What it does | How to use |
| --- | --- | --- |
| `commit-with-issue` | Creates a Conventional Commit (issue number at the end of the subject) on a matching `feature/#[issue]-[branch]` branch. | `/commit-with-issue` |
| `code-review-update` | Re-scans the source and updates `docs/codereview.md` (finding statuses + new findings) from architect, senior-engineer, hacker, and security perspectives. Writes only that file. | `/code-review-update` |
| `codereview-to-issues` | Opens a GitHub issue for every **Open** finding in `docs/codereview.md`, mapping defects to the bug-report template and design gaps to the feature-request template (via the GitHub MCP). | `/codereview-to-issues` |
| `readme-update` | Re-scans the source and updates this `README.md` (structure, API, C4 diagrams, tooling) to stay accurate for newcomers. Writes only this file. | `/readme-update` |

### `commit-with-issue`

When invoked, the skill:

1. Inspects your staged/unstaged changes.
2. Asks for the GitHub issue number (required) and uses it with a `#` prefix.
3. Picks the right Conventional Commit type (`feat`, `fix`, `docs`, `refactor`, …).
4. Creates a branch named `feature/#[issue]-[branchname]`.
5. Commits with a `type(scope): subject (#[issue])` message — the issue number
   goes at the end of the subject line.

Commits are authored by **you** (the developer making the change). Configure your
git identity once so authorship is attributed correctly:

```bash
git config user.name "Your Name"
git config user.email "you@example.com"
```

The skill definitions live in `.cursor/skills/<skill-name>/SKILL.md` — read or edit
them there to adjust the workflow.

### `code-review-update`

Run `/code-review-update` to refresh [`docs/codereview.md`](docs/codereview.md). It
re-scans the source, re-verifies each existing finding (flipping its status to
Resolved / Partial / Open), and adds new findings reviewed from four perspectives:
software architect, senior expert programmer, hacker/attacker, and security expert.
It only ever writes `docs/codereview.md`.

### `codereview-to-issues`

Run `/codereview-to-issues` to turn the **Open** findings in
[`docs/codereview.md`](docs/codereview.md) into GitHub issues — one per finding.
Defects are filed with the bug-report template and missing-capability/design
items with the feature-request template, created one by one (with dedup) through
the GitHub MCP server (`user-github`).

### `readme-update`

Run `/readme-update` to refresh this README from the code. It reconciles the
overview, project structure, API endpoints, architecture, and C4 diagrams with the
source, and keeps this Developer Tooling table in sync. It only ever writes
`README.md`.

## Roadmap / Future Work

- Hash passwords (currently stored/compared in plaintext) with bcrypt.
- Per-user authorization / data ownership (scope todos and categories to the
  authenticated user).
- Move secrets out of the repo; load keystore path/password from environment or
  a secret manager (TLS keys and the keystore are currently under `secrets/`).
- Replace the developer-machine absolute paths in `config/properties.toml` with
  container-relative paths / env overrides so the image runs outside the dev box.
- Secure the `/mcp` endpoint: add authentication (the MCP `create_user` tool is
  currently unauthenticated) and re-enable DNS-rebinding protection (drop
  `DisableLocalhostProtection`).
- Align the CI Go version (`.github/workflows/sonarcloud.yml` pins 1.22.0) with
  `go.mod` (1.25); the lint workflow and Dockerfile build stage already use 1.25.
- Persistent storage adapter (SQL or file-based) implementing the existing
  repository ports.
- Input validation and consistent JSON error responses (map errors to
  400/404/409).
- Tests for adapters, controllers, and the auth flow.
