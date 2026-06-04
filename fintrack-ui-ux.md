# Fintrack - UI/UX Design Specification

**Version**: 1.0  
**Date**: June 2026  
**Status**: Draft  
**Scope**: Web frontend MVP

---

## 1. Design Goals

Fintrack's web UI should feel simple, private, and practical. The main goal is to help users understand their financial position quickly and record transactions with minimal friction.

### Primary Goals
- Make net worth and account balances visible at a glance
- Make adding accounts and transactions fast
- Keep reports understandable without complex financial jargon
- Work comfortably on desktop and mobile web
- Avoid visual clutter; prioritize clear information hierarchy

### Product Personality
- Clean
- Calm
- Trustworthy
- Minimal
- Data-focused

---

## 2. Target User Context

Users are individuals tracking personal assets manually across:

- Bank accounts
- E-wallets
- Cash
- Gold as manually valued IDR balance
- Stock broker account as manually valued IDR balance

They likely update data manually and want quick answers:

- How much do I have overall?
- Where is my money stored?
- What did I spend this month?
- What are my recent transactions?

---

## 3. Information Architecture

### Primary Navigation

Authenticated app navigation:

```txt
Dashboard
Accounts
Transactions
Reports
```

Optional future navigation:

```txt
Settings
Profile
Budgets
Imports
Exports
```

### Public Pages

```txt
/login
/register
```

### Authenticated Pages

```txt
/dashboard
/accounts
/transactions
/transactions/new
/reports
```

Default redirect after login:

```txt
/dashboard
```

---

## 4. User Flows

### 4.1 First-Time User Flow

```txt
Register
  → Login automatically or redirect to Login
  → Dashboard empty state
  → Add first account
  → Add first transaction
  → Dashboard updates
```

### 4.2 Returning User Flow

```txt
Login
  → Dashboard
  → Review net worth and recent transactions
  → Add transaction if needed
```

### 4.3 Add Account Flow

```txt
Accounts page
  → Click Add Account
  → Fill name, account type, initial balance
  → Submit
  → Account appears in list
  → Dashboard net worth updates
```

### 4.4 Add Expense/Income Flow

```txt
Transactions page or Dashboard quick action
  → Click Add Transaction
  → Select type: Expense or Income
  → Select account
  → Select category
  → Fill amount, date, description
  → Submit
  → Transaction appears in history
  → Account balance updates
```

### 4.5 Add Transfer Flow

```txt
Transactions page
  → Click Add Transaction
  → Select type: Transfer
  → Select source account
  → Select destination account
  → Fill amount, date, description
  → Submit
  → Source balance decreases
  → Destination balance increases
```

### 4.6 Monthly Review Flow

```txt
Reports page
  → Select date range or month
  → View spending by category
  → Review top spending categories
```

---

## 5. Layout System

### Desktop Layout

Desktop uses a sidebar layout.

```txt
+-------------------------------------------------------------+
| Sidebar       | Topbar                                      |
|               |---------------------------------------------|
| Dashboard     | Page Content                                |
| Accounts      |                                             |
| Transactions  |                                             |
| Reports       |                                             |
|               |                                             |
| Logout        |                                             |
+-------------------------------------------------------------+
```

### Mobile Layout

Mobile uses a top header and bottom navigation.

```txt
+-----------------------------+
| Header / Page Title         |
|-----------------------------|
| Page Content                |
|                             |
|                             |
|-----------------------------|
| Dashboard Accounts Tx Report|
+-----------------------------+
```

### Breakpoint Recommendation

```txt
Mobile: < 768px
Tablet: 768px - 1023px
Desktop: >= 1024px
```

---

## 6. Visual Direction

### Color Direction

Recommended palette direction:

- Background: soft neutral gray / near-white
- Surface/cards: white
- Primary: blue or teal
- Success/income: green
- Expense/danger: red
- Transfer/info: blue/purple
- Text primary: near-black
- Text secondary: gray

Example semantic usage:

```txt
Income      → green
Expense     → red
Transfer    → blue
Inactive    → gray
Warning     → amber
```

### Typography

Use a modern sans-serif font. Recommended:

```txt
Inter, system-ui, sans-serif
```

### UI Style

- Rounded cards
- Subtle shadows or border-only cards
- Clear spacing
- Minimal icons
- Table on desktop, card list on mobile

### Number Formatting

Amounts should be displayed as Indonesian Rupiah:

```txt
Rp 5.000.000
Rp 50.000
```

Negative/expense amounts:

```txt
- Rp 50.000
```

Income amounts:

```txt
+ Rp 5.000.000
```

