import { useEffect, useState } from "react";
import { Card } from "../components/common/Card";
import { NeoDateInput } from "../components/common/NeoDateInput";
import { NeoEmptyState } from "../components/common/NeoEmptyState";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { NeoProgress } from "../components/common/NeoProgress";
import { NeoStatCard } from "../components/common/NeoStatCard";
import { SkeletonCard } from "../components/common/Skeleton";
import { useReportStore } from "../stores/reportStore";
import {
  clampFutureDateInput,
  getTodayDateInputValue,
} from "../utils/dateInput";
import { accountDisplayBalance } from "../utils/accountDisplay";
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

  const isLoading = isLoadingWorth || isLoadingSpending;

  useEffect(() => {
    fetchNetWorth();
    fetchSpending(spendingStartDate, spendingEndDate);
  }, [fetchNetWorth, fetchSpending, spendingStartDate, spendingEndDate]);

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

      {isLoading ? (
        <div className="space-y-6">
          <SkeletonCard />
          <SkeletonCard />
        </div>
      ) : (
        <>
          <NeoStatCard
            label="Net Worth"
            value={netWorth !== null ? formatIDR(netWorth) : "-"}
            icon="💰"
            tone="blue"
            helper={
              activeAccounts.length > 0 ? (
                <div className="space-y-2">
                  {activeAccounts.map((account) => (
                    <div
                      key={account.id}
                      className="flex flex-col gap-1 text-sm sm:flex-row sm:items-center sm:justify-between sm:gap-4"
                    >
                      <span className="min-w-0 truncate text-slate-600 dark:text-slate-200">
                        {account.name}
                      </span>
                      <span className="shrink-0 font-black text-slate-950 dark:text-slate-100">
                        {formatIDR(accountDisplayBalance(account))}
                      </span>
                    </div>
                  ))}
                </div>
              ) : null
            }
          />

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
                className="neo-button bg-blue-300 dark:text-slate-950"
                onClick={handleApply}
              >
                Apply
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
