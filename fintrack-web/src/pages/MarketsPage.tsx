import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "../components/common/Button";
import { Card } from "../components/common/Card";
import { NeoEmptyState } from "../components/common/NeoEmptyState";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { SkeletonCard } from "../components/common/Skeleton";
import { useAccountStore } from "../stores/accountStore";
import { useReportStore } from "../stores/reportStore";
import type { GoldPriceHistoryPoint, MarketChart } from "../types";
import { formatDate, formatIDR } from "../utils/format";
import { usePageTitle } from "../utils/usePageTitle";

const MARKET_RANGES = [
  { label: "5D", value: "5d" },
  { label: "1M", value: "1mo" },
  { label: "3M", value: "3mo" },
  { label: "6M", value: "6mo" },
  { label: "1Y", value: "1y" },
] as const;

type MarketRange = (typeof MARKET_RANGES)[number]["value"];

export function MarketsPage() {
  usePageTitle("Markets");
  const [selectedRange, setSelectedRange] = useState<MarketRange>("3mo");
  const {
    accounts,
    fetchAccounts,
    isLoading: loadingAccounts,
  } = useAccountStore();
  const {
    goldPrice,
    goldPriceHistory,
    fetchGoldPrice,
    fetchGoldPriceHistory,
    fetchMarketChart,
    marketCharts,
    isLoadingGoldPrice,
    isLoadingGoldPriceHistory,
    isLoadingMarketChart,
  } = useReportStore();

  const stockAccounts = useMemo(
    () =>
      accounts.filter(
        (account) => account.type === "stock_broker" && account.stock_symbol,
      ),
    [accounts],
  );
  const stockSymbols = useMemo(
    () =>
      Array.from(
        new Set(stockAccounts.map((account) => account.stock_symbol as string)),
      ),
    [stockAccounts],
  );
  const stockSymbolsKey = stockSymbols.join(",");
  const isLoading =
    loadingAccounts ||
    isLoadingGoldPrice ||
    isLoadingGoldPriceHistory ||
    isLoadingMarketChart;

  useEffect(() => {
    fetchAccounts();
    fetchGoldPrice().catch(() => undefined);
    fetchGoldPriceHistory(30).catch(() => undefined);
  }, [fetchAccounts, fetchGoldPrice, fetchGoldPriceHistory]);

  useEffect(() => {
    fetchMarketChart("IHSG", selectedRange).catch(() => undefined);
  }, [fetchMarketChart, selectedRange]);

  useEffect(() => {
    if (!stockSymbolsKey) return;
    stockSymbols.forEach((symbol) => {
      fetchMarketChart(symbol, selectedRange).catch(() => undefined);
    });
  }, [fetchMarketChart, selectedRange, stockSymbols, stockSymbolsKey]);

  return (
    <div className="space-y-6">
      <NeoPageHeader
        title="Markets"
        description="Track Antam gold price, IHSG, and stock holdings from your accounts."
        eyebrow="Market watch"
        icon="📊"
        actions={
          <Link to="/accounts/new">
            <Button>Add Stock Account</Button>
          </Link>
        }
      />

      <Card className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p className="text-sm font-black uppercase text-slate-700 dark:text-slate-200">
            Chart timeframe
          </p>
          <p className="text-xs font-semibold text-slate-500">
            Market data is cached and will use the last saved chart if the
            provider is unavailable.
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          {MARKET_RANGES.map((range) => (
            <button
              key={range.value}
              type="button"
              onClick={() => setSelectedRange(range.value)}
              className={`rounded-xl border-2 px-3 py-1.5 text-xs font-black transition ${
                selectedRange === range.value
                  ? "border-slate-950 bg-blue-500 text-white shadow-[2px_2px_0_0_#101828] dark:border-slate-100 dark:shadow-[2px_2px_0_0_#f8fafc]"
                  : "border-slate-300 bg-white text-slate-600 hover:border-slate-950 dark:border-slate-700 dark:bg-slate-900 dark:text-slate-300 dark:hover:border-slate-100"
              }`}
              aria-pressed={selectedRange === range.value}
            >
              {range.label}
            </button>
          ))}
        </div>
      </Card>

      {isLoading ? (
        <div className="grid gap-4 xl:grid-cols-2">
          <SkeletonCard />
          <SkeletonCard />
        </div>
      ) : null}

      <div className="grid gap-4 xl:grid-cols-2">
        <Card className="flex h-full flex-col">
          <p className="font-semibold text-slate-950 dark:text-slate-100">
            Antam Gold Price
          </p>
          {goldPrice ? (
            <>
              <p className="mt-2 text-3xl font-black text-yellow-700 dark:text-yellow-200">
                {formatIDR(goldPrice.price_per_gram)} / gr
              </p>
              <p className="mt-1 text-xs font-semibold text-slate-500">
                Updated {formatDate(goldPrice.fetched_at)} • {goldPrice.source}
              </p>
              <GoldPriceChart history={goldPriceHistory} />
            </>
          ) : (
            <p className="mt-4 text-sm text-slate-500">
              Gold price source is not configured yet.
            </p>
          )}
        </Card>

        <Card className="flex h-full flex-col">
          {marketCharts.IHSG?.points.length ? (
            <>
              <p className="font-semibold text-slate-950 dark:text-slate-100">
                IHSG
              </p>
              <p className="mt-2 text-3xl font-black text-blue-700 dark:text-blue-200">
                {formatIDR(
                  marketCharts.IHSG.points[marketCharts.IHSG.points.length - 1]
                    .close,
                )}
              </p>
              <p className="mt-1 text-xs font-semibold text-slate-500">
                {marketCharts.IHSG.name || "Jakarta Composite Index"} •{" "}
                {marketCharts.IHSG.source}
                {marketCharts.IHSG.fetched_at
                  ? ` • Updated ${formatDate(marketCharts.IHSG.fetched_at)}`
                  : ""}
              </p>
              <MarketLineChart
                chart={marketCharts.IHSG}
                fallbackLabel="IHSG"
                showHeader={false}
              />
            </>
          ) : (
            <>
              <p className="font-semibold text-slate-950 dark:text-slate-100">
                IHSG
              </p>
              <MarketLineChart
                chart={marketCharts.IHSG}
                fallbackLabel="IHSG"
                showHeader={false}
              />
            </>
          )}
        </Card>
      </div>

      <Card>
        <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="font-semibold text-slate-950 dark:text-slate-100">
              Portfolio Stocks
            </p>
            <p className="text-sm text-slate-500">
              Charts for every stock broker account you have added.
            </p>
          </div>
          <Link to="/accounts/new">
            <Button variant="secondary">Add Stock</Button>
          </Link>
        </div>

        {stockSymbols.length > 0 ? (
          <div className="mt-4 grid gap-4 lg:grid-cols-2">
            {stockSymbols.map((symbol) => {
              const symbolAccounts = stockAccounts.filter(
                (account) => account.stock_symbol === symbol,
              );
              const holdings = symbolAccounts.reduce(
                (sum, account) => sum + (account.stock_lots ?? 0),
                0,
              );
              const latestPrice = marketCharts[symbol]?.points.at(-1)?.close;
              const accountValue = symbolAccounts.reduce(
                (sum, account) => sum + account.balance,
                0,
              );
              const portfolioValue = latestPrice
                ? holdings * 100 * latestPrice
                : accountValue;
              return (
                <MarketLineChart
                  key={symbol}
                  chart={marketCharts[symbol]}
                  fallbackLabel={symbol}
                  subtitle={`${holdings} lot owned`}
                  portfolioValue={portfolioValue}
                />
              );
            })}
          </div>
        ) : (
          <NeoEmptyState
            className="mt-4"
            title="No stock accounts yet"
            description="Add a stock broker account to track IDX stock charts here."
            icon="📈"
            action={
              <Link to="/accounts/new">
                <Button>Add Stock Account</Button>
              </Link>
            }
          />
        )}
      </Card>
    </div>
  );
}

