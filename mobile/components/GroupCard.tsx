import React from 'react';
import { Card, Text } from 'react-native-paper';
import { Group } from '../lib/api';

type Props = {
  group: Group;
  onPress?: () => void;
};

export default function GroupCard({ group, onPress }: Props) {
  return (
    <Card onPress={onPress} style={{ marginVertical: 8 }}>
      <Card.Title title={group.name} subtitle={`ID: ${group.id}`} />
      {group.createdAt ? (
        <Card.Content>
          <Text variant="bodySmall">Created: {new Date(group.createdAt).toLocaleString()}</Text>
        </Card.Content>
      ) : null}
    </Card>
  );
}