---

## 7. Page Designs

## 7.1 Login Page

### Purpose
Allow existing users to access the app.

### Content
- App name/logo
- Email field
- Password field
- Login button
- Link to register
- Error message area

### Wireframe

```txt
+-----------------------------------+
|           Fintrack                |
|  Private personal finance tracker |
|                                   |
|  Email                            |
|  [ user@example.com             ] |
|                                   |
|  Password                         |
|  [ ********                     ] |
|                                   |
|  [ Login ]                        |
|                                   |
|  Don't have an account? Register  |
+-----------------------------------+
```

### States
- Loading while logging in
- Invalid credentials
- Validation error

---

## 7.2 Register Page

### Purpose
Create a new account.

### Content
- Email field
- Password field
- Register button
- Link to login

### Wireframe

```txt
+-----------------------------------+
|           Create Account          |
|                                   |
|  Email                            |
|  [ user@example.com             ] |
|                                   |
|  Password                         |
|  [ minimum 8 characters         ] |
|                                   |
|  [ Register ]                     |
|                                   |
|  Already have an account? Login   |
+-----------------------------------+
```

---

## 7.3 Dashboard Page

### Purpose
Give an at-a-glance summary of current finances.

### Content
- Net worth card
- Account summary cards
- Recent transactions
- Spending this month
- Quick action buttons

### Desktop Wireframe

```txt
+-------------------------------------------------------------+
| Dashboard                                      [+ Add Tx]    |
|-------------------------------------------------------------|
| Net Worth                                                   |
| Rp 5.450.000                                                |
|-------------------------------------------------------------|
| [ Accounts: 2 ] [ Spending This Month: Rp 50.000 ]          |
|-------------------------------------------------------------|
| Account Balances                                            |
| +--------------------+ +--------------------+               |
| | BCA Savings        | | Cash               |               |
| | Rp 4.850.000       | | Rp 600.000         |               |
| +--------------------+ +--------------------+               |
|-------------------------------------------------------------|
| Recent Transactions                                         |
| Date        Type      Account       Category      Amount    |
| 03 Jun      Expense   BCA Savings   Food       -Rp 50.000   |
| 03 Jun      Transfer  BCA Savings   Transfer   Rp 100.000   |
+-------------------------------------------------------------+
```

### Mobile Wireframe

```txt
+-----------------------------+
| Dashboard           + Tx    |
|-----------------------------|
| Net Worth                   |
| Rp 5.450.000                |
|-----------------------------|
| BCA Savings                 |
| Rp 4.850.000                |
|-----------------------------|
| Cash                        |
| Rp 600.000                  |
|-----------------------------|
| Recent Transactions         |
| Food              -Rp 50k   |
| Transfer           Rp 100k  |
+-----------------------------+
```

### Empty State

If user has no accounts:

```txt
No accounts yet
Add your first account to start tracking your net worth.
[ Add Account ]
```

---

## 7.4 Accounts Page

### Purpose
Manage accounts and balances.

### Content
- Account type filter optional
- Add account button
- Account list
- Edit action
- Deactivate/soft delete action

### Wireframe

```txt
+-------------------------------------------------------------+
| Accounts                                      [+ Add]        |
|-------------------------------------------------------------|
| Total Balance: Rp 5.450.000                                 |
|-------------------------------------------------------------|
| Name          Type       Balance        Status     Actions   |
| BCA Savings   Bank       Rp 4.850.000   Active     Edit      |
| Cash          Cash       Rp 600.000     Active     Edit      |
+-------------------------------------------------------------+
```

### Add/Edit Account Form

Fields:
- Name
- Account type
- Initial balance (create only)
- Active status (edit only)

### Mobile Behavior
Accounts appear as cards:

```txt
+-----------------------------+
| BCA Savings                 |
| Bank                        |
| Rp 4.850.000                |
| [Edit]                      |
+-----------------------------+
```

---

## 7.5 Transactions Page

### Purpose
View and create transactions.

### Content
- Date range filters
- Account filter
- Category filter
- Type filter
- Pagination
- Transaction list/table
- Add transaction button

### Desktop Wireframe

```txt
+-------------------------------------------------------------+
| Transactions                                  [+ Add]        |
|-------------------------------------------------------------|
| Start Date | End Date | Account | Category | Type | Filter  |
|-------------------------------------------------------------|
| Date        Type      Account       Category      Amount     |
| 03 Jun      Expense   BCA Savings   Food       -Rp 50.000    |
| 03 Jun      Transfer  BCA Savings   -           Rp 100.000   |
|-------------------------------------------------------------|
| [Previous]                         Showing 1-50    [Next]    |
+-------------------------------------------------------------+
```

