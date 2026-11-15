import { http, HttpError } from './httpClient';
import { authStore } from '../services/authStore';

function jsonResponse(status: number, body?: unknown): Response {
  return new Response(body === undefined ? null : JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

describe('httpClient', () => {
  const fetchMock = vi.fn<typeof fetch>();

  beforeEach(() => {
    fetchMock.mockReset();
    vi.stubGlobal('fetch', fetchMock);
    authStore.clear();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('attaches the bearer token when authenticated', async () => {
    authStore.setAuthData({ name: 'Jane', email: 'jane@example.com', role: 'ROLE_USER' }, 'abc');
    fetchMock.mockResolvedValueOnce(jsonResponse(200, { ok: true }));

    await http.get('/api/auth/me');

    const headers = fetchMock.mock.calls[0][1]!.headers as Headers;
    expect(headers.get('Authorization')).toBe('Bearer abc');
  });

  it('refreshes the access token on 401 and retries the request', async () => {
    authStore.setAuthData({ name: 'Jane', email: 'jane@example.com', role: 'ROLE_USER' }, 'old');
    fetchMock
      .mockResolvedValueOnce(jsonResponse(401))
      .mockResolvedValueOnce(jsonResponse(200, { accessToken: 'new', refreshToken: 'r' }))
      .mockResolvedValueOnce(jsonResponse(200, { name: 'Jane' }));

    const result = await http.get<{ name: string }>('/api/auth/me');

    expect(result.name).toBe('Jane');
    expect(authStore.getAccessToken()).toBe('new');
    const retryHeaders = fetchMock.mock.calls[2][1]!.headers as Headers;
    expect(retryHeaders.get('Authorization')).toBe('Bearer new');
  });

  it('logs out and notifies listeners when the refresh fails', async () => {
    authStore.setAuthData({ name: 'Jane', email: 'jane@example.com', role: 'ROLE_USER' }, 'old');
    const onExpired = vi.fn();
    const unsubscribe = authStore.onSessionExpired(onExpired);
    fetchMock.mockResolvedValueOnce(jsonResponse(401)).mockResolvedValueOnce(jsonResponse(401));

    await expect(http.get('/api/auth/me')).rejects.toBeInstanceOf(HttpError);

    expect(authStore.isAuthenticated()).toBe(false);
    expect(onExpired).toHaveBeenCalledOnce();
    unsubscribe();
  });

  it('maps network failures to status 0', async () => {
    fetchMock.mockRejectedValueOnce(new TypeError('Failed to fetch'));

    await expect(http.get('/api/anything')).rejects.toMatchObject({ status: 0 });
  });
});
