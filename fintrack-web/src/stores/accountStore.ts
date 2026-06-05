import { create } from "zustand";
import { api } from "../services/api";
import type { Account, AccountType } from "../types";

type AccountState = {
  accountTypes: AccountType[];
  accounts: Account[];
  isLoading: boolean;
  fetchAccountTypes: () => Promise<void>;
  fetchAccounts: () => Promise<void>;
  createAccount: (
    name: string,
    accountTypeID: number,
    balance: number,
    goldGrams?: number,
    stockSymbol?: string,
    stockLots?: number,
  ) => Promise<Account>;
  updateAccount: (
    accountID: string,
    values: { name?: string; isActive?: boolean },
  ) => Promise<Account>;
  deactivateAccount: (accountID: string) => Promise<void>;
  hardDeleteAccount: (accountID: string) => Promise<void>;
};

export const useAccountStore = create<AccountState>((set, get) => ({
  accountTypes: [],
  accounts: [],
  isLoading: false,

  async fetchAccountTypes() {
    if (get().accountTypes.length > 0) return;
    const { data } = await api.get<{ account_types: AccountType[] }>(
      "/account-types",
    );
    set({ accountTypes: data.account_types });
  },

  async fetchAccounts() {
    set({ isLoading: true });
    try {
      const { data } = await api.get<{ accounts: Account[] }>("/accounts");
      set({ accounts: data.accounts ?? [] });
    } finally {
      set({ isLoading: false });
    }
  },

  async createAccount(
    name: string,
    accountTypeID: number,
    balance: number,
    goldGrams?: number,
    stockSymbol?: string,
    stockLots?: number,
  ) {
    const payload: Record<string, unknown> = {
      name,
      account_type_id: accountTypeID,
      balance,
    };
    if (goldGrams !== undefined) payload.gold_grams = goldGrams;
    if (stockSymbol !== undefined) payload.stock_symbol = stockSymbol;
    if (stockLots !== undefined) payload.stock_lots = stockLots;
    const { data } = await api.post<Account>("/accounts", payload);
    set((state) => ({ accounts: [data, ...state.accounts] }));
    return data;
  },

  async updateAccount(
    accountID: string,
    values: { name?: string; isActive?: boolean },
  ) {
    const payload: Record<string, unknown> = {};
    if (values.name !== undefined) payload.name = values.name;
    if (values.isActive !== undefined) payload.is_active = values.isActive;
    const { data } = await api.put<Account>(`/accounts/${accountID}`, payload);
    set((state) => ({
      accounts: state.accounts.map((account) =>
        account.id === accountID ? { ...account, ...data } : account,
      ),
    }));
    return data;
  },

  async deactivateAccount(accountID: string) {
    await api.delete(`/accounts/${accountID}`);
    set((state) => ({
      accounts: state.accounts.map((account) =>
        account.id === accountID ? { ...account, is_active: false } : account,
      ),
    }));
  },

  async hardDeleteAccount(accountID: string) {
    await api.delete(`/accounts/${accountID}?hard=true`);
    set((state) => ({
      accounts: state.accounts.filter((account) => account.id !== accountID),
    }));
  },
}));
