import { Component, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { UserAvatar } from '../../components/user-avatar/user-avatar';

interface NavItem {
  label: string;
  icon: string;
  path: string;
}

@Component({
  selector: 'app-dashboard-layout',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [RouterLink, RouterLinkActive, RouterOutlet, UserAvatar],
  templateUrl: './dashboard-layout.html',
  host: {
    class: 'flex flex-col h-screen overflow-hidden'
  }
})
export class DashboardLayout {
  protected readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  protected readonly navItems: NavItem[] = [
    { label: 'Dashboard', icon: 'fa-home', path: '/dashboard' },
    { label: 'Profile', icon: 'fa-user', path: '/dashboard/profile' }
  ];

  protected readonly isUserMenuOpen = signal(false);

  protected toggleUserMenu(): void {
    this.isUserMenuOpen.update(v => !v);
  }

  protected closeUserMenu(): void {
    this.isUserMenuOpen.set(false);
  }

  protected logout(): void {
    this.isUserMenuOpen.set(false);
    this.authService.logout();
    this.router.navigate(['/']);
  }
}
