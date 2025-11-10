import React, { useState } from 'react';
import { View, ScrollView } from 'react-native';
import { Text, TextInput } from 'react-native-paper';
import Button from '../../components/Button';
import GroupCard from '../../components/GroupCard';
import { useGroups, createGroup, joinGroup } from '../../lib/api';
import { useRouter } from 'expo-router';

export default function Groups() {
  const { groups, isLoading, mutate } = useGroups();
  const [name, setName] = useState('');
  const [joinId, setJoinId] = useState('');
  const router = useRouter();

  const onCreate = async () => {
    if (!name.trim()) return;
    await createGroup(name.trim());
    setName('');
    mutate();
  };

  const onJoin = async () => {
    if (!joinId.trim()) return;
    await joinGroup(joinId.trim());
    setJoinId('');
    mutate();
  };

  return (
    <ScrollView contentContainerStyle={{ padding: 16 }}>
      <Text variant="titleLarge">Groups</Text>
      <View style={{ marginTop: 12, gap: 8 }}>
        <TextInput label="New group name" value={name} onChangeText={setName} />
        <Button onPress={onCreate}>Create</Button>
      </View>
      <View style={{ marginTop: 16, gap: 8 }}>
        <TextInput label="Join group by ID" value={joinId} onChangeText={setJoinId} />
        <Button onPress={onJoin}>Join</Button>
      </View>
      <View style={{ marginTop: 24 }}>
        {isLoading ? <Text>Loading…</Text> : groups.map((g) => (
          <GroupCard key={g.id} group={g} onPress={() => router.push(`/group/${g.id}`)} />
        ))}
      </View>
    </ScrollView>
  );
}