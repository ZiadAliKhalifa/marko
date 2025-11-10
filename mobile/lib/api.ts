import axios from 'axios';
import useSWR, { SWRConfiguration } from 'swr';
import { useAuthStore } from './store';

const API_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080';

export const axiosInstance = axios.create({
  baseURL: API_URL,
});

axiosInstance.interceptors.request.use(async (config) => {
  const token = useAuthStore.getState().sessionToken;
  if (token) {
    config.headers = {
      ...(config.headers || {}),
      Authorization: `Bearer ${token}`,
    };
  }
  return config;
});

export const swrFetcher = async (url: string) => {
  const res = await axiosInstance.get(url);
  return res.data;
};

export const swrConfig: SWRConfiguration = {
  fetcher: swrFetcher,
  revalidateOnFocus: true,
  shouldRetryOnError: false,
};

// API helpers
export type Group = { id: string; name: string; createdAt?: string };
export type Notification = { id: string; message: string; createdAt: string };
export type Member = { id: string; name?: string; joinedAt?: string };

export async function createGroup(name: string) {
  const res = await axiosInstance.post('/api/v1/groups', { name });
  return res.data as Group;
}

export async function joinGroup(groupId: string) {
  const res = await axiosInstance.post(`/api/v1/groups/${groupId}/join`);
  return res.data as { success: boolean };
}

export async function postLocationUpdate(payload: { countryCode: string; status: 'arrived' | 'left' }) {
  const res = await axiosInstance.post('/api/v1/locations', payload);
  return res.data;
}

// SWR hooks
export function useHealth() {
  const { data, error, isLoading } = useSWR('/healthz');
  return { ok: !!data, error, isLoading };
}

export function useGroups() {
  const { data, error, isLoading, mutate } = useSWR<Group[]>('/api/v1/groups');
  return { groups: data || [], error, isLoading, mutate };
}

export function useGroupMembers(groupId?: string) {
  const key = groupId ? `/api/v1/groups/${groupId}/members` : null;
  const { data, error, isLoading } = useSWR<Member[]>(key as string);
  return { members: data || [], error, isLoading };
}

export function useNotifications() {
  const { data, error, isLoading } = useSWR<Notification[]>('/api/v1/notifications');
  return { notifications: data || [], error, isLoading };
}