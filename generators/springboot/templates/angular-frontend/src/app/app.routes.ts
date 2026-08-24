import { Routes } from '@angular/router';
import { authGuard } from './guards/auth.guard';

export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./layout/default/layout').then((m) => m.Layout),
    children: [
      {
        path: '',
        loadComponent: () => import('./pages/home/home').then((m) => m.Home),
      },
      {
        path: 'login',
        loadComponent: () => import('./pages/login/login').then((m) => m.Login),
      },
      {
        path: 'register',
        loadComponent: () =>
          import('./pages/register/register').then((m) => m.Register),
      },
    ],
  },
  {
    path: 'dashboard',
    loadComponent: () =>
      import('./layout/dashboard-layout/dashboard-layout').then(
        (m) => m.DashboardLayout
      ),
    canActivate: [authGuard],
    children: [
      {
        path: '',
        loadComponent: () =>
          import('./pages/dashboard/dashboard-home/dashboard-home').then(
            (m) => m.DashboardHome
          ),
      },
      {
        path: 'profile',
        loadComponent: () =>
          import('./pages/dashboard/profile/profile').then((m) => m.Profile),
      },
    ],
  },
];