### Mobile Wireframe

```txt
+-----------------------------+
| Transactions          +     |
|-----------------------------|
| [Filters]                   |
|-----------------------------|
| Food                        |
| BCA Savings • 03 Jun        |
| -Rp 50.000                  |
|-----------------------------|
| Transfer                    |
| BCA Savings → Cash          |
| Rp 100.000                  |
+-----------------------------+
```

---

## 7.6 Add Transaction Page / Modal

### Recommendation
Use a dedicated page on mobile and a modal/drawer on desktop.

Route:

```txt
/transactions/new
```

### Transaction Type Selector

```txt
[ Expense ] [ Income ] [ Transfer ]
```

### Expense/Income Form

Fields:
- Type
- Account
- Category
- Amount
- Date
- Description

### Transfer Form

Fields:
- Source account
- Destination account
- Amount
- Date
- Description

### Wireframe

```txt
+-----------------------------------+
| Add Transaction                   |
|-----------------------------------|
| Type                              |
| [ Expense ] [ Income ] [Transfer] |
|                                   |
| Account                           |
| [ BCA Savings v ]                 |
|                                   |
| Category                          |
| [ Food v ]                        |
|                                   |
| Amount                            |
| [ 50000                         ] |
|                                   |
| Date                              |
| [ 2026-06-03                    ] |
|                                   |
| Description                       |
| [ Lunch                         ] |
|                                   |
| [ Cancel ] [ Save ]               |
+-----------------------------------+
```

### Validation
- Amount must be greater than zero
- Expense/income require category
- Transfer requires different source and destination accounts
- Show backend validation errors below fields

---

## 7.7 Reports Page

### Purpose
Help users understand spending patterns.

### Content
- Date range/month filter
- Spending by category summary
- Category list with percentage
- Net worth summary optional

### Wireframe

```txt
+-------------------------------------------------------------+
| Reports                                                     |
|-------------------------------------------------------------|
| Date Range: [2026-06-01] to [2026-06-30] [Apply]            |
|-------------------------------------------------------------|
| Total Spending                                              |
| Rp 50.000                                                   |
|-------------------------------------------------------------|
| Spending by Category                                        |
| Food          Rp 50.000        100%                         |
| Transport     Rp 0             0%                           |
+-------------------------------------------------------------+
```

### Chart Recommendation
For MVP:
- Horizontal bar chart is enough
- Pie chart optional later

---

## 8. Component Inventory

### Layout Components
- `AppLayout`
- `Sidebar`
- `MobileBottomNav`
- `Topbar`
- `PageHeader`

### Auth Components
- `AuthCard`
- `LoginForm`
- `RegisterForm`

### Common Components
- `Button`
- `Input`
- `Select`
- `DateInput`
- `AmountInput`
- `Card`
- `EmptyState`
- `LoadingState`
- `ErrorAlert`
- `PaginationControls`

### Finance Components
- `MoneyText`
- `NetWorthCard`
- `AccountCard`
- `AccountTable`
- `TransactionTable`
- `TransactionListItem`
- `TransactionTypeBadge`
- `CategoryBadge`
- `SpendingCategoryBar`

---

## 9. State Management Plan

Recommended Zustand stores:

### `authStore`

State:
- `token`
- `user`
- `isAuthenticated`

Actions:
- `login`
- `register`
- `logout`
- `restoreSession`

Persistence:
- token can be stored in `localStorage` for MVP

### `accountStore`

State:
- `accountTypes`
- `accounts`
- `isLoading`

Actions:
- `fetchAccountTypes`
- `fetchAccounts`
- `createAccount`
- `updateAccount`
- `deleteAccount`

### `categoryStore`

State:
- `categories`

Actions:
- `fetchCategories`
- `createCategory`
- `updateCategory`
- `deleteCategory`

### `transactionStore`

State:
- `transactions`
- `limit`
- `offset`
- `filters`

Actions:
- `fetchTransactions`
- `createTransaction`
- `setFilters`
- `nextPage`
- `previousPage`

### `reportStore`

State:
- `netWorth`
- `spendingByCategory`

Actions:
- `fetchNetWorth`
- `fetchSpendingByCategory`

---

## 10. API Client Plan

Use Axios with a shared client:

```txt
src/services/api.ts
```

Responsibilities:
- Set `baseURL`
- Attach JWT token to requests
- Handle 401 globally
- Normalize backend validation errors for forms

Example backend validation shape:

```json
{
  "error": "validation_error",
  "message": "invalid request body",
  "fields": {
    "email": "must be a valid email"
  }
}
```

