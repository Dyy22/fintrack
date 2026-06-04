# Fintrack — Project Roadmap

**Version**: 1.1  
**Last Updated**: 4 June 2026

---

## Legend

- ✅ Done
- 🚧 In Progress
- ⬜ Planned

---

## Phase 0: Specification

| Item | Status |
|------|--------|
| Product & Technical Specification | ✅ |
| UI/UX Design Specification | ✅ |

---

## Phase 1: Backend MVP

### Core API

| Endpoint | Status | Notes |
|----------|--------|-------|
| `GET /health` | ✅ | |
| `POST /auth/register` | ✅ | |
| `POST /auth/login` | ✅ | |
| `GET /account-types` | ✅ | |
| `GET /accounts` | ✅ | |
| `POST /accounts` | ✅ | |
| `PUT /accounts/:id` | ✅ | |
| `DELETE /accounts/:id` | ✅ | Supports `?hard=true` for permanent delete |
| `GET /categories` | ✅ | |
| `POST /categories` | ✅ | |
| `PUT /categories/:id` | ✅ | |
| `DELETE /categories/:id` | ✅ | |
| `GET /transactions` | ✅ | With pagination (`limit`/`offset`) + filters |
| `POST /transactions` | ✅ | |
| `GET /reports/net-worth` | ✅ | |
| `GET /reports/spending-by-category` | ✅ | |
| `GET /gold/price` | ✅ | Current cached Antam gold price |
| `GET /gold/prices/history` | ✅ | 7-day gold price history for dashboard chart |

### Database

| Item | Status | Notes |
|------|--------|-------|
| Schema: users, account_types, accounts, categories, transactions | ✅ | |
| Seed: account_types (5 types) | ✅ | bank, ewallet, cash, gold, stock_broker |
| Seed: default categories (8 items) | ✅ | Salary, Freelance, Investment, Food, Transport, Shopping, Bills, Entertainment |
| Constraints: amount > 0, transfer rules, unique categories | ✅ | |
| Migration files (up/down) | ✅ | |

### Testing

| Item | Status | Notes |
|------|--------|-------|
| Unit tests: config | ✅ | |
| Unit tests: middleware auth | ✅ | |
| Unit tests: response helper | ✅ | |
| Unit tests: password + JWT security | ✅ | |
| Unit tests: usecase validation | ✅ | |
| Unit tests: handler/router | ✅ | |
| Integration tests: PostgreSQL repository | ✅ | Uses build tag `integration` |
| Smoke test script | ✅ | `scripts/smoke-test.sh` |

### DevOps

| Item | Status | Notes |
|------|--------|-------|
| Dockerfile (multi-stage) | ✅ | |
| Docker Compose (dev) | ✅ | |
| Docker Compose (root) | ✅ | Includes postgres, backend, web |
| Makefile | ✅ | dev, test, smoke-test, db migration, etc. |
| `.env.example` | ✅ | |
| Idempotent DB migration script | ✅ | `scripts/db-migrate.sh` |
| Insufficient balance validation | ✅ | 422 response with `insufficient_balance` error |

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
| App layout (sidebar desktop, bottom nav mobile) | ✅ | |
| Protected route wrapper | ✅ | Redirect to /login if unauthenticated |
| Reusable components: Button, Card, ConfirmDialog, ErrorBoundary, Skeleton | ✅ | |
| cn() utility (clsx + tailwind-merge) | ✅ | |
| Path alias `@/` | ✅ | |
| Formatters: formatIDR, formatDate, transactionAmountLabel | ✅ | |
| API error parser (validation_error mapping) | ✅ | |
| Page titles (usePageTitle hook) | ✅ | |
| Favicon (SVG) | ✅ | |
| Dark mode toggle | ✅ | localStorage + system preference detect |
| Dockerfile (multi-stage, nginx) | ✅ | |
| Git repositories (3 separate) | ✅ | root, backend, web on GitHub |

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
| Recent transactions (last 5) | ✅ | |
| Loading state (skeleton cards) | ✅ | |
| Empty state (no accounts) | ✅ | |

### Accounts

| Item | Status | Notes |
|------|--------|-------|
| Account list (desktop table, mobile cards) | ✅ | |
| Total balance card | ✅ | |
| Add account form | ✅ | Name, type dropdown, initial balance |
| Edit account name (inline) | ✅ | |
| Deactivate account | ✅ | With confirm dialog |
| Activate account | ✅ | Reactivate deactivated accounts |
| Delete account (permanent) | ✅ | Hard delete with confirm dialog |

### Transactions

