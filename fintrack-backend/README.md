# Fintrack Backend

Go + Gin backend API for Fintrack.

## Stack

- Go 1.21+
- Gin
- PostgreSQL
- JWT authentication
- SQL migrations
- Docker for production builds

## Requirements

- Go 1.21+
- PostgreSQL 15+ for local development
- Docker for local Compose workflows and Docker image builds
- `psql` or Docker for remote database migrations

## Local setup

Start local PostgreSQL, apply migrations, and run the backend stack:

```bash
make dev
```

Useful development commands:

```bash
make db-up          # Start PostgreSQL only
make db-migrate     # Apply pending migrations
make db-reset       # Drop/recreate local schema using migration down/up
make docker-up      # Build/start backend and PostgreSQL
make docker-down    # Stop services
make logs           # Follow backend logs
make smoke-test     # Run an end-to-end API smoke test
```

Run the API directly with Go:

```bash
go run ./cmd/api
```

Health check:

```bash
curl http://localhost:8080/api/v1/health
```

## Configuration

Common environment variables:

```txt
APP_ENV=development
PORT=8080
DATABASE_URL=postgresql://user:password@host/db?sslmode=require
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fintrack
DB_USER=fintrack
DB_PASSWORD=fintrack
JWT_SECRET=change-me-in-production
JWT_EXPIRES_IN=24h
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
GOLD_PRICE_SOURCE_URL=https://logam-mulia-api.iamutaki.workers.dev/api/prices/logammulia
GOLD_PRICE_FALLBACK_PER_GRAM=0
GOLD_PRICE_REFRESH_INTERVAL=1h
```

`DATABASE_URL` takes precedence over the individual `DB_*` variables. In production (`APP_ENV=production`), `DATABASE_URL` is required.

## Database migrations

Migration files live in:

```txt
migrations/
```

Apply local migrations:

```bash
make db-migrate
```

Apply production migrations against Neon:

```bash
DATABASE_URL="postgresql://user:password@host/neondb?sslmode=require" make db-migrate
```

If `psql` is not installed locally, the migration script falls back to running the PostgreSQL client through Docker.

## Tests

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

Integration tests reset the target database by running migration down/up files. For safety, the database name must contain `test`.

## Production deployment

The backend is deployed to Render from the monorepo root using:

```txt
render.yaml
fintrack-backend/Dockerfile
```

Render service environment variables:

```txt
APP_ENV=production
DATABASE_URL=your_neon_pooled_connection_string
JWT_SECRET=replace-with-a-strong-secret
JWT_EXPIRES_IN=24h
CORS_ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
```

Render health check path:

```txt
/api/v1/health
```

GitHub Actions triggers Render deployment through the repository secret:

```txt
RENDER_DEPLOY_HOOK_URL
```

## Main endpoints

- `GET /api/v1/health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/account-types`
- `GET/POST/PUT/DELETE /api/v1/accounts`
- `GET/POST/PUT/DELETE /api/v1/categories`
- `GET/POST /api/v1/transactions`
- `GET /api/v1/gold/price`
- `GET /api/v1/gold/prices/history`
- `GET /api/v1/reports/net-worth`
- `GET /api/v1/reports/spending-by-category`
