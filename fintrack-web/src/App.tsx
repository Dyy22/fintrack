import { lazy, Suspense } from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import { AppLayout } from "./components/layout/AppLayout";
import { ErrorBoundary } from "./components/common/ErrorBoundary";
import { ProtectedRoute } from "./components/layout/ProtectedRoute";
import { AccountsPage } from "./pages/AccountsPage";
import { DashboardPage } from "./pages/DashboardPage";
import { LoginPage } from "./pages/LoginPage";
import { NewAccountPage } from "./pages/NewAccountPage";
import { NewTransactionPage } from "./pages/NewTransactionPage";
import { RegisterPage } from "./pages/RegisterPage";
import { ReportsPage } from "./pages/ReportsPage";
import { TransactionsPage } from "./pages/TransactionsPage";

const ApiDocsPage = lazy(() =>
  import("./pages/ApiDocsPage").then((module) => ({
    default: module.ApiDocsPage,
  })),
);

function DocsFallback() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-[#f7f1e3] px-6 text-center text-slate-950 dark:bg-slate-950 dark:text-slate-100">
      <div className="rounded-2xl border-2 border-slate-950 bg-white p-6 font-black uppercase shadow-[6px_6px_0_0_#101828] dark:border-slate-100 dark:bg-slate-900 dark:shadow-[6px_6px_0_0_#f8fafc]">
        Loading API docs...
      </div>
    </div>
  );
}

export default function App() {
  return (
    <ErrorBoundary>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route
          path="/docs/api"
          element={
            <Suspense fallback={<DocsFallback />}>
              <ApiDocsPage />
            </Suspense>
          }
        />
        <Route element={<ProtectedRoute />}>
          <Route element={<AppLayout />}>
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/accounts" element={<AccountsPage />} />
            <Route path="/accounts/new" element={<NewAccountPage />} />
            <Route path="/transactions" element={<TransactionsPage />} />
            <Route path="/transactions/new" element={<NewTransactionPage />} />
            <Route path="/reports" element={<ReportsPage />} />
          </Route>
        </Route>
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </ErrorBoundary>
  );
}
