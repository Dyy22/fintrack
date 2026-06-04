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

export type ValidationErrorResponse = {
  error: "validation_error";
  message: string;
  fields: Record<string, string>;
};

export type ApiErrorResponse = {
  error: string;
  message: string;
};
