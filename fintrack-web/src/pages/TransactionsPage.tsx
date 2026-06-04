import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "../components/common/Button";
import { Card } from "../components/common/Card";
import { NeoBadge } from "../components/common/NeoBadge";
import { NeoDateInput } from "../components/common/NeoDateInput";
import { NeoEmptyState } from "../components/common/NeoEmptyState";
import { NeoSelect } from "../components/common/NeoSelect";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { NeoTable } from "../components/common/NeoTable";
import { Skeleton } from "../components/common/Skeleton";
import { useAccountStore } from "../stores/accountStore";
import { useTransactionStore } from "../stores/transactionStore";
import {
  clampFutureDateInput,
  getTodayDateInputValue,
} from "../utils/dateInput";
import { formatDate, transactionAmountLabel } from "../utils/format";

import { usePageTitle } from "../utils/usePageTitle";

export function TransactionsPage() {
  usePageTitle("Transactions");
  const { accounts, fetchAccounts } = useAccountStore();
  const {
    transactions,
    isLoading,
    filters,
    limit,
    offset,
    setFilters,
    fetchTransactions,
    nextPage,
    prevPage,
  } = useTransactionStore();
  const maxDate = getTodayDateInputValue();
  const [localStartDate, setLocalStartDate] = useState(() =>
    clampFutureDateInput(filters.start_date ?? "", maxDate),
  );
  const [localEndDate, setLocalEndDate] = useState(() =>
    clampFutureDateInput(filters.end_date ?? "", maxDate),
  );
  const [localType, setLocalType] = useState(filters.type ?? "");
  const [localAccountID, setLocalAccountID] = useState(
    filters.account_id ?? "",
  );

  useEffect(() => {
    fetchAccounts();
  }, [fetchAccounts]);

  useEffect(() => {
    fetchTransactions(undefined, limit, offset);
  }, [fetchTransactions, limit, offset]);

  type TransactionRow = (typeof transactions)[number];

  const accountOptions = [
    { value: "", label: "All accounts" },
    ...accounts
      .filter((account) => account.is_active)
      .map((account) => ({ value: account.id, label: account.name })),
  ];
  const typeOptions = [
    { value: "", label: "All types" },
    { value: "income", label: "Income" },
    { value: "expense", label: "Expense" },
    { value: "transfer", label: "Transfer" },
  ];

  const transactionColumns = [
    {
      key: "date",
      header: "Date",
      className: "text-slate-500",
      cell: (tx: TransactionRow) => formatDate(tx.date),
    },
    {
      key: "type",
      header: "Type",
      cell: (tx: TransactionRow) => (
        <NeoBadge
          variant={
            tx.type === "expense"
              ? "danger"
              : tx.type === "income"
                ? "success"
                : "info"
          }
        >
          {tx.type}
        </NeoBadge>
      ),
    },
    {
      key: "account",
      header: "Account",
      cell: (tx: TransactionRow) => tx.account?.name ?? "-",
    },
    {
      key: "category",
      header: "Category",
      cell: (tx: TransactionRow) =>
        tx.category?.name ?? (tx.type === "transfer" ? "-" : "Uncategorized"),
    },
    {
      key: "description",
      header: "Description",
      className: "text-slate-500",
      cell: (tx: TransactionRow) => tx.description || "-",
    },
    {
      key: "amount",
      header: "Amount",
      headerClassName: "text-right",
      className: "text-right font-semibold",
      cell: (tx: TransactionRow) => (
        <span
          className={tx.type === "expense" ? "text-red-600" : "text-green-600"}
        >
          {transactionAmountLabel(tx.type, tx.amount)}
        </span>
      ),
    },
  ];

  function handleApply() {
    const startDate = clampFutureDateInput(localStartDate, maxDate);
    const endDate = clampFutureDateInput(localEndDate, maxDate);
    setLocalStartDate(startDate);
    setLocalEndDate(endDate);

    setFilters({
      start_date: startDate,
      end_date: endDate,
      type: localType || undefined,
      account_id: localAccountID || undefined,
    });
    fetchTransactions(
      {
        start_date: startDate,
        end_date: endDate,
        type: localType || undefined,
        account_id: localAccountID || undefined,
      },
      limit,
      0,
    );
  }

  return (
    <div className="space-y-6">
      <NeoPageHeader
        title="Transactions"
        description="View, filter, and paginate your transaction history."
        eyebrow="Money movement"
        icon="🧾"
        actions={
          <Link to="/transactions/new">
            <Button>Add Transaction</Button>
          </Link>
        }
      />

      <Card>
        <div className="grid gap-3 md:grid-cols-5">
          <NeoDateInput
            value={localStartDate}
            max={maxDate}
            onChange={(value) => setLocalStartDate(value)}
            ariaLabel="Filter start date"
          />
          <NeoDateInput
            value={localEndDate}
            max={maxDate}
            onChange={(value) => setLocalEndDate(value)}
            ariaLabel="Filter end date"
          />
          <NeoSelect
            value={localAccountID}
            options={accountOptions}
            onChange={(value) => setLocalAccountID(value)}
            placeholder="All accounts"
            ariaLabel="Filter by account"
          />
          <NeoSelect
            value={localType}
            options={typeOptions}
            onChange={(value) => setLocalType(value)}
            placeholder="All types"
            ariaLabel="Filter by transaction type"
          />
          <Button type="button" onClick={handleApply}>
            Apply
          </Button>
        </div>
      </Card>

      {isLoading ? (
        <Card>
          <Skeleton className="h-4 w-20" />
          <div className="mt-4 space-y-2">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        </Card>
      ) : transactions.length === 0 ? (
        <Card>
          <NeoEmptyState
            title="No transactions found"
            description="Try changing filters or add a new transaction."
            icon="🧾"
            action={
              <Link to="/transactions/new">
                <Button>Add Transaction</Button>
              </Link>
            }
          />
        </Card>
      ) : (
        <>
          <Card>
            <NeoTable
              columns={transactionColumns}
              data={transactions}
              getRowKey={(tx) => tx.id}
            />

            <div className="space-y-3 md:hidden">
              {transactions.map((tx) => (
                <div key={tx.id} className="neo-surface rounded-xl p-4">
                  <div className="flex items-start justify-between gap-3">
                    <NeoBadge
                      variant={
                        tx.type === "expense"
                          ? "danger"
                          : tx.type === "income"
                            ? "success"
                            : "info"
                      }
                    >
                      {tx.type}
                    </NeoBadge>
                    <span
                      className={`shrink-0 text-right text-sm font-semibold ${tx.type === "expense" ? "text-red-600" : "text-green-600"}`}
                    >
                      {transactionAmountLabel(tx.type, tx.amount)}
                    </span>
                  </div>
                  <p className="mt-2 text-sm font-medium text-slate-950 dark:text-slate-100">
                    {tx.category?.name ?? tx.type}
                  </p>
                  <p className="mt-1 truncate text-xs text-slate-500">
                    {tx.account?.name ? `${tx.account.name} • ` : ""}
                    {formatDate(tx.date)}
                  </p>
                  {tx.description ? (
                    <p className="mt-1 break-words text-xs text-slate-500">
                      {tx.description}
                    </p>
                  ) : null}
                  {tx.type === "transfer" && tx.transfer_account_id ? (
                    <p className="mt-1 text-xs text-slate-500">
                      →{" "}
                      {accounts.find((a) => a.id === tx.transfer_account_id)
                        ?.name ?? "Destination"}
                    </p>
                  ) : null}
                </div>
              ))}
            </div>
          </Card>

          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <p className="text-sm text-slate-500">
              Showing {offset + 1}–{offset + transactions.length}
            </p>
            <div className="grid grid-cols-2 gap-2 sm:flex">
              <Button
                variant="secondary"
                onClick={prevPage}
                disabled={offset === 0}
              >
                Previous
              </Button>
              <Button
                variant="secondary"
                onClick={nextPage}
                disabled={transactions.length < limit}
              >
                Next
              </Button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
