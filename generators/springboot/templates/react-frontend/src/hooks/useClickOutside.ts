import { useEffect, type RefObject } from 'react';

/**
 * Invokes `handler` when a mouse click happens outside the referenced element.
 */
export function useClickOutside<T extends HTMLElement>(
  ref: RefObject<T | null>,
  handler: () => void,
  enabled = true,
): void {
  useEffect(() => {
    if (!enabled) {
      return;
    }
    const listener = (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        handler();
      }
    };
    document.addEventListener('click', listener);
    return () => document.removeEventListener('click', listener);
  }, [ref, handler, enabled]);
}