| Item | Status | Notes |
|------|--------|-------|
| Transaction list | ✅ | Desktop table + mobile cards |
| Transaction filters | ✅ | Date range, account, type |
| Pagination controls | ✅ | Prev / Next with limit |
| Add transaction form (expense/income) | ✅ | Account, category, amount, date, description |
| Add transaction form (transfer) | ✅ | Source + destination account |
| Create transaction → redirect to list | ✅ | |
| Client-side validation | ✅ | Required fields checked before submit |

### Reports

| Item | Status | Notes |
|------|--------|-------|
| Net worth report | ✅ | Total + per account breakdown |
| Spending by category with date range | ✅ | Date filter with Apply button |
| Income vs spending donut chart on dashboard | ✅ | SVG donut chart showing income/spending split |
| Antam gold price trend chart | ✅ | Dashboard mini line chart for the latest 7 daily gold price snapshots |
| Simple horizontal bars | ✅ | Width = percentage of total |
| Percentage labels | ✅ | |

### UI Design — Neobrutalism

| Item | Status | Notes |
|------|--------|-------|
| Tailwind config — neo tokens (colors, shadow, radius, weight) | ✅ | |
| Global CSS — dot grid background, neobrutal form controls, `!important` overrides | ✅ | |
| Neobrutalism component system | ✅ | NeoButton, NeoCard, NeoInput, NeoTextarea, NeoSelect, NeoDateInput, NeoBadge, NeoAlert, NeoProgress, NeoTable, NeoEmptyState, NeoStatCard, NeoPageHeader |
| Refactor: Button, Card, Dialog, Skeleton → neo tokens | ✅ | |
| Refactor: Login, Register → NeoInput/NeoAlert/NeoCard | ✅ | |
| Refactor: NewAccount, NewTransaction → NeoInput/NeoTextarea/NeoSelect/NeoDateInput/NeoAlert/NeoPageHeader | ✅ | |
| Refactor: Accounts → NeoTable/NeoBadge/NeoStatCard/NeoPageHeader/NeoEmptyState | ✅ | |
| Refactor: Transactions → NeoTable/NeoBadge/NeoDateInput/NeoSelect/NeoPageHeader/NeoEmptyState | ✅ | |
| Refactor: Dashboard → NeoStatCard/NeoPageHeader/NeoEmptyState | ✅ | |
| Refactor: Reports → NeoStatCard/NeoProgress/NeoPageHeader/NeoEmptyState | ✅ | |
| Calendat picker — custom NeoDateInput (neobrutal popup, no native date input) | ✅ | |
| Dropdown list — custom NeoSelect (neobrutal popup, no native select) | ✅ | |
| Future date disabled in date picker and form | ✅ | `max=today`, clamp + validation | |
| Dark mode contrast fix (text-slate muted, border overrides) | ✅ | |

### Polish

| Item | Status | Notes |
|------|--------|-------|
| Empty states (all pages) | ✅ | NeoEmptyState |
| Loading states (all pages) | ✅ | Skeleton cards/rows |
| Skeleton loading | ✅ | Skeleton, SkeletonCard, SkeletonTable |
| Error boundaries | ✅ | ErrorBoundary component with Try Again |
| Client-side validation errors | ✅ | Field-level errors shown before submit |
| Dark mode | ✅ | Toggle in sidebar/header, persisted |
| Mobile responsiveness | ✅ | Responsive spacing, mobile nav, forms, cards, tables, and overflow pass |
| Accessibility pass | ✅ | Dialog semantics/focus, custom control labels, Escape handling, aria states, decorative icon cleanup |

---

## Phase 4: Future Enhancements (v2+)

| Item | Status | Notes |
|------|--------|-------|
| Gold gram tracking + Antam price display | ✅ | Gold accounts can store grams, backend values holdings using cached latest configured Antam price, dashboard shows current price |
| Gold price history snapshots + weekly chart | ✅ | Stores daily gold price snapshots and renders a 7-day trend chart on the dashboard |
| Cryptocurrency portfolio | ⬜ | |
| Multi-currency support | ⬜ | |
| Recurring transactions | ⬜ | |
| Budget tracking per category | ⬜ | |
| Bank API integration (open banking) | ⬜ | |
| Import CSV/Excel | ⬜ | |
| Export reports to PDF | ⬜ | |
| Shared accounts (family mode) | ⬜ | |
| React Native mobile app | ⬜ | |
| OpenAPI/Swagger docs | ⬜ | |
| CI pipeline (GitHub Actions) | ✅ | Runs backend tests, Go formatting check, frontend lint/build, and Docker build checks |
| Structured logging | ⬜ | |