function MarketLineChart({
  chart,
  fallbackLabel,
  subtitle,
  portfolioValue,
  showHeader = true,
}: {
  chart?: MarketChart;
  fallbackLabel: string;
  subtitle?: string;
  portfolioValue?: number;
  showHeader?: boolean;
}) {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);
  const points = chart?.points ?? [];

  if (points.length === 0) {
    return (
      <div className="mt-4 rounded-2xl border-2 border-dashed border-slate-300 bg-slate-50 p-4 text-sm font-semibold text-slate-500 dark:border-slate-700 dark:bg-slate-900/60">
        <div className="mb-3 h-3 w-28 animate-pulse rounded-full bg-slate-200 dark:bg-slate-700" />
        Fetching {fallbackLabel} chart or waiting for cached data...
      </div>
    );
  }

  const width = 320;
  const height = 120;
  const paddingX = 14;
  const paddingY = 14;
  const prices = points.map((point) => point.close);
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
        : paddingY + ((maxPrice - point.close) / priceRange) * chartHeight;
    return { ...point, x, y };
  });
  const path = coordinates.map((point) => `${point.x},${point.y}`).join(" ");
  const hoveredPoint = hoveredIndex === null ? null : coordinates[hoveredIndex];
  const latest = points[points.length - 1];
  const first = points[0];
  const delta = latest.close - first.close;
  const deltaPercent = first.close > 0 ? (delta / first.close) * 100 : 0;
  const label = chart?.symbol ?? fallbackLabel;

  return (
    <div
      className={`${showHeader ? "" : "mt-4 flex flex-1 flex-col justify-end"} rounded-2xl border-2 border-slate-950 bg-blue-50 p-4 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:bg-slate-900 dark:shadow-[3px_3px_0_0_#f8fafc]`}
    >
      {showHeader ? (
        <div className="mb-3 flex items-start justify-between gap-3">
          <div className="min-w-0">
            <p className="text-sm font-black uppercase text-slate-700 dark:text-slate-200">
              {label}
            </p>
            <p className="truncate text-xs font-semibold text-slate-500">
              {subtitle || chart?.name || chart?.source || "Market chart"}
            </p>
            <p className="mt-1 text-[0.65rem] font-semibold text-slate-400">
              {chart?.source || "Market data"}
              {chart?.fetched_at
                ? ` • Updated ${formatDate(chart.fetched_at)}`
                : ""}
            </p>
            {portfolioValue !== undefined ? (
              <p className="mt-1 text-xs font-black text-blue-700 dark:text-blue-300">
                Portfolio value: {formatIDR(portfolioValue)}
              </p>
            ) : null}
          </div>
          <div className="shrink-0 text-right">
            <p className="text-lg font-black text-slate-950 dark:text-slate-100">
              {formatIDR(latest.close)}
            </p>
            <p
              className={`text-xs font-black ${delta >= 0 ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}
            >
              {delta >= 0 ? "+" : ""}
              {formatIDR(delta)} ({deltaPercent >= 0 ? "+" : ""}
              {deltaPercent.toFixed(2)}%)
            </p>
          </div>
        </div>
      ) : null}
      <div className="relative">
        <svg
          className="h-36 w-full overflow-visible"
          viewBox={`0 0 ${width} ${height}`}
          role="img"
          aria-label={`${label} market chart`}
          preserveAspectRatio="none"
        >
          <polyline
            points={path}
            fill="none"
            stroke="currentColor"
            strokeWidth="4"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="text-blue-500"
          />
          {coordinates.map((point, index) => (
            <circle
              key={`${point.time}-${index}`}
              cx={point.x}
              cy={point.y}
              r={hoveredIndex === index ? 5 : 3.5}
              className="cursor-pointer fill-[#fffdf7] stroke-slate-950 outline-none dark:fill-slate-800 dark:stroke-slate-100"
              strokeWidth="3"
              tabIndex={0}
              aria-label={`${formatDate(point.time)}: ${formatIDR(point.close)}`}
              onFocus={() => setHoveredIndex(index)}
              onBlur={() => setHoveredIndex(null)}
              onMouseEnter={() => setHoveredIndex(index)}
            />
          ))}
        </svg>
        {hoveredPoint ? (
          <div className="pointer-events-none absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 rounded-xl border-2 border-slate-950 bg-[#fffdf7] px-2 py-1 text-center shadow-[2px_2px_0_0_#101828] dark:border-slate-100 dark:bg-slate-800 dark:shadow-[2px_2px_0_0_#f8fafc]">
            <p className="text-[0.6rem] font-black uppercase text-slate-500">
              {formatDate(hoveredPoint.time)}
            </p>
            <p className="text-xs font-black text-slate-950 dark:text-slate-100">
              {formatIDR(hoveredPoint.close)}
            </p>
          </div>
        ) : null}
      </div>
    </div>
  );
}

function GoldPriceChart({ history }: { history: GoldPriceHistoryPoint[] }) {
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);
  const points = history.slice(-30);

  if (points.length === 0) {
    return (
      <div className="mt-4 rounded-2xl border-2 border-dashed border-slate-300 p-4 text-sm font-semibold text-slate-500 dark:border-slate-700">
        Gold price history will appear after the next refresh.
      </div>
    );
  }

  const width = 320;
  const height = 120;
  const paddingX = 14;
  const paddingY = 14;
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
    <div className="mt-4 rounded-2xl border-2 border-slate-950 bg-yellow-50 p-4 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:bg-slate-900 dark:shadow-[3px_3px_0_0_#f8fafc]">
      <div className="mb-3 flex items-center justify-between gap-2">
        <p className="text-xs font-black uppercase text-slate-600 dark:text-slate-300">
          30D Trend
        </p>
        <div className="text-right">
          <p className="text-sm font-black text-slate-950 dark:text-slate-100">
            {formatIDR(latest.price_per_gram)} / gr
          </p>
          <p
            className={`text-xs font-black ${priceDelta >= 0 ? "text-green-700 dark:text-green-300" : "text-red-700 dark:text-red-300"}`}
          >
            {priceDelta >= 0 ? "+" : ""}
            {formatIDR(priceDelta)}
          </p>
        </div>
      </div>
      <div className="relative">
        <svg
          className="h-36 w-full overflow-visible"
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
              r={hoveredIndex === index ? 5 : 3.5}
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
    </div>
  );
}
