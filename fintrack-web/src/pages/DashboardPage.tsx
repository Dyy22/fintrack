import { useEffect } from "react";
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
          <Link to="/transactions/new">
            <Button>Add Transaction</Button>
          </Link>
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
                <MonthlySummary
                  income={totalIncome ?? 0}
                  spending={totalSpending ?? 0}
                />
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

function MonthlySummary({
  income,
  spending,
}: {
  income: number;
  spending: number;
}) {
  const net = income - spending;
  const spendingRatio = income > 0 ? Math.round((spending / income) * 100) : 0;
  const progress = Math.min(spendingRatio, 100);
  const isOverBudget = income > 0 && spending > income;

  return (
    <div className="mt-4 space-y-4">
      <div className="rounded-2xl border-2 border-slate-950 bg-yellow-100 p-4 shadow-[4px_4px_0_0_#101828] dark:border-slate-100 dark:bg-slate-800 dark:shadow-[4px_4px_0_0_#f8fafc]">
        <p className="text-xs font-black uppercase tracking-[0.16em] text-slate-600 dark:text-slate-200">
          Net This Month
        </p>
        <p
          className={`mt-2 break-words text-2xl font-black sm:text-3xl ${
            net < 0
              ? "text-red-700 dark:text-red-300"
              : "text-slate-950 dark:text-slate-100"
          }`}
        >
          {formatIDR(net)}
        </p>
        <p className="mt-1 text-xs font-semibold text-slate-600 dark:text-slate-300">
          {net >= 0 ? "Surplus after spending" : "Spending exceeded income"}
        </p>
      </div>

      <div>
        <div className="flex items-center justify-between gap-3 text-xs font-black uppercase tracking-wide text-slate-600 dark:text-slate-300">
          <span>Spent</span>
          <span>
            {income > 0 ? `${spendingRatio}% of income` : "No income"}
          </span>
        </div>
        <div className="mt-2 h-4 overflow-hidden rounded-full border-2 border-slate-950 bg-[#fffdf7] dark:border-slate-100 dark:bg-slate-800">
          <div
            className={`h-full rounded-full ${
              isOverBudget ? "bg-red-400" : "bg-emerald-300"
            }`}
            style={{ width: `${progress}%` }}
            role="progressbar"
            aria-label="Spending compared to income"
            aria-valuemin={0}
            aria-valuemax={100}
            aria-valuenow={progress}
          />
        </div>
      </div>

      <ul className="w-full space-y-3">
        <li className="flex items-start justify-between gap-3">
          <p className="text-sm font-medium text-emerald-700 dark:text-emerald-400">
            Income
          </p>
          <p className="shrink-0 text-right text-sm font-semibold text-emerald-700 dark:text-emerald-300">
            {formatIDR(income)}
          </p>
        </li>
        <li className="flex items-start justify-between gap-3">
          <p className="text-sm font-medium text-red-700 dark:text-red-400">
            Spending
          </p>
          <p className="shrink-0 text-right text-sm font-semibold text-red-700 dark:text-red-300">
            {formatIDR(spending)}
          </p>
        </li>
      </ul>
    </div>
  );
}
