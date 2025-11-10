import React from 'react';
import { Button as PaperButton } from 'react-native-paper';

type Props = {
  children: React.ReactNode;
  onPress?: () => void;
  mode?: 'text' | 'outlined' | 'contained';
  disabled?: boolean;
};

export default function Button({ children, onPress, mode = 'contained', disabled }: Props) {
  return (
    <PaperButton mode={mode} onPress={onPress} disabled={disabled}>
      {children}
    </PaperButton>
  );
}