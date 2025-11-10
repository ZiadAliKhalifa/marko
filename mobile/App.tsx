import React, { useEffect } from 'react';
import { Slot } from 'expo-router';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { Provider as PaperProvider } from 'react-native-paper';
import { SWRConfig } from 'swr';
import { swrConfig } from './lib/api';
import { useAuthStore } from './lib/store';
import { subscribeAuth, getSessionToken } from './lib/supabase';

export default function App() {
  const hydrate = useAuthStore((s) => s.hydrate);

  useEffect(() => {
    hydrate();
    subscribeAuth();
    // Prime token from Supabase if available
    getSessionToken().then((t) => useAuthStore.getState().setToken(t));
  }, [hydrate]);

  return (
    <SafeAreaProvider>
      <PaperProvider>
        <SWRConfig value={swrConfig}>
          <Slot />
        </SWRConfig>
      </PaperProvider>
    </SafeAreaProvider>
  );
}