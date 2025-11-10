import * as Location from 'expo-location';
import { postLocationUpdate } from './api';

export async function requestLocationPermissions() {
  const { status } = await Location.requestForegroundPermissionsAsync();
  return status === 'granted';
}

export async function updateCurrentCountry(status: 'arrived' | 'left') {
  const hasPerm = await requestLocationPermissions();
  if (!hasPerm) throw new Error('Location permission not granted');

  const pos = await Location.getCurrentPositionAsync({});
  const geocode = await Location.reverseGeocodeAsync({
    latitude: pos.coords.latitude,
    longitude: pos.coords.longitude,
  });
  const entry = geocode[0];
  const countryCode = (entry?.isoCountryCode || entry?.country || '').toUpperCase();
  if (!countryCode) throw new Error('Unable to determine country');

  return postLocationUpdate({ countryCode, status });
}