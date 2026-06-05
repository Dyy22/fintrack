import { FormEvent, useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Button } from "../components/common/Button";
import { Card } from "../components/common/Card";
import { NeoAlert } from "../components/common/NeoAlert";
import { NeoInput } from "../components/common/NeoInput";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { NeoSelect } from "../components/common/NeoSelect";
import { useAccountStore } from "../stores/accountStore";
import { useReportStore } from "../stores/reportStore";
import {
  getErrorMessage,
  getValidationErrors,
  type FormErrors,
} from "../utils/apiError";

import { formatIDR } from "../utils/format";
import { usePageTitle } from "../utils/usePageTitle";

const COMMON_IDX_SYMBOLS = [
  "AALI",
  "ACES",
  "ADRO",
  "AKRA",
  "ANTM",
  "ARTO",
  "ASII",
  "BBCA",
  "BBNI",
  "BBRI",
  "BBTN",
  "BMRI",
  "BRIS",
  "BRPT",
  "BUKA",
  "CPIN",
  "EMTK",
  "ERAA",
  "EXCL",
  "GOTO",
  "ICBP",
  "INCO",
  "INDF",
  "INKP",
  "INTP",
  "ITMG",
  "KLBF",
  "MDKA",
  "MEDC",
  "PGAS",
  "PTBA",
  "SMGR",
  "TLKM",
  "TOWR",
  "UNTR",
  "UNVR",
];

