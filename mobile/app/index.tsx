import React, { useEffect, useState } from 'react';
import { View } from 'react-native';
import { Text } from 'react-native-paper';
import { useRouter } from 'expo-router';
import { useAuthStore } from '../lib/store';
import Button from '../components/Button';
import { useHealth } from '../lib/api';

export default function Index() {
  const router = useRouter();
  const token = useAuthStore((s) => s.sessionToken);
  const { ok, isLoading } = useHealth();
  const [checked, setChecked] = useState(false);

  useEffect(() => {
    // after health check, route based on auth state
    if (!isLoading) {
      setChecked(true);
      if (token) router.replace('/(tabs)/home');
    }
  }, [isLoading, token, router]);

  if (token && checked) {
    return (
      <View style={{ padding: 16 }}>
        <Text>Redirecting…</Text>
      </View>
    );
  }

  return (
    <View style={{ padding: 16 }}>
      <Text variant="titleLarge">Welcome to Marko</Text>
      <Text style={{ marginTop: 8 }}>Connectivity: {isLoading ? 'Checking…' : ok ? 'OK' : 'Failed'}</Text>
      <Button onPress={() => router.replace('/(auth)/login')}>Continue</Button>
    </View>
  );
}