import React from 'react';
import { useLocalSearchParams } from 'expo-router';
import { ScrollView, View } from 'react-native';
import { Text } from 'react-native-paper';
import { useGroupMembers } from '../../lib/api';

export default function GroupDetail() {
  const { id } = useLocalSearchParams();
  const groupId = Array.isArray(id) ? id[0] : id;
  const { members, isLoading } = useGroupMembers(groupId);

  return (
    <ScrollView contentContainerStyle={{ padding: 16 }}>
      <Text variant="titleLarge">Group {groupId}</Text>
      <View style={{ marginTop: 16, gap: 8 }}>
        {isLoading ? (
          <Text>Loading members…</Text>
        ) : members.length ? (
          members.map((m) => <Text key={m.id}>{m.name || m.id}</Text>)
        ) : (
          <Text>No members yet</Text>
        )}
      </View>
    </ScrollView>
  );
}