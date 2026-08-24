import { Injectable, signal, computed, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {finalize, Observable, shareReplay, tap} from 'rxjs';
import {
  LoginResponse,
  RefreshTokenResponse,
  RegisterRequest,
  RegisterResponse,
  User} from '../models/auth.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly ACCESS_TOKEN_KEY = 'app_access_token';
  private readonly USER_KEY = 'app_auth_user';

  private readonly currentUser = signal<User | null>(
    this.loadUserFromStorage());
  private readonly accessToken = signal<string | null>(
    this.loadTokenFromStorage(this.ACCESS_TOKEN_KEY),
  );
  private refreshTokenInProgress$: Observable<RefreshTokenResponse> | null = null;
  readonly isAuthenticated = computed(
    () => this.currentUser() !== null && this.accessToken() !== null
  );
  readonly user = this.currentUser.asReadonly();

  private loadTokenFromStorage(key: string): string | null {
    if (typeof window !== 'undefined' && window.localStorage) {
      return localStorage.getItem(key);
    }
    return null;
  }

  private loadUserFromStorage(): User | null {
    if (typeof window !== 'undefined' && window.localStorage) {
      const userJson = localStorage.getItem(this.USER_KEY);
      return userJson ? JSON.parse(userJson) : null;
    }
    return null;
  }

  login(email: string, password: string): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(`/api/auth/login`, { email, password }).pipe(
      tap((response) => {
        this.setAuthData(response);
      })
    );
  }

  refreshAccessToken(): Observable<RefreshTokenResponse> {
    if (!this.refreshTokenInProgress$) {
      this.refreshTokenInProgress$ = this.http
        .post<RefreshTokenResponse>(`api/auth/refresh`, {}, { withCredentials: true })
        .pipe(
          tap((response) => this.setAccessToken(response.accessToken)),
          finalize(() => (this.refreshTokenInProgress$ = null)),
          shareReplay(1),
        );
    }
    return this.refreshTokenInProgress$;
  }

  register(data: RegisterRequest): Observable<RegisterResponse> {
    return this.http.post<RegisterResponse>(`/api/users`, data);
  }

  private setAuthData(response: LoginResponse): void {
    const user: User = {
      name: response.name,
      email: response.email,
      role: response.role
    };

    this.currentUser.set(user);
    this.setAccessToken(response.accessToken);

    if (typeof window !== 'undefined' && window.localStorage) {
      localStorage.setItem(this.USER_KEY, JSON.stringify(user));
    }
  }

  private setAccessToken(accessToken: string): void {
    this.accessToken.set(accessToken);

    if (typeof window !== 'undefined' && window.localStorage) {
      localStorage.setItem(this.ACCESS_TOKEN_KEY, accessToken);
    }
  }

  logout(): void {
    this.accessToken.set(null);
    this.currentUser.set(null);

    if (typeof window !== 'undefined' && window.localStorage) {
      localStorage.removeItem(this.ACCESS_TOKEN_KEY);
      localStorage.removeItem(this.USER_KEY);
    }
  }

  getAccessToken(): string | null {
    return this.accessToken();
  }
}
