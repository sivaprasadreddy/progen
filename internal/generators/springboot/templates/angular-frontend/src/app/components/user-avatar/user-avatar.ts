import { Component, ChangeDetectionStrategy, input, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { getUserInitials } from '../../utils/user.util';

/**
 * Size variants for the avatar.
 * - 'xs': Extra small (w-6 h-6) - for compact displays
 * - 'sm': Small (w-8 h-8) - for navbar, inline lists
 * - 'md': Medium (w-10 h-10) - for message cards, default
 * - 'lg': Large (w-16 h-16) - for profile headers
 * - 'xl': Extra large (w-24 h-24) - for prominent displays
 * - '2xl': 2X large (w-32 h-32) - for full profile pages
 */
export type AvatarSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';

/**
 * Color variants for the avatar background.
 * - 'primary': Blue background (bg-blue-600, text-white)
 * - 'secondary': Indigo background (bg-indigo-600, text-white)
 * - 'white': White background (bg-white, colored text)
 * - 'gradient': Blue gradient background (from-blue-500 to-blue-600)
 */
export type AvatarColorVariant = 'primary' | 'secondary' | 'white' | 'gradient';

/**
 * Reusable user avatar component for displaying user initials.
 * Supports different sizes, colors, and optional clickable links.
 * Integrates with getUserInitials utility for consistent initials display.
 */
@Component({
  selector: 'app-user-avatar',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [CommonModule, RouterLink],
  templateUrl: './user-avatar.html',
})
export class UserAvatar {

  /**
   * Full name of the user (optional).
   * Used for generating initials.
   */
  readonly fullName = input<string | null>(null);

  /**
   * Size of the avatar.
   * Default: 'md'
   */
  readonly size = input<AvatarSize>('md');

  /**
   * Color variant of the avatar.
   * Default: 'primary'
   */
  readonly colorVariant = input<AvatarColorVariant>('primary');

  /**
   * Optional router link path for the avatar.
   * If provided, avatar becomes clickable and navigates to this path.
   * Typically used to link to user profile: ['/user', username]
   */
  readonly routerLink = input<string[] | null>(null);

  /**
   * Whether to show hover effect on the avatar.
   * Default: true when routerLink is provided, false otherwise
   */
  readonly showHover = input<boolean | null>(null);

  /**
   * Optional ARIA label for the avatar.
   * If not provided, generates default label from username.
   */
  readonly ariaLabel = input<string | null>(null);

  /**
   * Computed: User initials derived from fullName or username.
   */
  readonly initials = computed(() => {
    return getUserInitials(this.fullName() || '', 'User');
  });

  /**
   * Computed: Whether hover effect should be shown.
   * Defaults to true if routerLink is provided.
   */
  readonly shouldShowHover = computed(() => {
    const explicitHover = this.showHover();
    if (explicitHover !== null) {
      return explicitHover;
    }
    return this.routerLink() !== null;
  });

  /**
   * Computed: ARIA label for accessibility.
   */
  readonly computedAriaLabel = computed(() => {
    const explicit = this.ariaLabel();
    if (explicit !== null) {
      return explicit;
    }
    return this.routerLink() !== null
      ? `View ${this.fullName()} profile`
      : `${this.fullName()} avatar`;
  });
}
