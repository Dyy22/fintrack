import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "../components/common/Button";
import { Card } from "../components/common/Card";
import { NeoEmptyState } from "../components/common/NeoEmptyState";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { NeoStatCard } from "../components/common/NeoStatCard";
import { useAccountStore } from "../stores/accountStore";
import { useReportStore } from "../stores/reportStore";
import { useTransactionStore } from "../stores/transactionStore";
import { SkeletonCard } from "../components/common/Skeleton";
import { usePageTitle } from "../utils/usePageTitle";
import type { GoldPriceHistoryPoint } from "../types";
import {
  formatDate,
  formatGoldGrams,
  formatIDR,
  transactionAmountLabel,
} from "../utils/format";

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
    goldPrice,
    goldPriceHistory,
    isLoadingWorth,
    isLoadingSpending,
    fetchNetWorth,
    fetchSpending,
    fetchGoldPrice,
    fetchGoldPriceHistory,
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
    fetchGoldPrice().catch(() => undefined);
    fetchGoldPriceHistory(7).catch(() => undefined);
  }, [
    fetchAccounts,
    fetchNetWorth,
    fetchRecent,
    fetchSpending,
    fetchGoldPrice,
    fetchGoldPriceHistory,
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
                Antam Gold Price
              </p>
              {goldPrice ? (
                <>
                  <p className="mt-2 text-2xl font-black text-yellow-700 dark:text-yellow-200">
                    {formatIDR(goldPrice.price_per_gram)} / gr
                  </p>
                  <p className="mt-1 text-xs font-semibold text-slate-500">
                    Updated {formatDate(goldPrice.fetched_at)} •{" "}
                    {goldPrice.source}
                  </p>
                  <GoldPriceChart history={goldPriceHistory} />
                </>
              ) : (
                <p className="mt-4 text-sm text-slate-500">
                  Gold price source is not configured yet.
                </p>
              )}
            </Card>

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
                          {account.type}
                          {account.type === "gold" && account.gold_grams != null
                            ? ` • ${formatGoldGrams(account.gold_grams)}`
                            : ""}
                        </p>
                      </div>
                      <p className="shrink-0 text-right text-sm font-semibold text-slate-950 dark:text-slate-100">
                        {formatIDR(account.balance)}
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

function GoldPriceChart({ history }: { history: GoldPriceHistoryPoint[] }) {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);
  const points = history.slice(-7);

  if (points.length === 0) {
    return (
      <div className="mt-4 rounded-2xl border-2 border-dashed border-slate-300 p-3 text-xs font-semibold text-slate-500 dark:border-slate-700">
        Gold price history will appear after the next refresh.
      </div>
    );
  }

  const width = 240;
  const height = 84;
  const paddingX = 12;
  const paddingY = 12;
  const prices = points.map((point) => point.price_per_gram);
  const minPrice = Math.min(...prices);
  const maxPrice = Math.max(...prices);
  const priceRange = maxPrice - minPrice;
  const chartWidth = width - paddingX * 2;
  const chartHeight = height - paddingY * 2;
  const coordinates = points.map((point, index) => {
    const x =
      points.length === 1
        ? width / 2
        : paddingX + (index / (points.length - 1)) * chartWidth;
    const y =
      priceRange === 0
        ? height / 2
        : paddingY +
          ((maxPrice - point.price_per_gram) / priceRange) * chartHeight;
    return { ...point, x, y };
  });
  const path = coordinates.map((point) => `${point.x},${point.y}`).join(" ");
  const hoveredPoint = hoveredIndex === null ? null : coordinates[hoveredIndex];
  const latest = points[points.length - 1];
  const first = points[0];
  const priceDelta = latest.price_per_gram - first.price_per_gram;

  return (
    <div className="mt-4 rounded-2xl border-2 border-slate-950 bg-yellow-50 p-3 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:bg-slate-900 dark:shadow-[3px_3px_0_0_#f8fafc]">
      <div className="mb-2 flex items-center justify-between gap-2">
        <p className="text-xs font-black uppercase text-slate-600 dark:text-slate-300">
          7D Trend
        </p>
        <p
          className={`text-xs font-black ${priceDelta >= 0 ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}
        >
          {priceDelta >= 0 ? "+" : ""}
          {formatIDR(priceDelta)}
        </p>
      </div>
      <div className="relative">
        <svg
          className="h-24 w-full overflow-visible"
          viewBox={`0 0 ${width} ${height}`}
          role="img"
          aria-label={`Gold price chart for the last ${points.length} days`}
          preserveAspectRatio="none"
        >
          <polyline
            points={path}
            fill="none"
            stroke="currentColor"
            strokeWidth="4"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="text-yellow-500"
          />
          {coordinates.map((point, index) => (
            <circle
              key={point.date}
              cx={point.x}
              cy={point.y}
              r={hoveredIndex === index ? 5 : 4}
              className="cursor-pointer fill-[#fffdf7] stroke-slate-950 outline-none dark:fill-slate-800 dark:stroke-slate-100"
              strokeWidth="3"
              tabIndex={0}
              aria-label={`${point.date}: ${formatIDR(point.price_per_gram)} per gram`}
              onFocus={() => setHoveredIndex(index)}
              onBlur={() => setHoveredIndex(null)}
              onMouseEnter={() => setHoveredIndex(index)}
            />
          ))}
        </svg>
        {hoveredPoint ? (
          <div className="pointer-events-none absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 rounded-xl border-2 border-slate-950 bg-[#fffdf7] px-2 py-1 text-center shadow-[2px_2px_0_0_#101828] dark:border-slate-100 dark:bg-slate-800 dark:shadow-[2px_2px_0_0_#f8fafc]">
            <p className="text-[0.6rem] font-black uppercase text-slate-500">
              {hoveredPoint.date}
            </p>
            <p className="text-xs font-black text-slate-950 dark:text-slate-100">
              {formatIDR(hoveredPoint.price_per_gram)}
            </p>
          </div>
        ) : null}
      </div>
      {points.length < 2 ? (
        <p className="mt-1 text-xs font-semibold text-slate-500">
          Need more daily snapshots to draw a full weekly trend.
        </p>
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
