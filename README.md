[![SonarCloud analysis](https://github.com/rahulkrishnanfs/todolist/actions/workflows/sonarcloud.yml/badge.svg)](https://github.com/rahulkrishnanfs/todolist/actions/workflows/sonarcloud.yml) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![Bugs](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=bugs)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![Coverage](https://sonarcloud.io/api/project_badges/measure?project=rahulkrishnanfs_todolist&metric=coverage)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist)

[![Quality gate](https://sonarcloud.io/api/project_badges/quality_gate?project=rahulkrishnanfs_todolist)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist) [![SonarQube Cloud](https://sonarcloud.io/images/project_badges/sonarcloud-light.svg)](https://sonarcloud.io/summary/new_code?id=rahulkrishnanfs_todolist)
 
# Todo List

A small TODO list application written in Go, structured around clean / hexagonal
architecture (ports and adapters). Domain models and persistence are decoupled
through repository interfaces, so the storage backend can be swapped without
touching business logic.

Currently the project ships with an in-memory store and an HTTP REST API
(`net/http`) served from `cmd/main.go` on port `:8080`.

## Features

- Domain models for `TODO` items and `Category` groupings.
- Repository interfaces (`ToDoRepository`, `CategoryRepository`) acting as ports.
- In-memory adapter implementation (`TodoMap`, `CategoryMap`).
- Controllers (`TODOController`, `CategoryController`) that depend on the
  abstractions, not concrete storage.
- Structured JSON logging via `log/slog`, created in `main` and injected into
  the controllers.
- Sentinel domain errors (`ErrObjectNotFound`, `ErrObjectAlreadyExists`,
  `ErrStoreEmpty`) returned by the stores.
- Meaningful HTTP status codes for writes (`201 Created` on create,
  `204 No Content` on update/delete).

## Project Structure

```text
todolist/
├── cmd/
│   ├── main.go                 # Entrypoint: wires stores, controllers, routes; starts HTTP server
│   ├── routes.go               # ToDoRoutes: registers /api/v1/todos handlers
│   └── category_routes.go      # CategoryRoutes: registers /api/v1/categories handlers
├── controller/
│   ├── todo_controller.go      # TODOController (depends on ToDoRepository)
│   └── category_controller.go  # CategoryController (depends on CategoryRepository)
├── memorystore/
│   ├── in_memory_todo.go       # In-memory adapter: TodoMap
│   └── in_memory_category.go   # In-memory adapter: CategoryMap
├── model/
│   └── model.go                # Domain entities + repository interfaces (ports)
├── go.mod                      # Module: todolist (Go 1.22)
└── README.md
```

## Getting Started

Requirements: Go 1.22+.

```bash
# From the repository root
go run ./cmd
```

This starts the HTTP server in `cmd/main.go`, which wires up the in-memory
stores, controllers, and routes, then listens on `:8080`.

### API Endpoints

| Method | Path | Description |
| --- | --- | --- |
| POST | `/api/v1/todos` | Create a TODO |
| GET | `/api/v1/todos` | List all TODOs |
| GET | `/api/v1/todos/{id}` | Get a TODO by id |
| PUT | `/api/v1/todos/{id}` | Update a TODO |
| DELETE | `/api/v1/todos/{id}` | Delete a TODO by id |
| POST | `/api/v1/categories` | Create a category |
| GET | `/api/v1/categories` | List all categories |
| GET | `/api/v1/categories/{id}` | Get a category by id |
| PUT | `/api/v1/categories/{id}` | Update a category |
| DELETE | `/api/v1/categories/{id}` | Delete a category by id |

Example:

```bash
curl -X POST localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -d '{"tid":"1","activity":"Write docs","description":"Update README","isdone":false}'

curl localhost:8080/api/v1/todos
```

## Architecture Overview

The application follows a ports-and-adapters layout:

- An HTTP layer (`ServeMux` + route files in `cmd/`) maps RESTful `/api/v1/*`
  routes to controller methods.
- Controllers (HTTP handlers) depend on repository **interfaces**, never on a
  concrete store.
- Repository interfaces (`ToDoRepository`, `CategoryRepository`) are the
  **ports** defined alongside the domain model.
- The in-memory store (`TodoMap`, `CategoryMap`) is one **adapter**
  implementing those ports. Other adapters (e.g. SQL, file) could be added
  without changing controllers or domain logic.
- Domain models (`TODO`, `Category`) are persistence-independent.
- A `*slog.Logger` (JSON handler writing to stdout) is constructed in
  `cmd/main.go` and injected into both controllers, which emit structured logs
  for each request.

```text
HTTP route  ->  Controller (handler)  ->  Repository interface (port)  ->  In-memory adapter  ->  Domain model
                      |
                      +--> structured logs (slog JSON -> stdout)
```

## C4 Architecture Diagrams

The diagrams below follow the [C4 model](https://c4model.com/) and are written
in Mermaid, which renders natively on GitHub.

### Level 1 - System Context

```mermaid
flowchart TD
    user["API Client / User<br/>[Person]<br/>Tracks tasks and categories"]
    system["Todo List API<br/>[Software System]<br/>Manages TODO items and categories over HTTP"]

    user -->|"Sends HTTP/JSON requests<br/>(create / update / delete / read)"| system
```

### Level 2 - Container

```mermaid
flowchart TD
    user["API Client / User<br/>[Person]"]

    subgraph todoApp [Todo List API]
        api["HTTP API Server<br/>[Container: Go net/http on :8080]<br/>Routes requests to controllers, logs via slog"]
        store["In-Memory Store<br/>[Container: Go maps]<br/>Holds TODOs and Categories in memory"]
    end

    logs["Structured Logs<br/>[stdout: JSON via log/slog]"]

    user -->|"HTTP/JSON over :8080<br/>/api/v1/todos, /api/v1/categories"| api
    api -->|"Reads / writes via repository ports"| store
    api -->|"Emits structured JSON logs"| logs
```

### Level 3 - Component

```mermaid
flowchart TD
    user["API Client / User<br/>[Person]"]

    subgraph server [HTTP API Server - cmd package]
        mux["ServeMux + Routes<br/>[Component]<br/>Maps /api/v1/* routes to handlers"]
    end

    subgraph handlers [Controllers - controller package]
        todoCtrl["TODOController<br/>[Component]<br/>HTTP handlers for TODO use cases"]
        catCtrl["CategoryController<br/>[Component]<br/>HTTP handlers for Category use cases"]
    end

    subgraph ports [Ports - model package]
        todoRepo["ToDoRepository<br/>[Interface]<br/>Create / Update / Delete / GetById / GetAll"]
        catRepo["CategoryRepository<br/>[Interface]<br/>Create / Update / Delete / GetByID / GetAll"]
    end

    subgraph adapters [Adapters - memorystore package]
        todoMap["TodoMap<br/>[Component]<br/>In-memory TODO store"]
        catMap["CategoryMap<br/>[Component]<br/>In-memory Category store"]
    end

    subgraph domain [Domain - model package]
        todoModel["TODO<br/>[Entity]"]
        catModel["Category<br/>[Entity]"]
        errs["Sentinel errors<br/>[ErrObjectNotFound, ErrObjectAlreadyExists, ErrStoreEmpty]"]
    end

    subgraph obs [Observability]
        logger["slog.Logger<br/>[Component]<br/>JSON handler writing to stdout"]
    end

    user -->|"HTTP/JSON"| mux
    mux -->|"routes to"| todoCtrl
    mux -->|"routes to"| catCtrl

    todoCtrl -->|"depends on"| todoRepo
    catCtrl -->|"depends on"| catRepo

    todoCtrl -->|"logs via"| logger
    catCtrl -->|"logs via"| logger

    todoMap -.->|"implements"| todoRepo
    catMap -.->|"implements"| catRepo

    todoMap -->|"returns"| errs
    catMap -->|"returns"| errs

    todoMap -->|"stores"| todoModel
    catMap -->|"stores"| catModel
```

## Roadmap / Future Work

- Persistent storage adapter (SQL or file-based) implementing the existing
  repository ports.
- Input validation and consistent JSON error responses.
- Authentication / authorization for the API.
- Tests for adapters and controllers.
