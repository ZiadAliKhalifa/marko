import React from 'react';
import { ScrollView, View } from 'react-native';
import { Text } from 'react-native-paper';
import { useNotifications } from '../../lib/api';
import NotificationCard from '../../components/NotificationCard';

export default function Activity() {
  const { notifications, isLoading } = useNotifications();
  return (
    <ScrollView contentContainerStyle={{ padding: 16 }}>
      <Text variant="titleLarge">Activity</Text>
      <View style={{ marginTop: 16 }}>
        {isLoading ? (
          <Text>Loading…</Text>
        ) : notifications.length ? (
          notifications.map((n) => <NotificationCard key={n.id} item={n} />)
        ) : (
          <Text>No recent notifications</Text>
        )}
      </View>
    </ScrollView>
  );
}