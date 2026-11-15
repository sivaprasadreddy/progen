import type { InputHTMLAttributes } from 'react';
import { cn } from '../utils/classNames';

interface IconInputProps extends InputHTMLAttributes<HTMLInputElement> {
  id: string;
  icon: string;
  error?: string | null;
}

/**
 * Text input with a leading Font Awesome icon and error-aware styling.
 */
export default function IconInput({ id, icon, error, className, ...inputProps }: IconInputProps) {
  return (
    <div className="relative">
      <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <i className={`fas ${icon} text-gray-400`} aria-hidden="true"></i>
      </div>
      <input
        id={id}
        className={cn(
          'block w-full pl-10 pr-3 py-3 border-2 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all',
          error ? 'border-red-300' : 'border-gray-300',
          className,
        )}
        aria-required="true"
        aria-invalid={error ? 'true' : 'false'}
        aria-describedby={error ? `${id}-error` : undefined}
        {...inputProps}
      />
    </div>
  );
}
