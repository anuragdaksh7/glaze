# Copilot Instructions for glaze.server

## Build, Test, and Lint Commands

- **Build:**
  - Docker: `docker build -t glaze-backend .`
  - Local: `go build -o main .`
- **Run:**
  - Docker Compose: `docker-compose up`
  - Local: `go run main.go`
- **Test:**
  - No Go test files detected. Add *_test.go files for test discovery.
- **Lint:**
  - No linting scripts detected. Add a linter (e.g., golangci-lint) for static analysis.

## High-Level Architecture

- **Backend:** Go (Gin framework)
- **Entrypoint:** `main.go` initializes config, logger, DB, and routers.
- **Config:** Loaded from `.env` via Viper in `config/`.
- **Database:** PostgreSQL (see `docker-compose.yml` and `config/dbConfig.go`).
- **Models:** Located in `models/` (User, Workspace, Project, Repository, Deployment, etc.)
- **Routing:** Defined in `router/router.go`, uses Gin and CORS middleware.
- **Auth:** JWT-based, handled in `middleware/requireAuth.go`.
- **Logging:** Zap logger, optionally integrates with Axiom for production.
- **Workspace/Project/Repo:**
  - Users belong to workspaces (via WorkspaceMember).
  - Workspaces contain projects, which are linked to repositories.
  - Deployments are tied to projects and triggered via webhooks.

## Key Conventions

- **Environment:** All config is loaded from `.env` (see `config/Config` struct for keys).
- **Database Migrations:** Auto-migrated on startup via `config/SyncDB()`.
- **UUIDs:** Used for primary keys and relationships.
- **CORS:** Only specific origins allowed (see `router/router.go`).
- **Deployment:** Deployment logic should run in a goroutine or worker, not in the HTTP handler.
- **Sensitive Data:** EnvVars should be encrypted at rest (see `future/steps.txt`).

## AI Agent & Backend Development Guidelines

When generating or editing code in this backend workspace, agents must rigidly adhere to the following best practices:

### 1. Architectural Layers & Separation of Concerns
- **Router (`router/router.go`)**: Only define routes, path signatures, and attach middleware (e.g., `middleware.RequireAuth`).
- **Handlers (`internal/<feature>/<feature>_handler.go`)**:
  - Extract auth context (`utils.ExtractUser(c)`).
  - Parse/Validate requests via `*gin.Context` using `ShouldBindJSON`, `ShouldBindUri`, or `ShouldBindQuery`. Never pass `*gin.Context` down to the service layer natively.
  - Return HTTP responses strictly via the `response` package (e.g., `response.OK(c, data)`, `response.BadRequest(c, err)`).
- **Services (`internal/<feature>/<feature>_service.go`)**:
  - Contain purely business logic and database interactions. Use interfaces for dependency injection and mocking.
  - Never parse HTTP context parameters directly here. Read UUIDs and structured DTOs as input.
  - Role-based Access Control (RBAC) validations and domain rules belong *here*, returning clean domain errors.

### 2. Payload Management (DTOs)
- Never use database models for REST payloads schemas. Use `dto/<feature>/` structs.
- Decorate JSON API interactions utilizing structural tags (`json:"name"`, `uri:"workspace_id" binding:"required,uuid"`).

### 3. Database Operations (GORM) & Transactions
- **Context/Timeout**: Enforce standard timeout variables or derived contexts where meaningful on extensive loops.
- **Relational Reads**: Extensively use `.Preload("Relation")` or standard JOINs prior to parsing.
- **Transactions**: For operations mutating multiple records/tables concurrently (e.g., Delete Workspace + cascading Projects + WorkspaceMembers), *always* embed your deletions inside a `s.DB.Transaction(func(tx *gorm.DB) error { ... })` block for atomicity.

### 4. Roles & Permissions Validation
- Implement standard `models.WorkspaceRole` checks mapping strings cleanly (`Owner`, `Admin`, `Member`, `Viewer`).
- Fail fast. Deny operations gracefully by querying member association constraints before interacting with main object bodies.

### 5. Error Handling & Logging
- **Logging**: Rely on Zap (`logger.Logger.Error("msg", zap.Error(err))`). Log inside handlers prior to triggering client responses, or inside service logic loops when handling deeper subsystem errors.
- Never write bare `panic()` or `log.Fatal()` (except in early init scripts).

---

This file summarizes build/run instructions, architecture, and conventions for Copilot and future contributors. Would you like to adjust anything or add coverage for additional areas (e.g., deployment, CI/CD, or secrets management)?
