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

## Deployment notes

Recommended deployment split:

- Frontend: Vercel
- Backend: Render, Railway, or Fly.io
- Database: Supabase Postgres, Neon, Railway Postgres, or Render Postgres

For Vercel frontend deployment:

- Root Directory: `fintrack-web`
- Build Command: `npm run build`
- Output Directory: `dist`
- Environment variable:

```txt
VITE_API_BASE_URL=https://your-backend-domain.com/api/v1
```

## Package READMEs

See package-specific documentation:

- [`fintrack-backend/README.md`](fintrack-backend/README.md)
- [`fintrack-web/README.md`](fintrack-web/README.md)
