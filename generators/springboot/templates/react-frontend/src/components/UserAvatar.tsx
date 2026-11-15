import { Link } from 'react-router';
import { getUserInitials } from '../utils/user';
import { cn } from '../utils/classNames';

/**
 * Size variants for the avatar.
 * - 'xs': w-6 h-6, 'sm': w-8 h-8 (navbar), 'md': w-10 h-10 (default),
 * - 'lg': w-16 h-16, 'xl': w-24 h-24, '2xl': w-32 h-32 (profile pages)
 */
export type AvatarSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl';

/**
 * Color variants for the avatar background.
 */
export type AvatarColorVariant = 'primary' | 'secondary' | 'white' | 'gradient';

interface UserAvatarProps {
  /** Full name of the user, used for generating initials. */
  fullName?: string | null;
  size?: AvatarSize;
  colorVariant?: AvatarColorVariant;
  /** Optional link target. When provided, the avatar becomes clickable. */
  to?: string | null;
  /** Whether to show a hover effect. Defaults to true when `to` is provided. */
  showHover?: boolean | null;
  /** Optional ARIA label. Defaults to a label derived from the full name. */
  ariaLabel?: string | null;
}

const sizeClasses: Record<AvatarSize, string> = {
  xs: 'w-6 h-6 text-xs',
  sm: 'w-8 h-8 text-sm',
  md: 'w-10 h-10 text-base',
  lg: 'w-16 h-16 text-xl',
  xl: 'w-24 h-24 text-2xl',
  '2xl': 'w-32 h-32 text-4xl',
};

const colorClasses: Record<AvatarColorVariant, string> = {
  primary: 'bg-blue-600 text-white',
  secondary: 'bg-indigo-600 text-white',
  white: 'bg-white text-indigo-600 shadow-lg',
  gradient: 'bg-gradient-to-br from-blue-500 to-blue-600 text-white',
};

const hoverClasses: Partial<Record<AvatarColorVariant, string>> = {
  primary: 'hover:bg-blue-700',
  secondary: 'hover:bg-indigo-700',
};

/**
 * Reusable user avatar component for displaying user initials.
 */
export default function UserAvatar({
  fullName = null,
  size = 'md',
  colorVariant = 'primary',
  to = null,
  showHover = null,
  ariaLabel = null,
}: UserAvatarProps) {
  const initials = getUserInitials(fullName || '', 'User');
  const shouldShowHover = showHover ?? to !== null;
  const label = ariaLabel ?? (to !== null ? `View ${fullName} profile` : `${fullName} avatar`);

  const className = cn(
    'shrink-0 rounded-full flex items-center justify-center font-semibold',
    sizeClasses[size],
    colorClasses[colorVariant],
    to !== null && 'transition-colors',
    to !== null && shouldShowHover && hoverClasses[colorVariant],
  );

  if (to !== null) {
    return (
      <Link to={to} className={className} aria-label={label}>
        {initials}
      </Link>
    );
  }

  return (
    <div className={className} aria-label={label} role="img">
      {initials}
    </div>
  );
}
