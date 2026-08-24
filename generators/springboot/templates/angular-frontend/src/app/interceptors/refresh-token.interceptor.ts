import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, switchMap, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';

const AUTH_ENDPOINTS_WITHOUT_REFRESH = ['/auth/login', '/auth/refresh'];

export const refreshTokenInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  if (AUTH_ENDPOINTS_WITHOUT_REFRESH.some((path) => req.url.includes(path))) {
    return next(req);
  }

  return next(req).pipe(
    catchError((error: unknown) => {
      if (!(error instanceof HttpErrorResponse) || error.status !== 401) {
        return throwError(() => error);
      }

      return authService.refreshAccessToken().pipe(
        switchMap(() => next(req)),
        catchError((refreshError) => {
          authService.logout();
          router.navigate(['/login'], { queryParams: { returnUrl: router.url } });
          return throwError(() => refreshError);
        }),
      );
    }),
  );
};
