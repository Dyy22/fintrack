import { FormEvent, useState } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";
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

export function RegisterPage() {
  usePageTitle("Create Account");
  const navigate = useNavigate();
  const register = useAuthStore((state) => state.register);
  const login = useAuthStore((state) => state.login);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const isLoading = useAuthStore((state) => state.isLoading);
  const { isDark, toggle } = useTheme();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<FormErrors>({});
  const [formError, setFormError] = useState("");

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFieldErrors({});
    setFormError("");

    if (password !== confirmPassword) {
      setFieldErrors({ confirm_password: "must match password" });
      return;
    }

    try {
      await register(email, password);
      await login(email, password);
      navigate("/dashboard", { replace: true });
    } catch (error) {
      setFieldErrors(getValidationErrors(error));
      setFormError(
        getErrorMessage(error, "Unable to register. Please try again."),
      );
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center px-3 py-8 sm:px-4 sm:py-12">
      <section className="neo-card w-full max-w-md bg-blue-100 p-5 dark:bg-slate-800 sm:p-8">
        <div className="mb-8 text-center">
          <div className="flex items-center justify-between gap-3">
            <h1 className="text-2xl font-black uppercase text-slate-950 dark:text-slate-100 sm:text-4xl">
              Create Account
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
            Start tracking your personal finances
          </p>
        </div>

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
              placeholder="minimum 8 characters"
              autoComplete="new-password"
            />
            {fieldErrors.password ? (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.password}
              </span>
            ) : null}
          </label>
          <label className="block">
            <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
              Confirm Password
            </span>
            <NeoInput
              type="password"
              value={confirmPassword}
              onChange={(event) => setConfirmPassword(event.target.value)}
              placeholder="repeat password"
              autoComplete="new-password"
            />
            {fieldErrors.confirm_password ? (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.confirm_password}
              </span>
            ) : null}
          </label>
          <Button className="w-full" type="submit" disabled={isLoading}>
            {isLoading ? "Creating account..." : "Register"}
          </Button>
        </form>
        <p className="mt-6 text-center text-sm text-slate-600 dark:text-slate-400">
          Already have an account?{" "}
          <Link className="neo-link" to="/login">
            Login
          </Link>
        </p>
      </section>
    </main>
  );
}
