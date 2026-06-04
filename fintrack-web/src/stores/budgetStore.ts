import { create } from "zustand";
import { api } from "../services/api";
import type { Budget, Category } from "../types";

type BudgetState = {
  budgets: Budget[];
  categories: Category[];
  loading: boolean;
  error: string | null;
  fetchBudgets: (month: number, year: number) => Promise<void>;
  fetchCategories: () => Promise<void>;
  createBudget: (data: {
    category_id: string;
    month: number;
    year: number;
    amount: number;
  }) => Promise<Budget>;
  updateBudget: (id: string, amount: number) => Promise<Budget>;
  deleteBudget: (id: string) => Promise<void>;
};

export const useBudgetStore = create<BudgetState>((set, get) => ({
  budgets: [],
  categories: [],
  loading: false,
  error: null,

  async fetchBudgets(month: number, year: number) {
    set({ loading: true, error: null });
    try {
      const { data } = await api.get<{ budgets: Budget[] }>("/budgets", {
        params: { month, year },
      });
      set({ budgets: data.budgets ?? [] });
    } catch {
      set({ error: "Failed to fetch budgets" });
    } finally {
      set({ loading: false });
    }
  },

  async fetchCategories() {
    if (get().categories.length > 0) return;
    try {
      const { data } = await api.get<{ categories: Category[] }>(
        "/categories",
      );
      set({ categories: data.categories ?? [] });
    } catch {
      // silent - categories fetch is ancillary
    }
  },

  async createBudget(data) {
    const res = await api.post<Budget>("/budgets", data);
    const budget = res.data;
    set((state) => ({ budgets: [...state.budgets, budget] }));
    return budget;
  },

  async updateBudget(id: string, amount: number) {
    const res = await api.put<Budget>(`/budgets/${id}`, { amount });
    const updated = res.data;
    set((state) => ({
      budgets: state.budgets.map((b) => (b.id === id ? updated : b)),
    }));
    return updated;
  },

  async deleteBudget(id: string) {
    await api.delete(`/budgets/${id}`);
    set((state) => ({
      budgets: state.budgets.filter((b) => b.id !== id),
    }));
  },
}));
