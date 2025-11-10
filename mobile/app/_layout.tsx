import React, { useEffect } from 'react';
import { Stack } from 'expo-router';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { Provider as PaperProvider } from 'react-native-paper';
import { SWRConfig } from 'swr';
import { swrConfig } from '../lib/api';
import { useAuthStore } from '../lib/store';
import { subscribeAuth, getSessionToken } from '../lib/supabase';

export default function RootLayout() {
  const hydrate = useAuthStore((s) => s.hydrate);

  useEffect(() => {
    hydrate();
    subscribeAuth();
    getSessionToken().then((t) => useAuthStore.getState().setToken(t));
  }, [hydrate]);

  return (
    <SafeAreaProvider>
      <PaperProvider>
        <SWRConfig value={swrConfig}>
          <Stack screenOptions={{ headerShown: false }} initialRouteName="index">
            <Stack.Screen name="index" />
            <Stack.Screen name="(auth)" />
            <Stack.Screen name="(tabs)" />
            <Stack.Screen name="group/[id]" />
          </Stack>
        </SWRConfig>
      </PaperProvider>
    </SafeAreaProvider>
  );
}