export type User = {
  id: string;
  email: string;
  created_at?: string;
  updated_at?: string;
};

export type AccountType = {
  id: number;
  name: string;
  description: string;
};

export type Account = {
  id: string;
  account_type_id?: number;
  type: string;
  name: string;
  balance: number;
  currency: string;
  gold_grams?: number;
  gold_price_per_gram?: number;
  stock_symbol?: string;
  stock_lots?: number;
  stock_price_per_share?: number;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
};

export type Category = {
  id: string;
  name: string;
  type: "income" | "expense";
  is_default: boolean;
  created_at?: string;
  updated_at?: string;
};

export type TransactionType = "income" | "expense" | "transfer";

export type Transaction = {
  id: string;
  account_id?: string;
  category_id?: string;
  transfer_account_id?: string;
  type: TransactionType;
  amount: number;
  gold_grams?: number;
  description?: string;
  date: string;
  created_at?: string;
  updated_at?: string;
  account?: Account;
  category?: Category | null;
};

export type SpendingCategory = {
  name: string;
  amount: number;
  percentage: number;
};

export type GoldPrice = {
  price_per_gram: number;
  source: string;
  fetched_at: string;
  updated_at?: string;
};

export type GoldPriceHistoryPoint = {
  date: string;
  price_per_gram: number;
  source: string;
};

export type MarketChartPoint = {
  time: string;
  close: number;
};

export type MarketCacheStatus = "live" | "cached" | "stale";

export type MarketChart = {
  symbol: string;
  name?: string;
  currency: string;
  source: string;
  fetched_at: string;
  cache_status?: MarketCacheStatus;
  points: MarketChartPoint[];
};

export type ValidationErrorResponse = {
  error: "validation_error";
  message: string;
  fields: Record<string, string>;
};

export type ApiErrorResponse = {
  error: string;
  message: string;
};

export type Budget = {
  id: string;
  category_id: string;
  category: { id: string; name: string; type: string; is_default: boolean };
  month: number;
  year: number;
  amount: number;
  spent: number;
  remaining: number;
  percent: number;
  created_at: string;
  updated_at: string;
};
