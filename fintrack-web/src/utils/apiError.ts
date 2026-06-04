import { AxiosError } from "axios";
import type { ApiErrorResponse, ValidationErrorResponse } from "../types";

export type FormErrors = Record<string, string>;

export function getValidationErrors(error: unknown): FormErrors {
  if (!isAxiosError(error)) return {};
  const data = error.response?.data;
  if (isValidationErrorResponse(data)) {
    return data.fields ?? {};
  }
  return {};
}

export function getErrorMessage(
  error: unknown,
  fallback = "Something went wrong",
): string {
  if (!isAxiosError(error)) return fallback;
  const data = error.response?.data;

  // Check for validation errors first
  if (isValidationErrorResponse(data)) {
    return data.message || fallback;
  }

  // Check for specific error codes
  if (isErrorCode(data, "insufficient_balance")) {
    return "Insufficient funds for this transaction";
  }

  // Generic API error
  if (isApiErrorResponse(data)) {
    return data.message || fallback;
  }

  if (error.response?.status === 401) {
    return "Invalid email or password";
  }

  return fallback;
}

function isErrorCode(data: unknown, code: string): boolean {
  return Boolean(
    data &&
    typeof data === "object" &&
    (data as { error?: unknown }).error === code,
  );
}

function isAxiosError(error: unknown): error is AxiosError<unknown> {
  return Boolean(error && typeof error === "object" && "isAxiosError" in error);
}

function isValidationErrorResponse(
  data: unknown,
): data is ValidationErrorResponse {
  return Boolean(
    data &&
    typeof data === "object" &&
    "error" in data &&
    (data as { error?: unknown }).error === "validation_error" &&
    "fields" in data,
  );
}

function isApiErrorResponse(data: unknown): data is ApiErrorResponse {
  return Boolean(
    data && typeof data === "object" && "error" in data && "message" in data,
  );
}
