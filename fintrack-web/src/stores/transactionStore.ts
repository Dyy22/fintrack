import { create } from "zustand";
import { api } from "../services/api";
import type { Transaction } from "../types";

type TransactionFilters = {
  start_date?: string;
  end_date?: string;
  account_id?: string;
  category_id?: string;
  type?: string;
};

type TransactionState = {
  transactions: Transaction[];
  isLoading: boolean;
  limit: number;
  offset: number;
  filters: TransactionFilters;
  createTransaction: (payload: Record<string, unknown>) => Promise<Transaction>;
  fetchTransactions: (
    filters?: TransactionFilters,
    limit?: number,
    offset?: number,
  ) => Promise<void>;
  setFilters: (filters: TransactionFilters) => void;
  nextPage: () => void;
  prevPage: () => void;
  fetchRecent: () => Promise<void>;
};

function todayString(): string {
  const now = new Date();
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`;
}

function thirtyDaysAgo(): string {
  const now = new Date();
  now.setDate(now.getDate() - 30);
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`;
}

export const useTransactionStore = create<TransactionState>((set, get) => ({
  transactions: [],
  isLoading: false,
  limit: 25,
  offset: 0,
  filters: {
    start_date: thirtyDaysAgo(),
    end_date: todayString(),
  },

  async createTransaction(payload: Record<string, unknown>) {
    const { data } = await api.post<Transaction>("/transactions", payload);
    return data;
  },

  async fetchTransactions(
    filters?: TransactionFilters,
    limit?: number,
    offset?: number,
  ) {
    set({ isLoading: true });
    try {
      const mergedFilters = filters ?? get().filters;
      const mergedLimit = limit ?? get().limit;
      const mergedOffset = offset ?? get().offset;
      const { data } = await api.get<{ transactions: Transaction[] }>(
        "/transactions",
        {
          params: {
            ...mergedFilters,
            limit: mergedLimit,
            offset: mergedOffset,
          },
        },
      );
      set({
        transactions: data.transactions ?? [],
        limit: mergedLimit,
        offset: mergedOffset,
        filters: mergedFilters,
      });
    } finally {
      set({ isLoading: false });
    }
  },

  setFilters(filters: TransactionFilters) {
    set({ filters, offset: 0 });
  },

  nextPage() {
    set((state) => ({ offset: state.offset + state.limit }));
  },

  prevPage() {
    set((state) => ({ offset: Math.max(0, state.offset - state.limit) }));
  },

  async fetchRecent() {
    set({ isLoading: true });
    try {
      const { data } = await api.get<{ transactions: Transaction[] }>(
        "/transactions",
        {
          params: { limit: 5, offset: 0 },
        },
      );
      set({ transactions: data.transactions ?? [] });
    } finally {
      set({ isLoading: false });
    }
  },
}));
