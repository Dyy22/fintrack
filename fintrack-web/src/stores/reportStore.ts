import { create } from "zustand";
import { api } from "../services/api";
import type {
  Account,
  GoldPrice,
  GoldPriceHistoryPoint,
  SpendingCategory,
} from "../types";

type NetWorthResponse = {
  net_worth: number;
  accounts: Account[];
};

type SpendingResponse = {
  start_date: string;
  end_date: string;
  total_spending: number;
  total_income: number;
  categories: SpendingCategory[];
};

type GoldPriceHistoryResponse = {
  history: GoldPriceHistoryPoint[];
};

function firstDayOfMonthString(): string {
  const now = new Date();
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-01`;
}

function todayString(): string {
  const now = new Date();
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`;
}

type ReportState = {
  netWorth: number | null;
  activeAccounts: Account[];
  totalSpending: number | null;
  totalIncome: number | null;
  spendingCategories: SpendingCategory[];
  goldPrice: GoldPrice | null;
  goldPriceHistory: GoldPriceHistoryPoint[];
  isLoadingWorth: boolean;
  isLoadingSpending: boolean;
  isLoadingGoldPrice: boolean;
  isLoadingGoldPriceHistory: boolean;
  spendingStartDate: string;
  spendingEndDate: string;
  fetchNetWorth: () => Promise<void>;
  fetchSpending: (startDate: string, endDate: string) => Promise<void>;
  fetchGoldPrice: () => Promise<void>;
  fetchGoldPriceHistory: (days?: number) => Promise<void>;
};

export const useReportStore = create<ReportState>((set) => ({
  netWorth: null,
  activeAccounts: [],
  totalSpending: null,
  totalIncome: null,
  spendingCategories: [],
  goldPrice: null,
  goldPriceHistory: [],
  isLoadingWorth: false,
  isLoadingSpending: false,
  isLoadingGoldPrice: false,
  isLoadingGoldPriceHistory: false,
  spendingStartDate: firstDayOfMonthString(),
  spendingEndDate: todayString(),

  async fetchNetWorth() {
    set({ isLoadingWorth: true });
    try {
      const { data } = await api.get<NetWorthResponse>("/reports/net-worth");
      set({ netWorth: data.net_worth, activeAccounts: data.accounts ?? [] });
    } finally {
      set({ isLoadingWorth: false });
    }
  },

  async fetchGoldPrice() {
    set({ isLoadingGoldPrice: true });
    try {
      const { data } = await api.get<GoldPrice>("/gold/price");
      set({ goldPrice: data });
    } finally {
      set({ isLoadingGoldPrice: false });
    }
  },

  async fetchGoldPriceHistory(days = 7) {
    set({ isLoadingGoldPriceHistory: true });
    try {
      const { data } = await api.get<GoldPriceHistoryResponse>(
        "/gold/prices/history",
        { params: { days } },
      );
      set({ goldPriceHistory: data.history ?? [] });
    } finally {
      set({ isLoadingGoldPriceHistory: false });
    }
  },

  async fetchSpending(startDate: string, endDate: string) {
    set({
      isLoadingSpending: true,
      spendingStartDate: startDate,
      spendingEndDate: endDate,
    });
    try {
      const { data } = await api.get<SpendingResponse>(
        "/reports/spending-by-category",
        {
          params: { start_date: startDate, end_date: endDate },
        },
      );
      set({
        totalSpending: data.total_spending,
        totalIncome: data.total_income,
        spendingCategories: data.categories ?? [],
      });
    } finally {
      set({ isLoadingSpending: false });
    }
  },
}));
