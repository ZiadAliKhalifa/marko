import React, { useState } from 'react';
import { View } from 'react-native';
import { TextInput, Text } from 'react-native-paper';
import Button from '../../components/Button';
import { supabase, signInWithEmail } from '../../lib/supabase';
import { useRouter } from 'expo-router';

export default function Login() {
  const [email, setEmail] = useState('');
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const onLogin = async () => {
    setLoading(true);
    setError(null);
    try {
      await signInWithEmail(email);
      // Magic link flow; alternatively allow password login
      router.replace('/(tabs)/home');
    } catch (e: any) {
      setError(e?.message || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <View style={{ padding: 16, gap: 12 }}>
      <Text variant="titleLarge">Login</Text>
      <TextInput label="Email" value={email} onChangeText={setEmail} autoCapitalize="none" keyboardType="email-address" />
      {error ? <Text style={{ color: 'red' }}>{error}</Text> : null}
      <Button onPress={onLogin} disabled={loading}>{loading ? 'Sending link…' : 'Send Magic Link'}</Button>
      <Button mode="text" onPress={() => router.push('/(auth)/signup')}>Create an account</Button>
    </View>
  );
}