# TaskFlow

> REST API for managing projects, tasks, project members, comments and task history.

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.26.2-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.12.0-00ADD8?style=for-the-badge)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![JWT](https://img.shields.io/badge/Auth-JWT-000000?style=for-the-badge&logo=jsonwebtokens)](https://jwt.io/)
[![Goose](https://img.shields.io/badge/Migrations-Goose-00ADD8?style=for-the-badge)](https://github.com/pressly/goose)

![Tests](https://img.shields.io/badge/tests-go%20test%20.--passed-brightgreen?style=flat-square)

</div>

## Overview

TaskFlow is a backend REST API written in Go. The application provides the core backend functionality of a lightweight project and task management system.

The API supports:

- user registration and login;
- JWT-based authentication for protected routes;
- project creation, retrieval, update and deletion;
- project membership management and role changes;
- task creation, retrieval, update and deletion;
- task assignment and unassignment;
- task comments;
- task status history retrieval.

The project is organized into separate HTTP, business-logic and data-access layers. Services depend on repository interfaces, which keeps the business logic independent from PostgreSQL implementations and makes the service layer straightforward to unit test.

## Features

### Authentication

- Registration with password hashing using bcrypt.
- Login with email and password verification.
- JWT generation after successful authentication.
- Authentication middleware for protected endpoints.
- `GET /me` for the authenticated user.

### Projects

- Create a project.
- Get a project by ID.
- Get projects owned by the authenticated user.
- Update project name and description.
- Delete a project.
- The project owner is automatically added to `project_members` with the owner role when a project is created.

### Project Members

- Add a user to a project.
- List project members.
- Get a member by project and user ID.
- Check project membership.
- Change a member role.
- Remove a member.
- Owner role cannot be changed or removed.

### Tasks

- Create a task inside a project.
- Get a task by ID.
- List tasks belonging to a project.
- Update task data.
- Delete a task.
- Assign a task to a user.
- Remove task assignment.

### Comments

- Create a comment for a task.
- List comments for a task.
- Get a comment by ID.
- Delete a comment.

### Task History

- Store task history records.
- Get a history record by ID.
- Get history for a task.
- Get history for all tasks in a project.

## Tech Stack

| Technology              | Purpose                               |
| ----------------------- | ------------------------------------- |
| Go 1.26.2               | Application language                  |
| Gin                     | HTTP server and routing               |
| PostgreSQL 16           | Relational database                   |
| pgx / pgxpool           | PostgreSQL driver and connection pool |
| JWT                     | Authentication                        |
| bcrypt                  | Password hashing                      |
| Goose                   | Database migrations                   |
| Docker / Docker Compose | Local containerized environment       |
| Make                    | Common development commands           |

## Architecture

TaskFlow uses a layered backend structure:

```text
                    HTTP Client
                         |
                         v
              +----------------------+
              |     Gin Router       |
              |  Auth Middleware     |
              +----------+-----------+
                         |
                         v
              +----------------------+
              |       Handlers        |
              |   HTTP / JSON layer   |
              +----------+-----------+
                         |
                         v
              +----------------------+
              |       Services       |
              |    Business Logic     |
              +----------+-----------+
                         |
                         v
              +----------------------+
              | Repository Interfaces|
              +----------+-----------+
                         |
                         v
              +----------------------+
              | PostgreSQL Repositories|
              +----------------------+
                         |
                         v
                    PostgreSQL
```

### Handlers

Handlers are responsible for HTTP concerns: reading path parameters, binding JSON requests, validating basic request shape, calling services and converting service errors into HTTP responses.

### Services

Services contain business rules and coordinate repository operations. They validate IDs and input data, check related entities and return domain-specific errors such as `ErrUserNotFound`, `ErrProjectNotFound` and `ErrTaskNotFound`.

### Repositories

Repositories contain PostgreSQL-specific data-access code implemented with `pgxpool`. The service layer depends on interfaces declared in `internal/service/interfaces.go`, not on concrete repository types.

This separation allows the service layer to be tested with mocks instead of a real database.

## Project Structure

```text
TaskFlow/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── notification/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── middleware.go
│   │
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── comment_handler.go
│   │   ├── project_handler.go
│   │   ├── project_member_handler.go
│   │   ├── task_handler.go
│   │   ├── task_history_handler.go
│   │   └── router.go
│   │
│   ├── model/
│   │   ├── user.go
│   │   ├── project.go
│   │   ├── project_member.go
│   │   ├── task.go
│   │   ├── comment.go
│   │   └── taskhistory.go
│   │
│   ├── repository/
│   │   ├── repository.go
│   │   ├── user_repository.go
│   │   ├── project_repository.go
│   │   ├── project_member_repository.go
│   │   ├── task_repository.go
│   │   ├── comment_repository.go
│   │   └── task_history_repository.go
│   │
│   └── service/
│       ├── interfaces.go
│       ├── errors.go
│       ├── user_service.go
│       ├── project_service.go
│       ├── project_member_service.go
│       ├── task_service.go
│       ├── comment_service.go
│       ├── task_history_service.go
│       └── *_test.go
│
├── migrations/
│   ├── *_create_users.sql
│   ├── *_create_projects.sql
│   ├── *_create_project_members.sql
│   ├── *_create_tasks.sql
│   ├── *_create_comments.sql
│   ├── *_create_task_history.sql
│   └── *_create_indexes.sql
│
├── .dockerignore
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## API

All routes except `/register` and `/login` require a valid JWT in the `Authorization` header.

### Authentication

| Method | Endpoint    | Description                    | Auth |
| ------ | ----------- | ------------------------------ | ---- |
| POST   | `/register` | Register a new user            | No   |
| POST   | `/login`    | Authenticate and receive a JWT | No   |
| GET    | `/me`       | Get the authenticated user     | Yes  |

### Projects

| Method | Endpoint        | Description                             | Auth |
| ------ | --------------- | --------------------------------------- | ---- |
| POST   | `/projects`     | Create a project                        | Yes  |
| GET    | `/projects`     | List projects owned by the current user | Yes  |
| GET    | `/projects/:id` | Get a project by ID                     | Yes  |
| PUT    | `/projects/:id` | Update a project                        | Yes  |
| DELETE | `/projects/:id` | Delete a project                        | Yes  |

### Project Members

| Method | Endpoint                              | Description          | Auth |
| ------ | ------------------------------------- | -------------------- | ---- |
| POST   | `/projects/:id/members`               | Add a member         | Yes  |
| GET    | `/projects/:id/members`               | List project members | Yes  |
| GET    | `/projects/:id/members/:userID`       | Get a member         | Yes  |
| GET    | `/projects/:id/members/:userID/check` | Check membership     | Yes  |
| PATCH  | `/projects/:id/members/:userID`       | Update member role   | Yes  |
| DELETE | `/projects/:id/members/:userID`       | Remove a member      | Yes  |

### Tasks

| Method | Endpoint              | Description                | Auth |
| ------ | --------------------- | -------------------------- | ---- |
| POST   | `/projects/:id/tasks` | Create a task in a project | Yes  |
| GET    | `/projects/:id/tasks` | List project tasks         | Yes  |
| GET    | `/tasks/:id`          | Get a task                 | Yes  |
| PUT    | `/tasks/:id`          | Update a task              | Yes  |
| DELETE | `/tasks/:id`          | Delete a task              | Yes  |
| PATCH  | `/tasks/:id/assign`   | Assign a task to a user    | Yes  |
| DELETE | `/tasks/:id/assign`   | Remove task assignment     | Yes  |

### Comments

| Method | Endpoint              | Description        | Auth |
| ------ | --------------------- | ------------------ | ---- |
| POST   | `/tasks/:id/comments` | Create a comment   | Yes  |
| GET    | `/tasks/:id/comments` | List task comments | Yes  |
| GET    | `/comments/:id`       | Get a comment      | Yes  |
| DELETE | `/comments/:id`       | Delete a comment   | Yes  |

### Task History

| Method | Endpoint                | Description                        | Auth |
| ------ | ----------------------- | ---------------------------------- | ---- |
| GET    | `/tasks/:id/history`    | Get history for a task             | Yes  |
| GET    | `/task-history/:id`     | Get one history record             | Yes  |
| GET    | `/projects/:id/history` | Get history for tasks in a project | Yes  |

## Authentication

The authentication flow is:

```text
POST /register
       |
       v
   User created
       |
       v
POST /login
       |
       v
   JWT returned
       |
       v
Authorization: Bearer <token>
       |
       v
AuthMiddleware
       |
       v
userID stored in Gin context
       |
       v
Protected handler
```

Passwords are stored as bcrypt hashes. Clients do not send or receive password hashes through the API response models.

Example protected request:

```bash
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer <token>"
```

## Configuration

### Local development

Create `.env` in the project root:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/taskflow
PORT=8080
```

`.env` is a local configuration file and should not be committed to Git. Use `.env.example` as a template.

### Docker Compose

Inside the Docker network, the API connects to PostgreSQL using the Compose service name `postgres`:

```text
postgres://postgres:postgres@postgres:5432/taskflow
```

The Docker Compose configuration provides the database settings to the API container through environment variables.

## Running Locally

Requirements:

- Go 1.26.2 or newer.
- PostgreSQL 16 or newer, or Docker Desktop.
- Docker Compose for the containerized setup.
- Make is optional.

Install dependencies:

```bash
go mod download
```

Run the API directly:

```bash
go run ./cmd/api
```

The API listens on:

```text
http://localhost:8080
```

## Running with Docker Compose

The Compose environment contains three services:

```text
+-------------------+
|    PostgreSQL     |
|      :5432        |
+---------+---------+
          |
          v
+-------------------+
|      Goose        |
|    migrations     |
+---------+---------+
          |
          v
+-------------------+
|    TaskFlow API   |
|      :8080        |
+-------------------+
```

Start the stack:

```bash
docker compose up -d --build
```

Check the services:

```bash
docker compose ps -a
```

Expected state:

```text
taskflow-postgres   Up (healthy)
taskflow-migrate    Exited (0)
taskflow-api        Up
```

The migration container is intentionally one-shot. `Exited (0)` means the migration command completed successfully.

View logs:

```bash
docker compose logs -f api
docker compose logs -f migrate
docker compose logs -f postgres
```

Stop the stack:

```bash
docker compose down
```

Stop the stack and remove the PostgreSQL volume:

```bash
docker compose down -v
```

> `docker compose down -v` removes the database volume and therefore deletes the PostgreSQL data stored in that volume.

## Database Migrations

TaskFlow uses [Goose](https://github.com/pressly/goose) for schema migrations.

Migration files are stored in `migrations/` and use Goose's `Up` / `Down` format:

```sql
-- +goose Up

CREATE TABLE example (
    id BIGSERIAL PRIMARY KEY
);

-- +goose Down

DROP TABLE example;
```

The current migration set creates the application's main tables and indexes, including:

- `users`
- `projects`
- `project_members`
- `tasks`
- `comments`
- `task_history`

Goose stores the applied migration version in:

```text
goose_db_version
```

In Docker Compose, migrations are executed before the API starts. The API service depends on successful completion of the migration service.

## Testing

The service layer is covered with unit tests using repository mocks. The tests do not require a live PostgreSQL connection.

Run the complete test suite:

```bash
go test ./...
```

Run service tests with verbose output:

```bash
go test ./internal/service -v
```

Run tests for a specific service:

```bash
go test ./internal/service -run TestUserService -v
go test ./internal/service -run TestProjectService -v
go test ./internal/service -run TestTaskService -v
go test ./internal/service -run TestProjectMemberService -v
go test ./internal/service -run TestCommentService -v
go test ./internal/service -run TestTaskHistoryService -v
```

The test suite covers successful operations together with invalid input and not-found scenarios for the service layer.

## Development Commands

The project includes a `Makefile` for common development tasks. The exact commands available in the current checkout can be inspected with:

```bash
make
```

Typical Go commands used during development are:

```bash
go fmt ./...
go vet ./...
go test ./...
go build ./cmd/api
docker compose up -d --build
docker compose down
```

## Example API Flow

### 1. Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "ivan",
    "email": "ivan@example.com",
    "password": "password123"
  }'
```

Example response:

```json
{
	"id": 1,
	"username": "ivan",
	"email": "ivan@example.com"
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ivan@example.com",
    "password": "password123"
  }'
```

Example response:

```json
{
	"token": "<jwt>"
}
```

### 3. Get current user

```bash
curl -X GET http://localhost:8080/me \
  -H "Authorization: Bearer <jwt>"
```

### 4. Create a project

```bash
curl -X POST http://localhost:8080/projects \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt>" \
  -d '{
    "name": "TaskFlow",
    "description": "Project and task management"
  }'
```

### 5. Create a task

```bash
curl -X POST http://localhost:8080/projects/1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt>" \
  -d '{
    "title": "Implement authentication",
    "description": "Add JWT authentication"
  }'
```

### 6. Create a comment

```bash
curl -X POST http://localhost:8080/tasks/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt>" \
  -d '{
    "text": "Authentication is ready for review"
  }'
```

## Design Decisions

### Repository aggregation

`internal/repository/repository.go` provides a single `Repository` struct that holds the concrete repository implementations:

```go
type Repository struct {
    Comments       *CommentRepository
    Projects       *ProjectRepository
    ProjectMembers *ProjectMemberRepository
    Tasks          *TaskRepository
    TaskHistory    *TaskHistoryRepository
    Users          *UserRepository
}
```

This keeps dependency wiring in `main.go` simple while preserving separate repositories for each domain.

### Dependency inversion

Services depend on repository interfaces declared in `internal/service/interfaces.go`:

```go
type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id int64) (*model.User, error)
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    Update(ctx context.Context, user *model.User) error
    Delete(ctx context.Context, id int64) error
}
```

Concrete PostgreSQL repositories satisfy these interfaces in `internal/repository`.

### Domain-specific errors

Business errors are defined centrally in `internal/service/errors.go`. Handlers translate these service errors into appropriate HTTP status codes.

For example:

```text
ErrUserNotFound        -> 404 Not Found
ErrProjectNotFound     -> 404 Not Found
ErrTaskNotFound        -> 404 Not Found
ErrInvalidTaskID       -> 400 Bad Request
ErrInvalidProjectName  -> 400 Bad Request
```

This keeps transport-specific response logic in handlers while allowing services to remain independent from HTTP.

## Verification

Before committing changes, the project can be checked with:

```bash
go fmt ./...
go vet ./...
go test ./...
go build ./cmd/api
docker compose config
docker compose up -d --build
docker compose ps -a
```

A successful Docker deployment should leave PostgreSQL healthy, the migration container at `Exited (0)`, and the API container running on port `8080`.

## Current Scope

The current repository contains the implemented backend API, service-level unit tests, PostgreSQL repositories, Goose migrations and Docker Compose environment.

The `cmd/notification` entry point is present in the repository, but the current documented API is the service exposed by `cmd/api`.

## License

No license is declared in the current project configuration. Add a `LICENSE` file before publishing the repository as an open-source project under a specific license.

---

<div align="center">

**TaskFlow — Go backend project**

</div>
