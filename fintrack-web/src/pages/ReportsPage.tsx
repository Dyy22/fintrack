import { useEffect, useRef, useState } from "react";
import { Card } from "../components/common/Card";
import { NeoDateInput } from "../components/common/NeoDateInput";
import { NeoEmptyState } from "../components/common/NeoEmptyState";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { NeoProgress } from "../components/common/NeoProgress";
import { SkeletonCard } from "../components/common/Skeleton";
import { useReportStore } from "../stores/reportStore";
import {
  clampFutureDateInput,
  getTodayDateInputValue,
} from "../utils/dateInput";
import {
  accountDisplayBalance,
  accountTypeLabel,
} from "../utils/accountDisplay";
import { formatIDR } from "../utils/format";

import { usePageTitle } from "../utils/usePageTitle";

export function ReportsPage() {
  usePageTitle("Reports");
  const {
    netWorth,
    activeAccounts,
    totalSpending,
    spendingCategories,
    isLoadingWorth,
    isLoadingSpending,
    spendingStartDate,
    spendingEndDate,
    fetchNetWorth,
    fetchSpending,
  } = useReportStore();

  const maxDate = getTodayDateInputValue();
  const [localStartDate, setLocalStartDate] = useState(() =>
    clampFutureDateInput(spendingStartDate, maxDate),
  );
  const [localEndDate, setLocalEndDate] = useState(() =>
    clampFutureDateInput(spendingEndDate, maxDate),
  );
  const initialSpendingStartDate = useRef(spendingStartDate);
  const initialSpendingEndDate = useRef(spendingEndDate);

  const isInitialLoading =
    (isLoadingWorth && netWorth === null) ||
    (isLoadingSpending && totalSpending === null);

  useEffect(() => {
    fetchNetWorth();
    fetchSpending(
      initialSpendingStartDate.current,
      initialSpendingEndDate.current,
    );
  }, [fetchNetWorth, fetchSpending]);

  function handleApply() {
    const startDate = clampFutureDateInput(localStartDate, maxDate);
    const endDate = clampFutureDateInput(localEndDate, maxDate);
    setLocalStartDate(startDate);
    setLocalEndDate(endDate);
    fetchSpending(startDate, endDate);
  }

  return (
    <div className="space-y-6">
      <NeoPageHeader
        title="Reports"
        description="Analyze net worth and spending by category."
        eyebrow="Financial insights"
        icon="📊"
      />

      {isInitialLoading ? (
        <div className="space-y-6">
          <SkeletonCard />
          <SkeletonCard />
        </div>
      ) : (
        <>
          <Card className="bg-blue-100">
            <div className="flex items-start justify-between gap-4">
              <div className="min-w-0">
                <p className="text-xs font-black uppercase tracking-[0.16em] text-slate-600 dark:text-slate-200">
                  Net Worth
                </p>
                <p className="mt-2 break-words text-2xl font-black text-slate-950 dark:text-slate-100 sm:text-3xl">
                  {netWorth !== null ? formatIDR(netWorth) : "-"}
                </p>
              </div>
              <div
                className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border-2 border-slate-950 bg-blue-300 text-lg font-black text-slate-950 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:shadow-[3px_3px_0_0_#f8fafc] sm:h-12 sm:w-12 sm:text-xl"
                aria-hidden="true"
              >
                💰
              </div>
            </div>

            <div className="mt-5 border-t-2 border-slate-950 pt-4 dark:border-slate-100">
              <div className="mb-3 flex items-center justify-between gap-3">
                <p className="text-sm font-black uppercase tracking-wide text-slate-700 dark:text-slate-200">
                  Active Accounts
                </p>
                <p className="shrink-0 text-xs font-bold text-slate-500 dark:text-slate-300">
                  {activeAccounts.length} account
                  {activeAccounts.length === 1 ? "" : "s"}
                </p>
              </div>

              {activeAccounts.length > 0 ? (
                <ul className="space-y-3">
                  {activeAccounts.map((account) => (
                    <li
                      key={account.id}
                      className="flex items-start justify-between gap-3"
                    >
                      <div className="min-w-0">
                        <p className="truncate text-sm font-medium text-slate-950 dark:text-slate-100">
                          {account.name}
                        </p>
                        <p className="text-xs text-slate-500 dark:text-slate-300">
                          {accountTypeLabel(account.type)}
                        </p>
                      </div>
                      <p className="shrink-0 text-right text-sm font-semibold text-slate-950 dark:text-slate-100">
                        {formatIDR(accountDisplayBalance(account))}
                      </p>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-sm text-slate-500 dark:text-slate-300">
                  No active accounts included in net worth.
                </p>
              )}
            </div>
          </Card>

          <Card>
            <p className="text-sm font-medium text-slate-500 mb-4">
              Spending by Category
            </p>
            <div className="grid gap-3 sm:grid-cols-[1fr_1fr_auto]">
              <NeoDateInput
                value={localStartDate}
                max={maxDate}
                onChange={(value) => setLocalStartDate(value)}
                ariaLabel="Spending report start date"
              />
              <NeoDateInput
                value={localEndDate}
                max={maxDate}
                onChange={(value) => setLocalEndDate(value)}
                ariaLabel="Spending report end date"
              />
              <button
                type="button"
                className="neo-button bg-blue-300 dark:text-slate-950 disabled:cursor-wait disabled:opacity-70"
                onClick={handleApply}
                disabled={isLoadingSpending}
              >
                {isLoadingSpending ? "Applying..." : "Apply"}
              </button>
            </div>

            {totalSpending !== null && totalSpending > 0 ? (
              <>
                <p className="mt-6 break-words text-2xl font-bold text-red-600 sm:text-3xl">
                  {formatIDR(totalSpending)}
                </p>
                <p className="mt-1 text-sm text-slate-500">total spending</p>
                <div className="mt-4 space-y-3">
                  {spendingCategories.map((category) => (
                    <div key={category.name}>
                      <div className="flex items-start justify-between gap-3 text-sm">
                        <span className="min-w-0 break-words font-medium text-slate-950 dark:text-slate-100">
                          {category.name}
                        </span>
                        <span className="shrink-0 text-right font-semibold text-slate-950 dark:text-slate-100">
                          {formatIDR(category.amount)}
                        </span>
                      </div>
                      <div className="mt-1 flex items-center gap-2">
                        <NeoProgress value={category.percentage} />
                        <span className="text-xs text-slate-500">
                          {category.percentage.toFixed(0)}%
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              </>
            ) : (
              <NeoEmptyState
                className="mt-6"
                title="No spending data"
                description="No spending data for this period. Add expense transactions to see category breakdown."
                icon="📊"
              />
            )}
          </Card>
        </>
      )}
    </div>
  );
}
