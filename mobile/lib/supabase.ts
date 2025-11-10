import { createClient } from '@supabase/supabase-js';
import { useAuthStore } from './store';

const SUPABASE_URL = process.env.EXPO_PUBLIC_SUPABASE_URL as string;
const SUPABASE_ANON_KEY = process.env.EXPO_PUBLIC_SUPABASE_ANON_KEY as string;

export const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY);

export async function signInWithEmail(email: string) {
  const { data, error } = await supabase.auth.signInWithOtp({ email, options: { shouldCreateUser: true } });
  if (error) throw error;
  return data;
}

export async function signUpWithEmail(email: string, password: string) {
  const { data, error } = await supabase.auth.signUp({ email, password });
  if (error) throw error;
  return data;
}

export function subscribeAuth() {
  supabase.auth.onAuthStateChange(async (_event, session) => {
    const token = session?.access_token || null;
    await useAuthStore.getState().setToken(token);
  });
}

export async function getSessionToken() {
  const { data } = await supabase.auth.getSession();
  return data.session?.access_token ?? null;
}