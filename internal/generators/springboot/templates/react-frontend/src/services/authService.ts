import { http, refreshAccessToken } from '../api/httpClient';
import { authStore } from './authStore';
import type {
  LoginResponse,
  RefreshTokenResponse,
  RegisterRequest,
  RegisterResponse,
} from '../models/auth';

export const authService = {
  async login(email: string, password: string): Promise<LoginResponse> {
    const response = await http.post<LoginResponse>(
      '/api/auth/login',
      { email, password },
      { credentials: 'include' },
    );
    authStore.setAuthData(
      { name: response.name, email: response.email, role: response.role },
      response.accessToken,
    );
    return response;
  },

  register(data: RegisterRequest): Promise<RegisterResponse> {
    return http.post<RegisterResponse>('/api/users', data);
  },

  refreshAccessToken(): Promise<RefreshTokenResponse> {
    return refreshAccessToken();
  },

  logout(): void {
    authStore.clear();
  },
};
