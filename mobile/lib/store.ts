import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';

type AuthState = {
  sessionToken: string | null;
  setToken: (token: string | null) => Promise<void>;
  hydrate: () => Promise<void>;
};

export const useAuthStore = create<AuthState>((set) => ({
  sessionToken: null,
  setToken: async (token) => {
    if (token) {
      await SecureStore.setItemAsync('session_token', token);
    } else {
      await SecureStore.deleteItemAsync('session_token');
    }
    set({ sessionToken: token });
  },
  hydrate: async () => {
    const token = await SecureStore.getItemAsync('session_token');
    set({ sessionToken: token });
  },
}));