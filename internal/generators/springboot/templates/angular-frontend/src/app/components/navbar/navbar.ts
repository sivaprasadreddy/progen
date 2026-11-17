import { Component, ChangeDetectionStrategy, signal, inject, computed, ElementRef } from '@angular/core';
import { RouterLink, RouterLinkActive, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { UserAvatar } from '../user-avatar/user-avatar';

@Component({
  selector: 'app-navbar',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [RouterLink, RouterLinkActive, UserAvatar],
  templateUrl: './navbar.html',
  host: {
    '(document:click)': 'onDocumentClick($event)'
  }
})
export class Navbar {
  protected readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly elementRef = inject(ElementRef);
  protected readonly isDropdownOpen = signal(false);

  protected readonly isAdmin = computed(() => {
    const user = this.authService.user();
    return user?.role === 'ADMIN';
  });

  protected onDocumentClick(event: MouseEvent): void {
    const clickedInside = this.elementRef.nativeElement.contains(event.target);
    if (!clickedInside && this.isDropdownOpen()) {
      this.closeDropdown();
    }
  }

  protected toggleDropdown(): void {
    this.isDropdownOpen.update(value => !value);
  }

  protected closeDropdown(): void {
    this.isDropdownOpen.set(false);
  }

  protected logout(): void {
    this.authService.logout();
    this.closeDropdown();
    this.router.navigate(['/']);
  }
}
