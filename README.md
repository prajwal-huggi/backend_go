# Go Backend Project – Production Architecture Guide

This document describes the **standard workflow and architecture** used to build a **production-ready Go backend service**.

The goal is to follow:

* Clean architecture principles
* Proper dependency management
* Scalable project structure
* Secure authentication and authorization
* Maintainable backend design

---

# 1. Initialize the Go Project

Create the project directory and initialize the Go module.

```bash
mkdir backend
cd backend

go mod init github.com/prajwal-huggi/backend_go
go mod tidy
```

### Install Core Dependencies

```bash
go get -u github.com/ilyakaznacheev/cleanenv
go get github.com/jackc/pgx/v5
go get github.com/joho/godotenv
go get github.com/golang-jwt/jwt/v5
go get github.com/go-chi/chi/v5
```

### Install Migration Tools

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

```
goose -dir migrations -s create create_users_table sql
```

---

# 2. Run PostgreSQL using Docker

Create a `docker-compose.yml` file.

Example:

```yaml
version: "3.9"

services:
  postgres:
    image: postgres:16
    container_name: go_postgres
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: go_app
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

Start PostgreSQL:

```bash
docker-compose up -d
docker ps
```

Connect to the database:

```bash
docker exec -it go_postgres psql -U postgres -d go_app
```

---

# 3. Environment Configuration

Create a `.env` file.

Example:

```
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_app
DB_SSLMODE=disable

JWT_SECRET=supersecretkey
JWT_EXPIRY=24h
```

Purpose:

* Separate configuration from code
* Allow different environments (local, staging, production)

---

# 4. Go Module Files

### go.mod

Defines project dependencies.

### go.sum

Stores cryptographic checksums of dependencies to ensure integrity.

---

# 5. Production Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   └── app.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── domain/
│   │   └── user.go
│   │
│   ├── repository/
│   │   └── user_repository.go
│   │
│   ├── service/
│   │   └── user_service.go
│   │
│   ├── auth/
│   │
│   ├── middleware/
│   │
│   └── transport/
│       └── http/
│           └── user_handler.go
│
├── migrations/
│
├── pkg/
│   └── logger/
│       └── logger.go
│
├── scripts/
│
├── docker/
│   └── Dockerfile
│
├── docker-compose.yml
├── .env
├── go.mod
├── go.sum
└── README.md
```

---

# 6. Application Entry Point

File:

```
cmd/server/main.go
```

Responsibilities:

* Load configuration
* Initialize database
* Initialize repositories
* Initialize services
* Start HTTP server

Important rule:

Main file should **only bootstrap the application**.

---

# 7. Configuration Layer

File:

```
internal/config/config.go
```

Responsibilities:

* Load environment variables
* Provide configuration to the application
* Prevent hardcoded secrets

---

# 8. Database Layer

File:

```
internal/database/postgres.go
```

Responsibilities:

* Create PostgreSQL connection pool
* Validate database connection
* Manage connection lifecycle

Connection pooling improves:

* performance
* scalability
* reliability

---

# 9. Why Use pgxpool

Database connections are expensive.

Opening a connection for every request causes performance issues.

Using connection pooling:

```
Application
      ↓
Connection Pool
      ↓
PostgreSQL
```

Benefits:

* faster queries
* fewer TCP handshakes
* better scalability

---

# 10. Database Migrations

Folder:

```
migrations/
```

Purpose:

* version control database schema
* enable safe schema changes
* support reproducible deployments

Migration tools used:

* goose
* golang-migrate

---

# 11. Repository Layer

Folder:

```
internal/repository
```

Responsibilities:

* execute SQL queries
* map database rows to Go structs
* isolate SQL from business logic

Example operations:

* create user
* fetch users
* update user
* delete user

---

# 12. Domain Layer

Folder:

```
internal/domain
```

Defines core business entities.

Example:

```
User
Order
Product
```

Domain models should remain **independent of infrastructure**.

---

# 13. Service Layer

Folder:

```
internal/service
```

Responsibilities:

* business rules
* validation
* orchestration

Examples:

* register user
* process order
* validate payment

Services use repositories to access data.

---

# 14. Authentication Module

Folder:

```
internal/auth
```

Responsibilities:

* JWT token generation
* token validation
* password hashing

Authentication flow:

```
Login request
      ↓
Verify password
      ↓
Generate JWT
      ↓
Return token
```

---

# 15. Middleware Layer

Folder:

```
internal/middleware
```

Responsibilities:

* JWT authentication
* RBAC authorization
* logging
* request validation

Middleware executes **before the request handler**.

---

# 16. RBAC (Role-Based Access Control)

Roles define what actions users can perform.

Example roles:

```
admin
user
manager
```

Example permissions:

```
admin → delete users
user → view profile
manager → manage resources
```

Middleware verifies permissions before allowing access.

---

# 17. HTTP Transport Layer

Folder:

```
internal/transport/http
```

Responsibilities:

* define API routes
* parse incoming requests
* call services
* return JSON responses

Popular Go routers:

* chi
* gin
* echo

---

# 18. API Routing Example

```
POST   /auth/login
POST   /users
GET    /users
GET    /users/{id}
PUT    /users/{id}
DELETE /users/{id}
```

Protected routes require JWT authentication.

---

# 19. Request Flow in Production Systems

```
HTTP Request
      ↓
Router
      ↓
Middleware (JWT / RBAC)
      ↓
Handler
      ↓
Service
      ↓
Repository
      ↓
Database
```

The response travels back through the same layers.

---

# 20. Additional Production Components

Most production backends also integrate:

### Logging

* zap
* zerolog

### Metrics

* Prometheus

### Distributed Tracing

* OpenTelemetry

### Caching

* Redis

### Messaging / Queueing

* Kafka
* RabbitMQ

---

# 21. Key Backend Engineering Principle

Each layer should have **one responsibility**.

```
Handler      → HTTP logic
Service      → Business logic
Repository   → Database access
Middleware   → Cross-cutting concerns
```

Maintaining this separation ensures the backend remains:

* maintainable
* scalable
* testable
* production-ready

domain acting similar to model but still different I don't know how and why
database models belong to the repository layer
dtos belong to the transport layer

## request flow
```
HTTP request
     ↓
CreateUserRequest (DTO)
     ↓
User (domain entity)
     ↓
Repository
     ↓
Database
```

## repsonse flow
```
Database row
     ↓
User (domain entity)
     ↓
UserResponse (DTO)
     ↓
JSON response
```

## visual architecture
```
transport/http
      ↓
DTOs
      ↓
service
      ↓
domain
      ↓
repository
      ↓
database
```