export function NewAccountPage() {
  usePageTitle("Add Account");
  const navigate = useNavigate();
  const { accountTypes, fetchAccountTypes, createAccount, isLoading } =
    useAccountStore();
  const { goldPrice, fetchGoldPrice } = useReportStore();
  const [name, setName] = useState("");
  const [accountTypeID, setAccountTypeID] = useState("");
  const [balance, setBalance] = useState("");
  const [goldGrams, setGoldGrams] = useState("");
  const [stockLots, setStockLots] = useState("");
  const [goldInputMode, setGoldInputMode] = useState<"idr" | "grams" | null>(
    null,
  );
  const [fieldErrors, setFieldErrors] = useState<FormErrors>({});
  const [formError, setFormError] = useState("");

  useEffect(() => {
    fetchAccountTypes();
    fetchGoldPrice().catch(() => undefined);
  }, [fetchAccountTypes, fetchGoldPrice]);

  const selectedAccountType = accountTypes.find(
    (type) => String(type.id) === accountTypeID,
  );
  const isGoldAccount = selectedAccountType?.name === "gold";
  const isStockAccount = selectedAccountType?.name === "stock_broker";

  const accountTypeOptions = [
    { value: "", label: "Select account type" },
    ...accountTypes.map((type) => ({
      value: String(type.id),
      label:
        type.name === "ewallet"
          ? "E-Wallet"
          : type.name
              .replace(/_/g, " ")
              .replace(/\b\w/g, (character) => character.toUpperCase()),
    })),
  ];

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setFieldErrors({});
    setFormError("");

    const errors: FormErrors = {};
    if (!name.trim()) errors.name = "is required";
    if (!accountTypeID) errors.account_type_id = "is required";
    if (isGoldAccount) {
      if (!goldPrice?.price_per_gram) {
        errors.gold_grams = "gold price is not available";
      }
      if (!goldGrams || Number(goldGrams) <= 0) {
        errors.gold_grams = "must be greater than 0";
      }
      if (!balance || Number(balance) <= 0) {
        errors.balance = "must be greater than 0";
      }
    }
    if (isStockAccount) {
      if (!/^[A-Z]{4}$/.test(name.trim().toUpperCase())) {
        errors.name = "must be a 4-letter IDX ticker, e.g. BBCA";
      }
      if (!stockLots || Number(stockLots) <= 0) {
        errors.stock_lots = "must be greater than 0";
      }
    }

    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors);
      return;
    }

    try {
      await createAccount(
        name,
        Number(accountTypeID),
        isStockAccount ? 0 : Number(balance) || 0,
        isGoldAccount ? Number(goldGrams) || 0 : undefined,
        isStockAccount ? name.trim().toUpperCase() : undefined,
        isStockAccount ? Number(stockLots) || 0 : undefined,
      );
      navigate("/accounts", { replace: true });
    } catch (error) {
      setFieldErrors(getValidationErrors(error));
      setFormError(getErrorMessage(error, "Unable to create account."));
    }
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <NeoPageHeader
        title="Add Account"
        description="Add a new bank, e-wallet, cash, gold, or brokerage account."
        eyebrow="New balance source"
        icon="➕"
      />
      <Card>
        {formError && (
          <NeoAlert className="mb-4" variant="danger">
            {formError}
          </NeoAlert>
        )}
        <form className="space-y-4" onSubmit={handleSubmit}>
          <label className="block">
            <span className="text-sm font-medium text-slate-700">Name</span>
            <NeoInput
              type="text"
              value={name}
              onChange={(e) =>
                setName(
                  isStockAccount
                    ? e.target.value.toUpperCase()
                    : e.target.value,
                )
              }
              placeholder={isStockAccount ? "BBCA" : "BCA Savings"}
              list={isStockAccount ? "idx-symbols" : undefined}
            />
            {isStockAccount ? (
              <datalist id="idx-symbols">
                {COMMON_IDX_SYMBOLS.map((symbol) => (
                  <option key={symbol} value={symbol} />
                ))}
              </datalist>
            ) : null}
            {fieldErrors.name && (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.name}
              </span>
            )}
          </label>
          <label className="block">
            <span className="text-sm font-medium text-slate-700">
              Account Type
            </span>
            <NeoSelect
              className="mt-1"
              value={accountTypeID}
              options={accountTypeOptions}
              onChange={(value) => {
                setAccountTypeID(value);
                setBalance("");
                setGoldGrams("");
                setStockLots("");
                setGoldInputMode(null);
              }}
              placeholder="Select account type"
              ariaLabel="Account type"
            />
            {fieldErrors.account_type_id && (
              <span className="mt-1 block text-sm text-red-600">
                {fieldErrors.account_type_id}
              </span>
            )}
          </label>
          {isGoldAccount ? (
            <div className="grid gap-4 sm:grid-cols-2">
              <label className="block">
                <span className="text-sm font-medium text-slate-700">
                  Nominal IDR
                </span>
                <NeoInput
                  type="number"
                  min="0"
                  value={balance}
                  disabled={goldInputMode === "grams"}
                  onChange={(e) => {
                    const value = e.target.value;
                    setBalance(value);
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
                {fieldErrors.balance && (
                  <span className="mt-1 block text-sm text-red-600">
                    {fieldErrors.balance}
                  </span>
                )}
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
                      setBalance(
                        String(
                          Math.round(Number(value) * goldPrice.price_per_gram),
                        ),
                      );
                    } else if (!value) {
                      setBalance("");
                    }
                  }}
                  placeholder="1"
                />
                {fieldErrors.gold_grams && (
                  <span className="mt-1 block text-sm text-red-600">
                    {fieldErrors.gold_grams}
                  </span>
                )}
              </label>
              <p className="text-xs font-semibold text-slate-500 sm:col-span-2">
                Current Antam price: {formatIDR(goldPrice?.price_per_gram)} /
                gr. Fill one field and Fintrack will calculate the other.
              </p>
            </div>
          ) : isStockAccount ? (
            <div className="space-y-3 rounded-2xl border-2 border-slate-950 bg-blue-50 p-4 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:bg-slate-900 dark:shadow-[3px_3px_0_0_#f8fafc]">
              <label className="block">
                <span className="text-sm font-medium text-slate-700 dark:text-slate-200">
                  Lot Saham
                </span>
                <NeoInput
                  type="number"
                  min="0"
                  step="0.01"
                  value={stockLots}
                  onChange={(e) => setStockLots(e.target.value)}
                  placeholder="10"
                />
                {fieldErrors.stock_lots && (
                  <span className="mt-1 block text-sm text-red-600">
                    {fieldErrors.stock_lots}
                  </span>
                )}
              </label>
              <p className="text-xs font-semibold text-slate-500">
                1 lot = 100 lembar. Fintrack akan validasi emiten BEI dan
                menghitung nilai rupiah otomatis dari harga pasar terbaru.
              </p>
            </div>
          ) : (
            <label className="block">
              <span className="text-sm font-medium text-slate-700">
                Initial Balance
              </span>
              <NeoInput
                type="number"
                value={balance}
                onChange={(e) => setBalance(e.target.value)}
                placeholder="5000000"
              />
              {fieldErrors.balance && (
                <span className="mt-1 block text-sm text-red-600">
                  {fieldErrors.balance}
                </span>
              )}
            </label>
          )}
          <div className="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
            <Link to="/accounts" className="sm:inline-flex">
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
              {isLoading ? "Creating..." : "Create Account"}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}
