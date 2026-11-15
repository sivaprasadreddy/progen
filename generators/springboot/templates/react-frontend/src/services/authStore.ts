import type { User } from '../models/auth';

/**
 * Framework-agnostic store holding the authenticated user and access token.
 * State is persisted to localStorage so the session survives page reloads.
 * React components subscribe to it via the `useAuth` hook; the HTTP client
 * reads the access token from it directly.
 */
export interface AuthState {
  user: User | null;
  accessToken: string | null;
}

const ACCESS_TOKEN_KEY = 'app_access_token';
const USER_KEY = 'app_auth_user';

type Listener = () => void;

const listeners = new Set<Listener>();
const sessionExpiredListeners = new Set<Listener>();

function hasStorage(): boolean {
  return typeof window !== 'undefined' && !!window.localStorage;
}

function loadState(): AuthState {
  if (!hasStorage()) {
    return { user: null, accessToken: null };
  }
  const userJson = localStorage.getItem(USER_KEY);
  let user: User | null = null;
  try {
    user = userJson ? (JSON.parse(userJson) as User) : null;
  } catch {
    localStorage.removeItem(USER_KEY);
  }
  return { user, accessToken: localStorage.getItem(ACCESS_TOKEN_KEY) };
}

let state: AuthState = loadState();

function setState(next: AuthState): void {
  state = next;
  listeners.forEach((listener) => listener());
}

export const authStore = {
  getState(): AuthState {
    return state;
  },

  subscribe(listener: Listener): () => void {
    listeners.add(listener);
    return () => listeners.delete(listener);
  },

  getAccessToken(): string | null {
    return state.accessToken;
  },

  isAuthenticated(): boolean {
    return state.user !== null && state.accessToken !== null;
  },

  setAuthData(user: User, accessToken: string): void {
    if (hasStorage()) {
      localStorage.setItem(USER_KEY, JSON.stringify(user));
      localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    }
    setState({ user, accessToken });
  },

  setAccessToken(accessToken: string): void {
    if (hasStorage()) {
      localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
    }
    setState({ ...state, accessToken });
  },

  clear(): void {
    if (hasStorage()) {
      localStorage.removeItem(ACCESS_TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
    }
    setState({ user: null, accessToken: null });
  },

  /** Registers a callback invoked when the session expires and the token cannot be refreshed. */
  onSessionExpired(listener: Listener): () => void {
    sessionExpiredListeners.add(listener);
    return () => sessionExpiredListeners.delete(listener);
  },

  notifySessionExpired(): void {
    sessionExpiredListeners.forEach((listener) => listener());
  },

  /** Re-reads state from localStorage. Intended for tests. */
  reload(): void {
    setState(loadState());
  },
};
