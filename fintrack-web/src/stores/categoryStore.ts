import { create } from "zustand";
import { api } from "../services/api";
import type { Category } from "../types";

type CategoryState = {
  categories: Category[];
  isLoading: boolean;
  fetchCategories: (type?: string) => Promise<void>;
};

export const useCategoryStore = create<CategoryState>((set) => ({
  categories: [],
  isLoading: false,

  async fetchCategories(type?: string) {
    set({ isLoading: true });
    try {
      const params: Record<string, string> = {};
      if (type) params.type = type;
      const { data } = await api.get<{ categories: Category[] }>("/categories", { params });
      set({ categories: data.categories ?? [] });
    } finally {
      set({ isLoading: false });
    }
  },
}));
