import { useState, type InputHTMLAttributes } from 'react';
import { cn } from '../utils/classNames';

interface PasswordInputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  id: string;
  error?: string | null;
}

/**
 * Password field with a lock icon and a show/hide toggle button.
 */
export default function PasswordInput({ id, error, className, ...inputProps }: PasswordInputProps) {
  const [showPassword, setShowPassword] = useState(false);

  return (
    <div className="relative">
      <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <i className="fas fa-lock text-gray-400" aria-hidden="true"></i>
      </div>
      <input
        id={id}
        type={showPassword ? 'text' : 'password'}
        className={cn(
          'block w-full pl-10 pr-12 py-3 border-2 rounded-lg text-gray-900 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all',
          error ? 'border-red-300' : 'border-gray-300',
          className,
        )}
        aria-required="true"
        aria-invalid={error ? 'true' : 'false'}
        aria-describedby={error ? `${id}-error` : undefined}
        {...inputProps}
      />
      <button
        type="button"
        onClick={() => setShowPassword((show) => !show)}
        className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600 focus:outline-none focus:text-gray-600"
        aria-label={showPassword ? 'Hide password' : 'Show password'}
      >
        <i className={showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'} aria-hidden="true"></i>
      </button>
    </div>
  );
}
