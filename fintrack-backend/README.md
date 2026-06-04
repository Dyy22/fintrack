# Fintrack Backend

Go + Gin backend for Fintrack.

## Requirements

- Go 1.21+
- PostgreSQL 15+
- golang-migrate CLI

## Local setup

```bash
cp .env.example .env
make dev
make smoke-test
```

Useful development commands:

```bash
make db-up          # Start PostgreSQL only
make db-migrate     # Apply migration if not already applied
make db-reset       # Drop/recreate schema using migration down/up
make docker-up      # Build/start backend and PostgreSQL
make docker-down    # Stop services
make logs           # Follow backend logs
```

The smoke test runs an end-to-end API flow: health check, register, login, account types, accounts, categories, expense, transfer, and reports.

Run unit tests:

```bash
go test ./...
```

Run PostgreSQL integration tests:

```bash
createdb fintrack_test
FINTRACK_TEST_DB_URL="postgres://fintrack:fintrack@localhost:5432/fintrack_test?sslmode=disable" go test -tags=integration ./internal/repository/postgres
```

Or via Makefile:

```bash
make test-integration
```

Integration tests reset the target database by running the migration down/up files. For safety, the database name must contain `test`.

Health check:

```bash
curl http://localhost:8080/api/v1/health
```

## Main endpoints

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/account-types`
- `GET/POST/PUT/DELETE /api/v1/accounts`
- `GET/POST/PUT/DELETE /api/v1/categories`
- `GET/POST /api/v1/transactions`
- `GET /api/v1/reports/net-worth`
- `GET /api/v1/reports/spending-by-category`
