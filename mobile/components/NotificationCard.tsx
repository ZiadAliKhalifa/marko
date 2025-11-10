import React from 'react';
import { Card, Text } from 'react-native-paper';
import { Notification } from '../lib/api';

type Props = {
  item: Notification;
};

export default function NotificationCard({ item }: Props) {
  return (
    <Card style={{ marginVertical: 8 }}>
      <Card.Content>
        <Text>{item.message}</Text>
        <Text variant="bodySmall" style={{ marginTop: 4 }}>
          {new Date(item.createdAt).toLocaleString()}
        </Text>
      </Card.Content>
    </Card>
  );
}