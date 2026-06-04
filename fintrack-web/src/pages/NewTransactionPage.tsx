import { FormEvent, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button } from "../components/common/Button";
import { Card } from "../components/common/Card";
import { NeoAlert } from "../components/common/NeoAlert";
import { NeoDateInput } from "../components/common/NeoDateInput";
import { NeoInput } from "../components/common/NeoInput";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { NeoSelect } from "../components/common/NeoSelect";
import { NeoTextarea } from "../components/common/NeoTextarea";
import { useAccountStore } from "../stores/accountStore";
import { useCategoryStore } from "../stores/categoryStore";
import { useTransactionStore } from "../stores/transactionStore";
import { useReportStore } from "../stores/reportStore";
import {
  getErrorMessage,
  getValidationErrors,
  type FormErrors,
} from "../utils/apiError";
import { getTodayDateInputValue } from "../utils/dateInput";
import { formatIDR } from "../utils/format";

type TransactionType = "expense" | "income" | "transfer";

import { usePageTitle } from "../utils/usePageTitle";

export function NewTransactionPage() {
  usePageTitle("Add Transaction");
  const navigate = useNavigate();
  const { accounts, fetchAccounts } = useAccountStore();
  const { categories, fetchCategories } = useCategoryStore();
  const { createTransaction, isLoading } = useTransactionStore();
  const { goldPrice, fetchGoldPrice } = useReportStore();

  const [type, setType] = useState<TransactionType>("expense");
  const [accountID, setAccountID] = useState("");
  const [transferAccountID, setTransferAccountID] = useState("");
  const [categoryID, setCategoryID] = useState("");
  const [amount, setAmount] = useState("");
  const [goldGrams, setGoldGrams] = useState("");
  const [goldInputMode, setGoldInputMode] = useState<"idr" | "grams" | null>(
    null,
  );
  const [date, setDate] = useState(() => getTodayDateInputValue());
  const [description, setDescription] = useState("");
  const maxDate = getTodayDateInputValue();
  const [fieldErrors, setFieldErrors] = useState<FormErrors>({});
  const [formError, setFormError] = useState("");

  useEffect(() => {
    fetchAccounts();
    fetchCategories();
    fetchGoldPrice().catch(() => undefined);
  }, [fetchAccounts, fetchCategories, fetchGoldPrice]);

  const activeAccounts = accounts.filter((account) => account.is_active);
  const accountOptions = [
    { value: "", label: "Select account" },
    ...activeAccounts.map((account) => ({
      value: account.id,
      label: account.name,
    })),
  ];
  const transferAccountOptions = [
    { value: "", label: "Select destination account" },
    ...activeAccounts
      .filter((account) => account.id !== accountID)
      .map((account) => ({ value: account.id, label: account.name })),
  ];
  const selectedAccount = activeAccounts.find(
    (account) => account.id === accountID,
  );
  const selectedTransferAccount = activeAccounts.find(
    (account) => account.id === transferAccountID,
  );
  const isGoldTransaction =
    selectedAccount?.type === "gold" ||
    selectedTransferAccount?.type === "gold";

  const filteredCategories = categories.filter(
    (c) => c.type === type || type === "transfer",
  );
  const categoryOptions = [
    { value: "", label: "Select category" },
    ...filteredCategories.map((category) => ({
      value: category.id,
      label: category.name,
    })),
  ];

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFieldErrors({});
    setFormError("");

    const errors: FormErrors = {};
    if (!accountID) errors.account_id = "is required";
    if (!amount || Number(amount) <= 0)
      errors.amount = "must be greater than 0";
    if (isGoldTransaction) {
      if (!goldPrice?.price_per_gram) {
        errors.gold_grams = "gold price is not available";
      }
      if (!goldGrams || Number(goldGrams) <= 0) {
        errors.gold_grams = "must be greater than 0";
      }
    }
    if (date > maxDate) errors.date = "cannot be in the future";
    if (type === "transfer") {
      if (!transferAccountID) errors.transfer_account_id = "is required";
    } else {
      if (!categoryID) errors.category_id = "is required";
    }

    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors);
      return;
    }

    const payload: Record<string, unknown> = {
      type,
      account_id: accountID,
      amount: Number(amount),
      date: `${date}T12:00:00Z`,
      description,
    };

    if (isGoldTransaction) {
      payload.gold_grams = Number(goldGrams);
    }

    if (type === "transfer") {
      payload.transfer_account_id = transferAccountID;
    } else {
      payload.category_id = categoryID;
    }

    try {
      await createTransaction(payload);
      navigate("/transactions", { replace: true });
    } catch (error) {
      setFieldErrors(getValidationErrors(error));
      setFormError(getErrorMessage(error, "Unable to create transaction."));
    }
  }

  const typeOptions: { key: TransactionType; label: string }[] = [
    { key: "expense", label: "Expense" },
    { key: "income", label: "Income" },
    { key: "transfer", label: "Transfer" },
  ];

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <NeoPageHeader
        title="Add Transaction"
        description="Create income, expense, or transfer records."
        eyebrow="New money movement"
        icon="➕"
      />

      <Card>
        {formError ? (
          <NeoAlert className="mb-4" variant="danger">
            {formError}
          </NeoAlert>
        ) : null}

        <form className="space-y-4" onSubmit={handleSubmit}>
          <div className="grid gap-2 sm:grid-cols-3">
            {typeOptions.map((opt) => (
              <Button
                key={opt.key}
                type="button"
                variant={type === opt.key ? "primary" : "secondary"}
                aria-pressed={type === opt.key}
                onClick={() => {
                  setType(opt.key);
                  setCategoryID("");
                  setTransferAccountID("");
                  setAmount("");
                  setGoldGrams("");
                  setGoldInputMode(null);
                }}
              >
                {opt.label}
              </Button>
            ))}
          </div>

          <label className="block">
            <span className="text-sm font-medium text-slate-700">
              {type === "transfer" ? "Source Account" : "Account"}
            </span>
            <NeoSelect
              className="mt-1"
              value={accountID}
              options={accountOptions}
              onChange={(value) => {
                setAccountID(value);
                if (transferAccountID === value) setTransferAccountID("");
                setAmount("");
                setGoldGrams("");
                setGoldInputMode(null);
              }}
              placeholder="Select account"
              ariaLabel={type === "transfer" ? "Source account" : "Account"}
            />
            {fieldErrors.account_id ? (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.account_id}
              </span>
            ) : null}
          </label>

          {type === "transfer" ? (
            <label className="block">
              <span className="text-sm font-medium text-slate-700">
                Destination Account
              </span>
              <NeoSelect
                className="mt-1"
                value={transferAccountID}
                options={transferAccountOptions}
                onChange={(value) => {
                  setTransferAccountID(value);
                  setAmount("");
                  setGoldGrams("");
                  setGoldInputMode(null);
                }}
                placeholder="Select destination account"
                ariaLabel="Destination account"
              />
              {fieldErrors.transfer_account_id ? (
                <span className="mt-1 block text-sm text-red-600">
                  {fieldErrors.transfer_account_id}
                </span>
              ) : null}
            </label>
          ) : null}

          {type !== "transfer" ? (
            <label className="block">
              <span className="text-sm font-medium text-slate-700">
                Category
              </span>
              <NeoSelect
                className="mt-1"
                value={categoryID}
                options={categoryOptions}
                onChange={(value) => setCategoryID(value)}
                placeholder="Select category"
                ariaLabel="Category"
              />
              {fieldErrors.category_id ? (
                <span className="mt-1 block text-sm text-red-600">
                  {fieldErrors.category_id}
                </span>
              ) : null}
            </label>
          ) : null}

          {isGoldTransaction ? (
            <div className="grid gap-4 sm:grid-cols-2">
              <label className="block">
                <span className="text-sm font-medium text-slate-700">
                  Nominal IDR
                </span>
                <NeoInput
                  type="number"
                  min="0"
                  value={amount}
                  disabled={goldInputMode === "grams"}
                  onChange={(e) => {
                    const value = e.target.value;
                    setAmount(value);
                    setGoldInputMode(value ? "idr" : null);
                    if (goldPrice?.price_per_gram && value) {
                      setGoldGrams(
                        String(Number(value) / goldPrice.price_per_gram),
                      );
                    } else if (!value) {
                      setGoldGrams("");
                    }
                  }}
                  placeholder="2759000"
                />
                {fieldErrors.amount ? (
                  <span className="mt-1 block text-sm text-red-600">
                    {fieldErrors.amount}
                  </span>
                ) : null}
              </label>

              <label className="block">
                <span className="text-sm font-medium text-slate-700">
                  Nominal Gram
                </span>
                <NeoInput
                  type="number"
                  min="0"
                  step="0.0001"
                  value={goldGrams}
                  disabled={goldInputMode === "idr"}
                  onChange={(e) => {
                    const value = e.target.value;
                    setGoldGrams(value);
                    setGoldInputMode(value ? "grams" : null);
                    if (goldPrice?.price_per_gram && value) {
                      setAmount(
                        String(
                          Math.round(Number(value) * goldPrice.price_per_gram),
                        ),
                      );
                    } else if (!value) {
                      setAmount("");
                    }
                  }}
                  placeholder="1"
                />
                {fieldErrors.gold_grams ? (
                  <span className="mt-1 block text-sm text-red-600">
                    {fieldErrors.gold_grams}
                  </span>
                ) : null}
              </label>
              <p className="text-xs font-semibold text-slate-500 sm:col-span-2">
                Current Antam price: {formatIDR(goldPrice?.price_per_gram)} /
                gr. Fill one field and Fintrack will calculate the other.
              </p>
            </div>
          ) : (
            <label className="block">
              <span className="text-sm font-medium text-slate-700">Amount</span>
              <NeoInput
                type="number"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                placeholder="50000"
              />
              {fieldErrors.amount ? (
                <span className="mt-1 block text-sm text-red-600">
                  {fieldErrors.amount}
                </span>
              ) : null}
            </label>
          )}

          <label className="block">
            <span className="text-sm font-medium text-slate-700">Date</span>
            <NeoDateInput
              className="mt-1"
              value={date}
              max={maxDate}
              onChange={(value) => setDate(value)}
              ariaLabel="Transaction date"
            />
            {fieldErrors.date ? (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.date}
              </span>
            ) : null}
          </label>

          <label className="block">
            <span className="text-sm font-medium text-slate-700">
              Description
            </span>
            <NeoTextarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Lunch"
              rows={3}
            />
            {fieldErrors.description ? (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.description}
              </span>
            ) : null}
          </label>

          <div className="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
            <Link to="/transactions" className="sm:inline-flex">
              <Button
                className="w-full sm:w-auto"
                variant="secondary"
                type="button"
              >
                Cancel
              </Button>
            </Link>
            <Button
              className="w-full sm:w-auto"
              type="submit"
              disabled={isLoading}
            >
              {isLoading ? "Saving..." : "Save Transaction"}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
