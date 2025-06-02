import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { AuthModel } from "@/api/types.gen";

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  user: AuthModel | null;
  setTokens: (accessToken: string, refreshToken: string) => void;
  setUser: (user: AuthModel | null) => void;
  clearTokens: () => void;
  clearUser: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      refreshToken: null,
      user: null,
      setTokens: (accessToken: string, refreshToken: string) =>
        set({ accessToken, refreshToken }),
      setUser: (user: AuthModel | null) => set({ user }),
      clearTokens: () => set({ accessToken: null, refreshToken: null, user: null }),
      clearUser: () => set({ user: null }),
    }),
    {
      name: 'auth-storage',
    }
  )
);
