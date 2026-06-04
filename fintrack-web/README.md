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
cp .env.example .env
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

## Scripts

```bash
npm run dev
npm run build
npm run preview
npm run lint
```

## Current status

This is the initial scaffold. It includes routing, layout, placeholder pages, API client, auth store, and formatting utilities. Feature integration will be implemented page by page.
