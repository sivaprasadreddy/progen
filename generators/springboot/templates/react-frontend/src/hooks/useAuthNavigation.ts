import { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { authStore } from '../services/authStore';

/**
 * Redirects to the login page (preserving the current URL as `returnUrl`)
 * when the session expires and the access token cannot be refreshed.
 */
export function useAuthNavigation(): void {
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    return authStore.onSessionExpired(() => {
      const returnUrl = location.pathname + location.search;
      navigate(`/login?returnUrl=${encodeURIComponent(returnUrl)}`, { replace: true });
    });
  }, [navigate, location]);
}
