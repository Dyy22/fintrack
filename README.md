# Fintrack

Fintrack is a personal finance tracker built as a monorepo with a Go backend and a React frontend.

## Project structure

```txt
fintrack/
├── fintrack-backend/   # Go + Gin API, PostgreSQL migrations, tests
├── fintrack-web/       # React + Vite + TypeScript frontend
├── docker-compose.yml  # Local full-stack compose setup
└── .github/workflows/  # CI pipeline
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
cp .env.example .env
make dev
```

Run tests:

```bash
cd fintrack-backend
go test ./...
```

Apply database migrations:

```bash
cd fintrack-backend
make db-migrate
```

### Frontend only

```bash
cd fintrack-web
cp .env.example .env
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

## Deployment

Production deployment is configured through GitHub Actions:

- Frontend: Vercel
- Backend: Render
- Database: Neon Postgres or any hosted PostgreSQL provider

Workflow file:

```txt
.github/workflows/deploy.yml
```

It runs after the `CI` workflow succeeds on `main` and can also be started manually from GitHub Actions.

### Required GitHub repository secrets

Set these in GitHub: `Settings` → `Secrets and variables` → `Actions`.

#### Render

```txt
RENDER_DEPLOY_HOOK_URL=your_render_deploy_hook_url
```

The backend service is configured for Render in `render.yaml` and deploys from `fintrack-backend` using the existing backend `Dockerfile`.

#### Vercel

```txt
VERCEL_TOKEN=your_vercel_token
VERCEL_ORG_ID=your_vercel_team_or_user_id
VERCEL_PROJECT_ID=your_vercel_project_id
```

The frontend deploys from `fintrack-web` using Vercel CLI.

### Required production environment variables

#### Backend service on Render

Set these in the Render backend service environment variables:

```txt
APP_ENV=production
DATABASE_URL=your_neon_pooled_connection_string
JWT_SECRET=replace-with-a-strong-secret
JWT_EXPIRES_IN=24h
CORS_ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
```

`DATABASE_URL` is supported directly for Neon/Postgres deployments. If you do not use `DATABASE_URL`, the backend falls back to `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, and `DB_PASSWORD`.

#### Frontend project on Vercel

Set this in the Vercel project environment variables:

```txt
VITE_API_BASE_URL=https://your-render-backend-domain.onrender.com/api/v1
```

### Database migrations

Apply migrations from `fintrack-backend/migrations` to the production database before using the deployed app. You can run them manually with the backend migration tooling against the Neon production database.

## Package READMEs

See package-specific documentation:

- [`fintrack-backend/README.md`](fintrack-backend/README.md)
- [`fintrack-web/README.md`](fintrack-web/README.md)
