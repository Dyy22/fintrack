# Fintrack

[![CI](https://github.com/Dyy22/fintract/actions/workflows/ci.yml/badge.svg)](https://github.com/Dyy22/fintract/actions/workflows/ci.yml)
[![Deploy](https://github.com/Dyy22/fintract/actions/workflows/deploy.yml/badge.svg)](https://github.com/Dyy22/fintract/actions/workflows/deploy.yml)
[![Production Database Migrations](https://github.com/Dyy22/fintract/actions/workflows/migrate.yml/badge.svg)](https://github.com/Dyy22/fintract/actions/workflows/migrate.yml)

Fintrack is a personal finance tracker built as a monorepo with a Go backend and a React frontend.

## Deployment status

| Component | Platform | Status |
| --- | --- | --- |
| Frontend | Vercel | Production deployed |
| Backend API | Render | Production deployed |
| Database | Neon Postgres | Production configured and migrated |
| CI/CD | GitHub Actions | `CI` then `Deploy` on `main`; manual `Production Database Migrations` |

Production flow:

```txt
push to main -> CI -> Render backend deploy hook; Vercel deploys frontend from Git integration
manual workflow_dispatch -> Production Database Migrations -> Neon
```

## Project structure

```txt
fintrack/
├── fintrack-backend/   # Go + Gin API, PostgreSQL migrations, tests
├── fintrack-web/       # React + Vite + TypeScript frontend
├── render.yaml         # Render backend service blueprint
├── docker-compose.yml  # Local full-stack compose setup
└── .github/workflows/  # CI and deployment workflows
```

## Stack

### Backend

- Go 1.21+
- Gin
- PostgreSQL
- JWT authentication
- SQL migrations

### Frontend

- React
- TypeScript
- Vite
- Tailwind CSS
- Zustand
- Axios

### Production hosting

- Vercel for the frontend
- Render for the backend API
- Neon Postgres for the database
- GitHub Actions for CI/CD

## Features

- Authentication
- Account tracking
- Transaction tracking
- Transfers
- Category reports
- Net worth dashboard
- Income vs spending summary
- Gold account tracking in grams
- Antam gold price display and 7-day trend snapshots
- Dark mode
- Responsive neobrutal UI

## Local development

### Full stack with Docker Compose

From the repository root:

```bash
docker compose up --build
```

Default services:

- Backend API: `http://localhost:8080/api/v1`
- Frontend: `http://localhost:3000`
- PostgreSQL: `localhost:5432`

### Backend only

```bash
cd fintrack-backend
make dev
```

Run tests:

```bash
cd fintrack-backend
go test ./...
```

Apply local database migrations:

```bash
cd fintrack-backend
make db-migrate
```

### Frontend only

```bash
cd fintrack-web
npm install
npm run dev
```

Set the backend API URL with:

```txt
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

Build and lint:

```bash
cd fintrack-web
npm run lint
npm run build
```

## CI

GitHub Actions runs on `main` / `master` pushes and pull requests:

- Backend tests
- Go formatting check
- Frontend lint
- Frontend build
- Docker image build checks

Workflow file:

```txt
.github/workflows/ci.yml
```

## API documentation

The backend API is documented with OpenAPI in [`docs/openapi.yaml`](docs/openapi.yaml). You can open it with Swagger Editor, Redoc, Scalar, or any OpenAPI-compatible tooling.

## Deployment

Production deployment is configured through GitHub Actions:

- Frontend: Vercel
- Backend: Render
- Database: Neon Postgres

Workflow file:

```txt
.github/workflows/deploy.yml
```

The deploy workflow runs after the `CI` workflow succeeds on `main`. It can also be started manually from GitHub Actions.

### Required GitHub repository secrets

Set these in GitHub: `Settings` → `Secrets and variables` → `Actions`.

#### Render

```txt
RENDER_DEPLOY_HOOK_URL=your_render_deploy_hook_url
```

The backend deploys from `fintrack-backend` using the existing backend `Dockerfile`. The Render service blueprint is stored in `render.yaml`.

#### Vercel

No GitHub Actions secret is required for the frontend when Vercel Git integration is enabled. Vercel deploys `fintrack-web` directly from the connected GitHub repository.

### Required production environment variables

#### Backend service on Render

```txt
APP_ENV=production
DATABASE_URL=your_neon_pooled_connection_string
JWT_SECRET=replace-with-a-strong-secret
JWT_EXPIRES_IN=24h
CORS_ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
```

`DATABASE_URL` is supported directly for Neon/Postgres deployments. If `APP_ENV=production`, the backend requires `DATABASE_URL` to be configured.

#### Frontend project on Vercel

```txt
VITE_API_BASE_URL=https://your-render-backend-domain.onrender.com/api/v1
```

After changing Vercel environment variables, redeploy the frontend so Vite can bake the value into the production build.

### Database migrations

Apply migrations from `fintrack-backend/migrations` to the production database before using the deployed app.

Recommended production path:

1. Add this repository secret in GitHub Actions:

   ```txt
   PRODUCTION_DATABASE_URL=your_neon_pooled_connection_string
   ```

2. Open **Actions** → **Production Database Migrations**.
3. Click **Run workflow** from the `main` branch.
4. Type `MIGRATE` in the confirmation field.

The same migration script also supports local execution against Neon if needed:

```bash
cd fintrack-backend
DATABASE_URL="your_neon_pooled_connection_string" make db-migrate
```

If `psql` is not installed locally, the script falls back to running the PostgreSQL client through Docker.

## Package READMEs

See package-specific documentation:

- [`fintrack-backend/README.md`](fintrack-backend/README.md)
- [`fintrack-web/README.md`](fintrack-web/README.md)

## License

This project is open source under the [MIT License](LICENSE).
