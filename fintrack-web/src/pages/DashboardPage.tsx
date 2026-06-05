import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "../components/common/Button";
import { Card } from "../components/common/Card";
import { NeoEmptyState } from "../components/common/NeoEmptyState";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { NeoStatCard } from "../components/common/NeoStatCard";
import { SkeletonCard } from "../components/common/Skeleton";
import { useAccountStore } from "../stores/accountStore";
import { useReportStore } from "../stores/reportStore";
import { useTransactionStore } from "../stores/transactionStore";
import {
  accountDisplayBalance,
  accountTypeLabel,
} from "../utils/accountDisplay";
import {
  formatDate,
  formatGoldGrams,
  formatIDR,
  transactionAmountLabel,
} from "../utils/format";
import { usePageTitle } from "../utils/usePageTitle";

export function DashboardPage() {
  usePageTitle("Dashboard");
  const {
    accounts,
    isLoading: loadingAccounts,
    fetchAccounts,
  } = useAccountStore();
  const {
    netWorth,
    totalSpending,
    totalIncome,
    isLoadingWorth,
    isLoadingSpending,
    fetchNetWorth,
    fetchSpending,
    spendingStartDate,
    spendingEndDate,
  } = useReportStore();
  const {
    transactions,
    isLoading: loadingTx,
    fetchRecent,
  } = useTransactionStore();

  const isLoading =
    loadingAccounts || isLoadingWorth || isLoadingSpending || loadingTx;

  useEffect(() => {
    fetchAccounts();
    fetchNetWorth();
    fetchRecent();
    fetchSpending(spendingStartDate, spendingEndDate);
  }, [
    fetchAccounts,
    fetchNetWorth,
    fetchRecent,
    fetchSpending,
    spendingStartDate,
    spendingEndDate,
  ]);

  return (
    <div className="space-y-6">
      <NeoPageHeader
        title="Dashboard"
        description="Overview of your net worth, accounts, and recent activity."
        eyebrow="Fintrack overview"
        icon="📈"
        actions={
          <>
            <Link to="/markets">
              <Button variant="secondary">Markets</Button>
            </Link>
            <Link to="/accounts">
              <Button variant="secondary">Manage Accounts</Button>
            </Link>
            <Link to="/transactions/new">
              <Button>Add Transaction</Button>
            </Link>
          </>
        }
      />

      {isLoading ? (
        <div className="space-y-6">
          <SkeletonCard />
          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <SkeletonCard />
            <SkeletonCard />
            <SkeletonCard />
          </div>
        </div>
      ) : null}

      {!isLoading && netWorth !== null ? (
        <>
          <NeoStatCard
            label="Net Worth"
            value={formatIDR(netWorth)}
            icon="💰"
            tone="blue"
          />

          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
            <Card>
              <p className="font-semibold text-slate-950 dark:text-slate-100">
                Account Balances
              </p>
              {accounts?.length === 0 ? (
                <div className="mt-4">
                  <p className="text-sm text-slate-500">No accounts yet.</p>
                  <Link to="/accounts">
                    <Button variant="secondary" className="mt-3">
                      Add Account
                    </Button>
                  </Link>
                </div>
              ) : (
                <ul className="mt-4 space-y-3">
                  {accounts.map((account) => (
                    <li
                      key={account.id}
                      className="flex items-start justify-between gap-3"
                    >
                      <div className="min-w-0">
                        <p className="truncate text-sm font-medium text-slate-950 dark:text-slate-100">
                          {account.name}
                        </p>
                        <p className="text-xs text-slate-500">
                          {accountTypeLabel(account.type)}
                          {account.type === "gold" && account.gold_grams != null
                            ? ` • ${formatGoldGrams(account.gold_grams)}`
                            : account.type === "stock_broker" &&
                                account.stock_lots != null
                              ? ` • ${account.stock_lots} lot @ ${formatIDR(account.stock_price_per_share)}`
                              : ""}
                        </p>
                      </div>
                      <p className="shrink-0 text-right text-sm font-semibold text-slate-950 dark:text-slate-100">
                        {formatIDR(accountDisplayBalance(account))}
                      </p>
                    </li>
                  ))}
                </ul>
              )}
            </Card>

            <Card>
              <p className="font-semibold text-slate-950 dark:text-slate-100">
                Summary This Month
              </p>

              {totalSpending === null && totalIncome === null ? (
                <p className="mt-4 text-sm text-slate-500">
                  No data this month.
                </p>
              ) : (
                <div className="mt-4 flex flex-col gap-4 sm:flex-row sm:items-center">
                  <DonutChart
                    income={totalIncome ?? 0}
                    expense={totalSpending ?? 0}
                  />
                  <div className="min-w-0 space-y-3">
                    <div>
                      <p className="text-xs font-bold text-green-600 uppercase">
                        Income
                      </p>
                      <p className="text-lg font-black text-green-700 dark:text-green-300">
                        {formatIDR(totalIncome)}
                      </p>
                    </div>
                    <div>
                      <p className="text-xs font-bold text-red-600 uppercase">
                        Spending
                      </p>
                      <p className="text-lg font-black text-red-700 dark:text-red-300">
                        {formatIDR(totalSpending)}
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </Card>

            <Card>
              <p className="font-semibold text-slate-950 dark:text-slate-100">
                Recent Transactions
              </p>
              {transactions.length === 0 ? (
                <div className="mt-4">
                  <p className="text-sm text-slate-500">
                    No recent transactions.
                  </p>
                  <Link to="/transactions/new">
                    <Button variant="secondary" className="mt-3">
                      Add Transaction
                    </Button>
                  </Link>
                </div>
              ) : (
                <ul className="mt-4 space-y-3">
                  {transactions.slice(0, 5).map((tx) => (
                    <li
                      key={tx.id}
                      className="flex items-center justify-between"
                    >
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-sm font-medium text-slate-950 dark:text-slate-100">
                          {tx.category?.name ?? tx.description ?? tx.type}
                        </p>
                        <p className="truncate text-xs text-slate-500">
                          {tx.account?.name} • {formatDate(tx.date)}
                        </p>
                      </div>
                      <p
                        className={`ml-2 whitespace-nowrap text-sm font-semibold ${tx.type === "expense" ? "text-red-600" : "text-green-600"}`}
                      >
                        {transactionAmountLabel(tx.type, tx.amount)}
                      </p>
                    </li>
                  ))}
                </ul>
              )}
              {transactions.length > 5 ? (
                <Link
                  to="/transactions"
                  className="neo-link mt-3 block text-sm"
                >
                  View all transactions
                </Link>
              ) : null}
            </Card>
          </div>
        </>
      ) : null}

      {!isLoading && netWorth === null ? (
        <Card>
          <NeoEmptyState
            title="No accounts yet"
            description="Add your first account to start tracking your net worth."
            icon="🏦"
            action={
              <Link to="/accounts">
                <Button>Add Account</Button>
              </Link>
            }
          />
        </Card>
      ) : null}
    </div>
  );
}

function DonutChart({ income, expense }: { income: number; expense: number }) {
  const [hoveredSegment, setHoveredSegment] = useState<
    "income" | "spending" | null
  >(null);
  const total = income + expense;
  const incomeRatio = total > 0 ? income / total : 0;
  const spendingRatio = total > 0 ? expense / total : 0;
  const radius = 36;
  const circumference = 2 * Math.PI * radius;
  const incomeDash = incomeRatio * circumference;
  const spendingDash = spendingRatio * circumference;
  const incomePercentage = Math.round(incomeRatio * 100);
  const spendingPercentage = Math.round(spendingRatio * 100);
  const hoveredLabel =
    hoveredSegment === "income"
      ? { label: "Income", percentage: incomePercentage }
      : hoveredSegment === "spending"
        ? { label: "Spending", percentage: spendingPercentage }
        : null;

  return (
    <div
      className="neo-surface relative h-32 w-32 shrink-0 rounded-full bg-yellow-100 p-3 dark:bg-slate-900"
      onMouseLeave={() => setHoveredSegment(null)}
    >
      <svg
        className="h-full w-full -rotate-90"
        viewBox="0 0 100 100"
        role="img"
        aria-label={`Income ${incomePercentage}%, spending ${spendingPercentage}%`}
      >
        <circle
          cx="50"
          cy="50"
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth="14"
          className="text-[#fffdf7] dark:text-slate-800"
        />
        <circle
          cx="50"
          cy="50"
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth="14"
          strokeLinecap="butt"
          strokeDasharray={`${incomeDash} ${circumference - incomeDash}`}
          className="cursor-pointer text-emerald-300 outline-none transition-opacity hover:opacity-80 focus:opacity-80"
          tabIndex={0}
          aria-label={`Income ${incomePercentage}%`}
          onFocus={() => setHoveredSegment("income")}
          onBlur={() => setHoveredSegment(null)}
          onMouseEnter={() => setHoveredSegment("income")}
        />
        <circle
          cx="50"
          cy="50"
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth="14"
          strokeLinecap="butt"
          strokeDasharray={`${spendingDash} ${circumference - spendingDash}`}
          strokeDashoffset={-incomeDash}
          className="cursor-pointer text-red-300 outline-none transition-opacity hover:opacity-80 focus:opacity-80"
          tabIndex={0}
          aria-label={`Spending ${spendingPercentage}%`}
          onFocus={() => setHoveredSegment("spending")}
          onBlur={() => setHoveredSegment(null)}
          onMouseEnter={() => setHoveredSegment("spending")}
        />
      </svg>
      {hoveredLabel ? (
        <div className="pointer-events-none absolute inset-0 flex items-center justify-center">
          <div className="rounded-xl border-2 border-slate-950 bg-[#fffdf7] px-2 py-1 text-center shadow-[2px_2px_0_0_#101828] dark:border-slate-100 dark:bg-slate-800 dark:shadow-[2px_2px_0_0_#f8fafc]">
            <p className="text-[0.6rem] font-black uppercase text-slate-500">
              {hoveredLabel.label}
            </p>
            <p className="text-sm font-black text-slate-950 dark:text-slate-100">
              {hoveredLabel.percentage}%
            </p>
          </div>
        </div>
      ) : null}
    </div>
  );
}
