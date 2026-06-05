# Fintrack — Project Roadmap

**Version**: v1.1.0-roadmap  
**Last Updated**: 5 June 2026

> Latest tagged release: [`v1.0.0`](https://github.com/Dyy22/fintrack/releases/tag/v1.0.0)  
> Current feature PR: [#4 Add IDX markets and persistent cache](https://github.com/Dyy22/fintrack/pull/4)

---

## Legend

- ✅ Done
- 🚧 In Progress
- ⬜ Planned

---

## Current Production Status

| Component | Status | Notes |
|-----------|--------|-------|
| Repository visibility | ✅ | Public GitHub repository |
| License | ✅ | MIT License added |
| Monorepo structure | ✅ | Root repository contains backend, frontend, Docker Compose, CI/CD, and deployment config |
| Frontend deployment | ✅ | Vercel Git integration deploys `fintrack-web` |
| Backend deployment | ✅ | Render deploy hook triggered by GitHub Actions |
| Production database | ✅ | Neon Postgres configured and migrated |
| CI | ✅ | GitHub Actions validates backend, frontend, and Docker builds |
| Deploy workflow | ✅ | Runs after CI success on `main`; triggers Render backend deploy |
| Production docs | ✅ | Root/backend/frontend READMEs updated with deployment and license info |
| IDX stock market data | 🚧 | Implemented in PR #4; pending merge/deploy/migration |
| Persistent market cache | 🚧 | Implemented in PR #4 using DB-backed quote/chart cache; pending production migration |

Production flow:

```txt
push to main -> CI -> Render backend deploy hook; Vercel deploys frontend from Git integration
```

---

## Phase 0: Specification & Repository Setup

| Item | Status | Notes |
|------|--------|-------|
| Product & Technical Specification | ✅ | Maintained on documentation branch |
| UI/UX Design Specification | ✅ | Maintained on documentation branch |
| Root monorepo | ✅ | Backend and frontend consolidated into one root repository |
| Public repository | ✅ | Repo is public |
| MIT License | ✅ | `LICENSE` added; README license sections updated |
| README documentation | ✅ | Root, backend, and frontend READMEs aligned with current deployment setup |
| Markdown hygiene | ✅ | App branch keeps README/license docs; larger docs can remain on `dev` as needed |

---

## Phase 1: Backend MVP

### Core API

| Endpoint | Status | Notes |
|----------|--------|-------|
| `GET /api/v1/health` | ✅ | Render health check path |
| `POST /api/v1/auth/register` | ✅ | |
| `POST /api/v1/auth/login` | ✅ | |
| `GET /api/v1/account-types` | ✅ | |
| `GET /api/v1/accounts` | ✅ | |
| `POST /api/v1/accounts` | ✅ | |
| `PUT /api/v1/accounts/:id` | ✅ | |
| `DELETE /api/v1/accounts/:id` | ✅ | Supports `?hard=true` for permanent delete |
| `GET /api/v1/categories` | ✅ | |
| `POST /api/v1/categories` | ✅ | |
| `PUT /api/v1/categories/:id` | ✅ | |
| `DELETE /api/v1/categories/:id` | ✅ | |
| `GET /api/v1/transactions` | ✅ | With pagination (`limit`/`offset`) + filters |
| `POST /api/v1/transactions` | ✅ | |
| `GET /api/v1/reports/net-worth` | ✅ | |
| `GET /api/v1/reports/spending-by-category` | ✅ | |
| `GET /api/v1/gold/price` | ✅ | Current cached Antam gold price |
| `GET /api/v1/gold/prices/history` | ✅ | Gold price history for Markets page chart |
| `GET /api/v1/market/chart` | 🚧 | Implemented in PR #4; returns IHSG/IDX stock charts with persistent cache fallback |

### Database

| Item | Status | Notes |
|------|--------|-------|
| Schema: users, account_types, accounts, categories, transactions | ✅ | |
| Seed: account_types | ✅ | bank, ewallet, cash, gold, stock_broker |
| Seed: default categories | ✅ | Salary, Freelance, Investment, Food, Transport, Shopping, Bills, Entertainment |
| Constraints: amount > 0, transfer rules, unique categories | ✅ | |
| Migration files (up/down) | ✅ | 000001-000007 migrations; 000006 stock holdings and 000007 market cache pending merge/deploy |
| Neon production database | ✅ | Production schema migrated through v1.0.0; market cache migration pending PR #4 deployment |
| Migration script supports `DATABASE_URL` | ✅ | Can run local or remote Neon migrations |

### Testing

| Item | Status | Notes |
|------|--------|-------|
| Unit tests: config | ✅ | Includes production `DATABASE_URL` validation |
| Unit tests: middleware auth | ✅ | |
| Unit tests: response helper | ✅ | |
| Unit tests: password + JWT security | ✅ | |
| Unit tests: usecase validation | ✅ | |
| Unit tests: handler/router | ✅ | |
| Integration tests: PostgreSQL repository | ✅ | Uses build tag `integration` |
| Smoke test script | ✅ | `scripts/smoke-test.sh` |

### Backend DevOps

| Item | Status | Notes |
|------|--------|-------|
| Dockerfile (multi-stage) | ✅ | Render uses Docker build |
| Docker Compose (backend dev) | ✅ | |
| Docker Compose (root) | ✅ | Includes postgres, backend, web |
| Makefile | ✅ | dev, test, smoke-test, db migration, etc. |
| `.env.example` | ✅ | Placeholder-only examples |
| Idempotent DB migration script | ✅ | `scripts/db-migrate.sh` |
| Production config validation | ✅ | Fails fast when `APP_ENV=production` and `DATABASE_URL` is missing |
| Render deployment | ✅ | `render.yaml` + deploy hook |

---

## Phase 2: Frontend Scaffold

| Item | Status | Notes |
|------|--------|-------|
| Vite + React + TypeScript | ✅ | |
| Tailwind CSS | ✅ | |
| React Router | ✅ | |
| Zustand | ✅ | accountStore, authStore, categoryStore, reportStore, transactionStore, themeStore |
| Axios client with JWT interceptor | ✅ | Auto-attach token, handle 401 |
| Auth store | ✅ | login, register, logout, restore session |
| App layout | ✅ | Sidebar desktop, bottom nav mobile |
| Protected route wrapper | ✅ | Redirect to `/login` if unauthenticated |
| Reusable components | ✅ | Button, Card, ConfirmDialog, ErrorBoundary, Skeleton, neo components |
| Path alias `@/` | ✅ | |
| Formatters | ✅ | formatIDR, formatDate, transactionAmountLabel |
| API error parser | ✅ | validation_error mapping |
| Page titles | ✅ | usePageTitle hook |
| Favicon | ✅ | SVG |
| Dark mode toggle | ✅ | localStorage + system preference detect |
| Vercel SPA rewrites | ✅ | `fintrack-web/vercel.json` |
| Frontend production deployment | ✅ | Vercel Git integration |

---

## Phase 3: Frontend Feature Implementation

### Authentication

| Item | Status | Notes |
|------|--------|-------|
| Login page | ✅ | Controlled form, loading state, field errors, redirect |
| Register page | ✅ | Email, password, confirm password, auto-login after register |
| Logout | ✅ | Via sidebar and mobile header |
| Client-side validation | ✅ | Required fields checked before API call |

### Dashboard

| Item | Status | Notes |
|------|--------|-------|
| Net worth card | ✅ | Live from API |
| Account balances | ✅ | List with formatIDR |
| Spending this month | ✅ | Total + category breakdown |
| Recent transactions | ✅ | Last 5 |
| Loading state | ✅ | Skeleton cards |
| Empty state | ✅ | No accounts |
| Income vs spending donut chart | ✅ | SVG donut chart showing split |
| Summary hover/focus percentages | ✅ | Donut center percentage appears on segment hover/focus |
| Antam gold price trend chart | ✅ | Moved to dedicated Markets page in PR #4 |

### Accounts

| Item | Status | Notes |
|------|--------|-------|
| Account list | ✅ | Desktop table, mobile cards |
| Total balance card | ✅ | |
| Add account form | ✅ | Name, type dropdown, initial balance |
| Edit account name | ✅ | Inline |
| Deactivate account | ✅ | Confirm dialog |
| Activate account | ✅ | Reactivate deactivated accounts |
| Delete account | ✅ | Hard delete with confirm dialog |
| Gold account holdings | ✅ | Gram-based account balance using latest gold price |
| IDX stock account holdings | 🚧 | Implemented in PR #4; validates IDX tickers, stores lots, calculates value from latest quote |

### Transactions

| Item | Status | Notes |
|------|--------|-------|
| Transaction list | ✅ | Desktop table + mobile cards |
| Transaction filters | ✅ | Date range, account, type |
| Pagination controls | ✅ | Prev / Next with limit |
| Add expense/income transaction | ✅ | Account, category, amount, date, description |
| Add transfer transaction | ✅ | Source + destination account |
| Create transaction redirect | ✅ | Redirect to list |
| Client-side validation | ✅ | Required fields checked before submit |

### Reports

| Item | Status | Notes |
|------|--------|-------|
| Net worth report | ✅ | Total + per account breakdown |
| Spending by category with date range | ✅ | Date filter with Apply button |
| Simple horizontal bars | ✅ | Width = percentage of total |
| Percentage labels | ✅ | |

### Markets

| Item | Status | Notes |
|------|--------|-------|
| Dedicated Markets page | 🚧 | Implemented in PR #4; separate menu/page for gold, IHSG, and portfolio stock charts |
| IHSG chart | 🚧 | Implemented in PR #4 via Yahoo Finance chart provider |
| IDX stock charts | 🚧 | Implemented in PR #4 for symbols from stock accounts |
| Chart timeframe selector | 🚧 | Implemented in PR #4: 5D, 1M, 3M, 6M, 1Y |
| Market data timestamps | 🚧 | Implemented in PR #4 using `fetched_at` metadata |
| Persistent quote cache | 🚧 | Implemented in PR #4; DB cache TTL 5 minutes with stale fallback |
| Persistent chart cache | 🚧 | Implemented in PR #4; DB cache TTL 15 minutes with stale fallback |
| Yahoo Finance IDX compatibility | 🚧 | Implemented in PR #4; accepts JKT/Jakarta exchange metadata |

### UI Design — Neobrutalism

| Item | Status | Notes |
|------|--------|-------|
| Tailwind neo tokens | ✅ | Colors, shadow, radius, weight |
| Global CSS | ✅ | Dot grid background, neobrutal form controls, overrides |
| Neobrutalism component system | ✅ | NeoButton, NeoCard, NeoInput, NeoTextarea, NeoSelect, NeoDateInput, NeoBadge, NeoAlert, NeoProgress, NeoTable, NeoEmptyState, NeoStatCard, NeoPageHeader |
| Forms refactor | ✅ | Login, Register, NewAccount, NewTransaction |
| Data pages refactor | ✅ | Accounts, Transactions, Dashboard, Reports, Markets |
| Custom NeoDateInput | ✅ | Neobrutal popup, future date disabled |
| Custom NeoSelect | ✅ | Neobrutal popup, no native select |
| Dark mode contrast fix | ✅ | Muted text and border overrides |

### Polish

| Item | Status | Notes |
|------|--------|-------|
| Empty states | ✅ | NeoEmptyState |
| Loading states | ✅ | Skeleton cards/rows |
| Error boundaries | ✅ | ErrorBoundary component with Try Again |
| Client-side validation errors | ✅ | Field-level errors |
| Dark mode | ✅ | Toggle persisted |
| Mobile responsiveness | ✅ | Responsive nav, forms, cards, tables, overflow |
| Accessibility pass | ✅ | Dialog semantics/focus, custom control labels, Escape handling, aria states |

---

## Phase 4: CI/CD, Deployment & Open Source

| Item | Status | Notes |
|------|--------|-------|
| CI pipeline | ✅ | Backend tests, Go formatting check, frontend lint/build, Docker build checks |
| Deploy workflow | ✅ | Runs after CI success on `main` and triggers Render deploy hook |
| Vercel Git deployment | ✅ | Frontend deploys directly via Vercel Git integration |
| Render backend deployment | ✅ | Backend deployed with Docker and health check |
| Neon production DB | ✅ | Configured and migrated |
| Production CORS | ✅ | Render backend allows production frontend origins, including approved subdomain setup |
| Secret hygiene check | ✅ | Current tree checked for obvious committed secrets; production secrets rotated by user |
| MIT License | ✅ | Project is now open source licensed |
| Deployment docs | ✅ | Root/backend/frontend READMEs updated |
| Branch protection | ✅ | Require PR + passing CI before merge to `main` |
| Manual migration workflow | ✅ | GitHub Action to run production migrations with a protected secret |
| DB-aware health check | ✅ | Health check verifies database connectivity, returns `degraded` status on failure |
| Secret scanning automation | ⬜ | Add GitHub secret scanning/gitleaks workflow |
| Market cache production migration | 🚧 | PR #4 adds `000007_market_cache`; apply via migration workflow before/with backend deploy |

---

## Phase 5: Future Enhancements (v2+)

| Item | Status | Notes |
|------|--------|-------|
| Market cache status badge | ⬜ | Show Live/Cached/Stale status in Markets UI |
| Market data error states | ⬜ | Per-card retry/unavailable state when provider and cache both fail |
| Stock watchlist | ⬜ | Track IDX tickers without requiring stock account ownership |
| Cryptocurrency portfolio | ⬜ | |
| Multi-currency support | ⬜ | |
| Recurring transactions | ⬜ | |
| Bank API integration (open banking) | ⬜ | |
| Import CSV/Excel | ⬜ | |
| Export reports to PDF | ⬜ | |
| Shared accounts (family mode) | ⬜ | |
| React Native mobile app | ⬜ | |
| Audit log / account activity history | ⬜ | |
| User profile and settings page | ⬜ | |

---

## Recommended Next Priorities

1. Merge PR #4 and run production migrations `000006_stock_holdings` and `000007_market_cache`.
2. Verify Markets page in production after backend restart: IHSG, stock charts, timestamps, and cache fallback.
3. Add market cache status badge (`Live` / `Cached` / `Stale`) and clearer per-card error states.
4. Add secret scanning automation (GitHub secret scanning or gitleaks).
5. Audit log / account activity history.
6. User profile and settings page.
7. Recurring transactions.
