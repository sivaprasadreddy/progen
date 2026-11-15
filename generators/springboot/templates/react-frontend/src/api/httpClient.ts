import { authStore } from '../services/authStore';
import type { RefreshTokenResponse } from '../models/auth';

/**
 * Error thrown for non-2xx responses.
 * `status` is 0 when the server could not be reached.
 */
export class HttpError extends Error {
  readonly status: number;
  readonly body: unknown;

  constructor(status: number, message: string, body?: unknown) {
    super(message);
    this.name = 'HttpError';
    this.status = status;
    this.body = body;
  }
}

export interface RequestOptions extends Omit<RequestInit, 'body'> {
  body?: unknown;
}

const AUTH_ENDPOINTS_WITHOUT_REFRESH = ['/api/auth/login', '/api/auth/refresh'];

let refreshInProgress: Promise<RefreshTokenResponse> | null = null;

async function parseBody(response: Response): Promise<unknown> {
  const text = await response.text();
  if (!text) {
    return undefined;
  }
  const contentType = response.headers.get('content-type') ?? '';
  if (contentType.includes('json')) {
    try {
      return JSON.parse(text);
    } catch {
      return text;
    }
  }
  return text;
}

async function send(url: string, options: RequestOptions): Promise<Response> {
  const { body, headers, ...rest } = options;
  const requestHeaders = new Headers(headers);

  // Equivalent of an auth interceptor: attach the bearer token when available.
  const token = authStore.getAccessToken();
  if (token && !requestHeaders.has('Authorization')) {
    requestHeaders.set('Authorization', `Bearer ${token}`);
  }

  let requestBody: BodyInit | undefined;
  if (body !== undefined) {
    if (!requestHeaders.has('Content-Type')) {
      requestHeaders.set('Content-Type', 'application/json');
    }
    requestBody = typeof body === 'string' ? body : JSON.stringify(body);
  }
  if (!requestHeaders.has('Accept')) {
    requestHeaders.set('Accept', 'application/json');
  }

  try {
    return await fetch(url, { ...rest, headers: requestHeaders, body: requestBody });
  } catch (error) {
    throw new HttpError(0, 'Unable to connect to the server.', error);
  }
}

async function toResult<T>(response: Response): Promise<T> {
  const body = await parseBody(response);
  if (!response.ok) {
    throw new HttpError(response.status, response.statusText || 'Request failed', body);
  }
  return body as T;
}

/**
 * Requests a new access token using the HTTP-only refresh token cookie.
 * Concurrent callers share the same in-flight request.
 */
export function refreshAccessToken(): Promise<RefreshTokenResponse> {
  if (!refreshInProgress) {
    refreshInProgress = send('/api/auth/refresh', { method: 'POST', credentials: 'include' })
      .then((response) => toResult<RefreshTokenResponse>(response))
      .then((result) => {
        authStore.setAccessToken(result.accessToken);
        return result;
      })
      .finally(() => {
        refreshInProgress = null;
      });
  }
  return refreshInProgress;
}

/**
 * Performs an HTTP request against the backend API.
 * On a 401 response it transparently refreshes the access token and retries once.
 * If the refresh fails, the user is logged out and session-expired listeners are notified.
 */
export async function request<T>(url: string, options: RequestOptions = {}): Promise<T> {
  const response = await send(url, options);

  const skipRefresh = AUTH_ENDPOINTS_WITHOUT_REFRESH.some((path) => url.includes(path));
  if (response.status !== 401 || skipRefresh) {
    return toResult<T>(response);
  }

  try {
    await refreshAccessToken();
  } catch (refreshError) {
    authStore.clear();
    authStore.notifySessionExpired();
    throw refreshError;
  }
  return toResult<T>(await send(url, options));
}

export const http = {
  get: <T>(url: string, options?: RequestOptions) => request<T>(url, { ...options, method: 'GET' }),
  post: <T>(url: string, body?: unknown, options?: RequestOptions) =>
    request<T>(url, { ...options, method: 'POST', body }),
  put: <T>(url: string, body?: unknown, options?: RequestOptions) =>
    request<T>(url, { ...options, method: 'PUT', body }),
  patch: <T>(url: string, body?: unknown, options?: RequestOptions) =>
    request<T>(url, { ...options, method: 'PATCH', body }),
  delete: <T>(url: string, options?: RequestOptions) =>
    request<T>(url, { ...options, method: 'DELETE' }),
};
