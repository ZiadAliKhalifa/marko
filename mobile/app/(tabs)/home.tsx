import React, { useEffect, useState } from 'react';
import { View } from 'react-native';
import { Text } from 'react-native-paper';
import Button from '../../components/Button';
import { registerForPushNotificationsAsync, sendPushTokenToBackend } from '../../lib/notifications';
import { useHealth } from '../../lib/api';

export default function Home() {
  const { ok, isLoading } = useHealth();
  const [pushToken, setPushToken] = useState<string | undefined>();

  useEffect(() => {
    registerForPushNotificationsAsync().then((token) => {
      setPushToken(token);
      if (token) sendPushTokenToBackend(token);
    });
  }, []);

  return (
    <View style={{ padding: 16, gap: 12 }}>
      <Text variant="titleLarge">Home</Text>
      <Text>Backend Health: {isLoading ? 'Checking…' : ok ? 'OK' : 'Failed'}</Text>
      <Text>Push Token: {pushToken ? pushToken.slice(0, 12) + '…' : 'not registered'}</Text>
      <Button mode="outlined">Explore</Button>
    </View>
  );
}