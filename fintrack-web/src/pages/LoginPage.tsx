import { FormEvent, useState } from "react";
import { Link, Navigate, useLocation, useNavigate } from "react-router-dom";
import { Button } from "../components/common/Button";
import { NeoAlert } from "../components/common/NeoAlert";
import { NeoInput } from "../components/common/NeoInput";
import { useAuthStore } from "../stores/authStore";
import { useTheme } from "../stores/themeStore";
import {
  getErrorMessage,
  getValidationErrors,
  type FormErrors,
} from "../utils/apiError";
import { usePageTitle } from "../utils/usePageTitle";

type LocationState = {
  from?: {
    pathname?: string;
  };
  message?: string;
};

export function LoginPage() {
  usePageTitle("Login");
  const navigate = useNavigate();
  const location = useLocation();
  const locationState = location.state as LocationState | null;
  const redirectTo = locationState?.from?.pathname ?? "/dashboard";
  const login = useAuthStore((state) => state.login);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const isLoading = useAuthStore((state) => state.isLoading);
  const { isDark, toggle } = useTheme();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<FormErrors>({});
  const [formError, setFormError] = useState("");

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFieldErrors({});
    setFormError("");

    try {
      await login(email, password);
      navigate(redirectTo, { replace: true });
    } catch (error) {
      setFieldErrors(getValidationErrors(error));
      setFormError(
        getErrorMessage(error, "Unable to login. Please try again."),
      );
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center px-3 py-8 sm:px-4 sm:py-12">
      <section className="neo-card w-full max-w-md bg-blue-100 p-5 dark:bg-slate-800 sm:p-8">
        <div className="mb-8 text-center">
          <div className="flex items-center justify-between gap-3">
            <h1 className="text-3xl font-black uppercase text-slate-950 dark:text-slate-100 sm:text-4xl">
              Fintrack
            </h1>
            <button
              className="neo-button bg-blue-300 px-3 py-1 dark:text-slate-950"
              onClick={toggle}
              aria-label={
                isDark ? "Switch to light mode" : "Switch to dark mode"
              }
              aria-pressed={isDark}
            >
              {isDark ? "☀️" : "🌙"}
            </button>
          </div>
          <p className="mt-2 text-sm font-bold text-slate-600 dark:text-slate-300">
            Private personal finance tracker
          </p>
        </div>

        {locationState?.message ? (
          <NeoAlert className="mb-4" variant="success">
            {locationState.message}
          </NeoAlert>
        ) : null}
        {formError ? (
          <NeoAlert className="mb-4" variant="danger">
            {formError}
          </NeoAlert>
        ) : null}

        <form className="space-y-4" onSubmit={handleSubmit}>
          <label className="block">
            <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
              Email
            </span>
            <NeoInput
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="user@example.com"
              autoComplete="email"
            />
            {fieldErrors.email ? (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.email}
              </span>
            ) : null}
          </label>
          <label className="block">
            <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
              Password
            </span>
            <NeoInput
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="••••••••"
              autoComplete="current-password"
            />
            {fieldErrors.password ? (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.password}
              </span>
            ) : null}
          </label>
          <Button className="w-full" type="submit" disabled={isLoading}>
            {isLoading ? "Logging in..." : "Login"}
          </Button>
        </form>
        <p className="mt-6 text-center text-sm text-slate-600 dark:text-slate-400">
          Don't have an account?{" "}
          <Link className="neo-link" to="/register">
            Register
          </Link>
        </p>
      </section>
    </main>
  );
}
