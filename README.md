# BookShare API

A modular, production-grade REST API built with **Go**, focused on demonstrating clean architecture, background processing, secure auth, Docker orchestration, and testable code structure.  

This project is designed to showcase **backend engineering capabilities** suitable for real-world production systems.

---

## Tech Stack Overview

| Layer | Tech |
|-------|------|
| Language | [Go 1.23+](https://golang.org) |
| Web framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Background Jobs | [Asynq (Redis)](https://github.com/hibiken/asynq) |
| Auth | JWT (with refresh token via Redis) |
| Database | PostgreSQL |
| Rate Limiting | Custom Redis-based per-user, per-route |
| API Docs | [Swagger UI via Swag](https://github.com/swaggo/swag) |
| Testing | `testing`, `testify`, integration with Redis and DB |
| DevOps | Docker, Docker Compose, GitHub Actions (CI) |

---

## Features

### Auth System
- JWT access + refresh tokens (Redis-backed)
- Secure login, logout, token refresh
- Email verification via background queue
- Role-based access (e.g. Admin)

### Books API (CRUD)
- Authenticated user access
- Only creators can update/delete their own books

### Admin Panel (API-level)
- View all users
- Change user roles (promote to admin)

### Background Processing
- Email sending handled via Redis + Asynq
- Worker service runs independently of API

### Rate Limiting (Advanced)
- Per-route, per-role limits (e.g. 5/min for `/login`)
- Redis-backed counters
- `X-RateLimit-*` headers returned

### Testing
- Unit and integration tests for auth, registration, CRUD, workers
- Real Redis + PostgreSQL used in tests (isolated DB)
- Token store and worker logic tested with mocks/fakes

---

## Project Structure

bookshare-api/
├── cmd/               # Entry points
│   ├── api/           # HTTP server (main.go)
│   └── worker/        # Background worker
├── internal/          # All application logic
│   ├── user/          # Registration, auth, user info
│   ├── books/         # CRUD logic
│   ├── admin/         # Admin-only handlers
│   ├── middleware/    # JWT, AdminOnly, RateLimiter
│   ├── task/          # Redis/Asynq distributor & processor
│   ├── db/            # GORM + migrate setup
├── migrations/        # SQL schema migrations
├── tests/             # Integration test setup
├── Dockerfile, docker-compose.yml
├── .env.example       # Sample environment config

---

## Usage & Setup

### 1. Setup Env

Copy `.env.example` and fill in:

```env
HOST=<api_host>
PORT=<api_port>

DB_SOURCE=postgresql://postgres:postgres@postgres:5432/bookshare?sslmode=disable
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=bookshare

REDIS_ADDR=redis:6379

JWT_SECRET=<your_secret>

SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=<user_email>
SMTP_PASS=<pass>

DEBUG=1

TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5433/bookshare_test?sslmode=disable
TEST_REDIS_URL=localhost:6379
```



---

### 2. Run with Docker Compose

docker-compose up --build
- API runs on http://localhost:8080
- Swagger docs: http://localhost:8080/swagger/index.html
- Background worker started as separate service
- PostgreSQL and Redis containers are bootstrapped

---

### 3. API Demo: Register + Verify + Login

### Register user
```
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "strongpass"}'
```

### Simulate email verification (normally queued)
```
curl "http://localhost:8080/api/v1/verify?token=...&uid=..."
```

### Login
```
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "strongpass"}'
```


---

### Run Tests

`go test ./... -v`

Make sure Redis and PostgreSQL are running. You can isolate test DB using a separate schema or connection string.

---

### Swagger API Docs

After launching, visit:

http://localhost:8080/swagger/index.html

---

### Why This Project?

This project is designed not to showcase product features, but backend engineering capabilities:
	•	Clean structure
	•	Scalable task queues
	•	Real-world auth
	•	Middleware design
	•	Dockerized workflows
	•	Production-ready practices

---

### Author

DMaryanskiy
Backend Developer — Go / Python / PostgreSQL / Redis
