/**
 * Joins truthy class names into a single string.
 *
 * @example
 * cn('px-4', isActive && 'bg-blue-600', undefined) // 'px-4 bg-blue-600'
 */
export function cn(...classes: Array<string | false | null | undefined>): string {
  return classes.filter(Boolean).join(' ');
}