Frontend should map `fields.email` to the email input error.

---

## 11. Responsive Behavior

### Desktop
- Sidebar visible
- Tables for accounts and transactions
- Cards in multi-column grid

### Mobile
- Bottom navigation
- Tables become cards
- Forms full-screen or page-based
- Filters collapse into expandable panel

### Tablet
- Sidebar may collapse
- Cards use 2-column layout

---

## 12. Accessibility Notes

- Use semantic buttons and labels
- Every input must have visible label
- Color should not be the only signal for income/expense
- Use readable contrast
- Loading states must not trap interaction
- Error messages should be associated with fields

---

## 13. Empty, Loading, and Error States

### Empty States

Accounts:
```txt
No accounts yet
Add your first account to start tracking your net worth.
[ Add Account ]
```

Transactions:
```txt
No transactions found
Try changing filters or add a new transaction.
[ Add Transaction ]
```

Reports:
```txt
No spending data for this period
Add expense transactions to see category breakdown.
```

### Loading States
- Skeleton cards for dashboard
- Spinner/button loading for form submit
- Table skeleton for transactions

### Error States
- Use inline form errors for validation
- Use top alert for API/network errors
- Provide retry action for failed fetches

---

## 14. Implementation Phases

### Phase 1: Frontend Scaffold
- Setup Vite + React + TypeScript
- Setup Tailwind CSS
- Setup React Router
- Setup Axios client
- Setup Zustand stores
- Create route placeholders
- Create app layout

### Phase 2: Authentication
- Login page
- Register page
- Auth store
- Protected routes
- Logout

### Phase 3: Dashboard + Accounts
- Dashboard net worth
- Account cards
- Accounts page
- Add/edit account form

### Phase 4: Categories + Transactions
- Fetch categories
- Transaction list with filters and pagination
- Add transaction form
- Transfer support

### Phase 5: Reports
- Net worth report
- Spending by category report
- Date range filter
- Simple horizontal bar visualization

### Phase 6: Polish
- Mobile responsiveness pass
- Empty states
- Loading states
- Error states
- Basic accessibility pass

---

## 15. MVP Screen Priority

Build screens in this order:

1. Login
2. Register
3. App layout / protected routes
4. Dashboard
5. Accounts
6. Transactions list
7. Add transaction
8. Reports

Reasoning:
- Auth is required before accessing API data
- Dashboard validates overall app structure
- Accounts are required before useful transactions
- Transactions are the core daily workflow
- Reports depend on transaction data

---

## 16. MVP Acceptance Criteria

### Authentication
- User can register
- User can login
- Token is persisted across refresh
- User can logout
- Protected pages redirect unauthenticated users to login

### Dashboard
- User can see net worth
- User can see account balances
- User can see recent transactions
- Empty state appears when there are no accounts

### Accounts
- User can view account types
- User can create an account
- User can edit account name/status
- User can deactivate an account

### Transactions
- User can view paginated transactions
- User can filter transactions by date range and type
- User can create income transaction
- User can create expense transaction
- User can create transfer transaction
- Backend validation errors appear below relevant fields

### Reports
- User can view net worth report
- User can view spending by category for a date range
- Empty state appears when there is no spending data

---

## 17. Key UI Decisions

### Use dedicated pages first, modals later
For MVP, prefer dedicated pages/forms over complex modals. This makes mobile behavior simpler and reduces implementation complexity.

### Use tables on desktop and cards on mobile
Transactions and accounts should be tables on desktop, but cards on mobile for readability.

### Keep charts simple
Use horizontal bars for category spending in MVP. Avoid heavy chart dependencies unless needed.

### Store JWT in localStorage for MVP
This is acceptable for the current self-hosted MVP. If security requirements increase later, revisit token storage and session strategy.

### Keep currency fixed to IDR
All amount inputs and displays assume IDR for v1.

---

## 18. Open UI Questions

- Should registration automatically log the user in, or redirect to login?
- Should add transaction be a page, drawer, or modal on desktop?
- Should dashboard show income this month in v1, or only spending?
- Should accounts page show inactive accounts by default?
- Should transaction filters auto-apply or require clicking an Apply button?

Recommended MVP decisions:
- Registration redirects to login
- Add transaction uses a dedicated page
- Dashboard shows net worth, accounts, recent transactions, and spending this month
- Accounts page hides inactive accounts by default but provides a toggle later
- Transaction filters use an Apply button to avoid excessive API calls

---

**Document End**

---

*This UI/UX design specification is a living document and should be updated as the frontend evolves.*