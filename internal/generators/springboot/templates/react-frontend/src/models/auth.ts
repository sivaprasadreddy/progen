export interface User {
  name: string;
  email: string;
  role: string;
}

export interface LoginResponse {
  accessToken: string;
  name: string;
  email: string;
  role: string;
}

export interface RefreshTokenResponse {
  accessToken: string;
  refreshToken: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  name: string;
  email: string;
  role: string;
}
