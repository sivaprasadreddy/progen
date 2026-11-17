/**
 * Returns `returnUrl` if it is a safe, same-origin relative path; otherwise `fallback`.
 * Prevents open redirects via crafted `?returnUrl=` query parameters.
 */
export function safeReturnUrl(returnUrl: string | null | undefined, fallback = '/'): string {
  if (!returnUrl || !returnUrl.startsWith('/') || returnUrl.startsWith('//')) {
    return fallback;
  }
  return returnUrl;
}
