# Fintrack Web

React + TypeScript frontend for Fintrack.

## Stack

- Vite
- React
- TypeScript
- Tailwind CSS
- React Router
- Zustand
- Axios

## Local setup

```bash
npm install
npm run dev
```

The frontend expects the backend API at:

```txt
http://localhost:8080/api/v1
```

Override with:

```txt
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

For production builds, set `VITE_API_BASE_URL` before building because Vite embeds environment variables at build time.

## Scripts

```bash
npm run dev      # Start Vite dev server
npm run build    # Type-check and build production assets
npm run preview  # Preview built assets locally
npm run lint     # Run ESLint
```

## Production deployment

The frontend is deployed to Vercel from the monorepo subdirectory:

```txt
fintrack-web
```

Vercel project settings:

```txt
Framework Preset: Vite
Root Directory: fintrack-web
Build Command: npm run build
Output Directory: dist
```

Production environment variable:

```txt
VITE_API_BASE_URL=https://your-render-backend-domain.onrender.com/api/v1
```

The app includes `vercel.json` rewrites so React Router routes can be refreshed directly in production.

GitHub Actions deploys the frontend through Vercel CLI using these repository secrets:

```txt
VERCEL_TOKEN
VERCEL_ORG_ID
VERCEL_PROJECT_ID
```

## Backend integration

The API client is configured in:

```txt
src/services/api.ts
```

It reads `VITE_API_BASE_URL` and falls back to local development:

```txt
http://localhost:8080/api/v1
```

If requests from Vercel fail with a CORS error, update the backend Render environment variable:

```txt
CORS_ALLOWED_ORIGINS=https://your-vercel-app.vercel.app
```

Then redeploy the Render backend.
