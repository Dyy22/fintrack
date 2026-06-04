import { create } from 'zustand';
import { api, clearStoredToken, getStoredToken, setStoredToken } from '../services/api';
import type { User } from '../types';

type LoginResponse = {
  token: string;
  user: User;
};

type AuthState = {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => void;
  restoreSession: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  token: getStoredToken(),
  user: null,
  isAuthenticated: Boolean(getStoredToken()),
  isLoading: false,

  async login(email, password) {
    set({ isLoading: true });
    try {
      const { data } = await api.post<LoginResponse>('/auth/login', { email, password });
      setStoredToken(data.token);
      set({ token: data.token, user: data.user, isAuthenticated: true });
    } finally {
      set({ isLoading: false });
    }
  },

  async register(email, password) {
    set({ isLoading: true });
    try {
      await api.post<User>('/auth/register', { email, password });
    } finally {
      set({ isLoading: false });
    }
  },

  logout() {
    clearStoredToken();
    set({ token: null, user: null, isAuthenticated: false });
  },

  restoreSession() {
    const token = getStoredToken();
    set({ token, isAuthenticated: Boolean(token) });
  },
}));
