import { useSyncExternalStore } from 'react';
import { authStore } from '../services/authStore';
import { authService } from '../services/authService';

/**
 * Exposes the current authentication state and auth actions to components.
 */
export function useAuth() {
  const state = useSyncExternalStore(authStore.subscribe, authStore.getState, authStore.getState);

  return {
    user: state.user,
    isAuthenticated: state.user !== null && state.accessToken !== null,
    login: authService.login,
    register: authService.register,
    logout: authService.logout,
  };
}
