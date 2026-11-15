import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router';
import { useAuth } from '../hooks/useAuth';
import { useForm } from '../hooks/useForm';
import { Validators } from '../utils/validators';
import { HttpError } from '../api/httpClient';
import type { RegisterRequest } from '../models/auth';
import IconInput from '../components/IconInput';
import PasswordInput from '../components/PasswordInput';
import FieldError from '../components/FieldError';

type RegisterFormValues = {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
};

const passwordMatchValidator = (values: RegisterFormValues) =>
  values.password === values.confirmPassword ? null : { passwordMismatch: true };

export default function Register() {
  const { register, login } = useAuth();
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const form = useForm<RegisterFormValues>({
    initialValues: { name: '', email: '', password: '', confirmPassword: '' },
    validators: {
      name: [Validators.required, Validators.minLength(2)],
      email: [Validators.required, Validators.email],
      password: [Validators.required, Validators.minLength(4)],
      confirmPassword: [Validators.required],
    },
    formValidator: passwordMatchValidator,
  });

  const nameError = form.getFieldError('name');
  const emailError = form.getFieldError('email');
  const passwordError = form.getFieldError('password');
  const confirmPasswordError =
    form.getFieldError('confirmPassword') ?? form.getFormError('passwordMismatch');

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!form.isValid) {
      form.markAllAsTouched();
      return;
    }

    setIsLoading(true);
    setErrorMessage(null);

    const { name, email, password } = form.values;
    const registerData: RegisterRequest = { name, email, password };

    try {
      await register(registerData);
    } catch (error) {
      setIsLoading(false);
      const status = error instanceof HttpError ? error.status : -1;
      if (status === 400) {
        setErrorMessage('Invalid registration data. Please check your inputs.');
      } else if (status === 409) {
        setErrorMessage('Email already exists. Please try another.');
      } else if (status === 0) {
        setErrorMessage('Unable to connect to the server. Please check your connection.');
      } else {
        setErrorMessage('An error occurred during registration. Please try again later.');
      }
      console.error('Registration error:', error);
      return;
    }

    // Auto-login after registration; fall back to the login page if it fails.
    try {
      await login(email, password);
      navigate('/', { replace: true });
    } catch {
      navigate('/login');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-2xl shadow-2xl p-8 md:p-10">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Create Account</h1>
          </div>

          {errorMessage && (
            <div
              className="bg-red-50 border-2 border-red-200 rounded-lg p-4 mb-6"
              role="alert"
              aria-live="polite"
            >
              <div className="flex items-start">
                <i
                  className="fas fa-exclamation-circle text-red-500 mt-0.5 mr-3"
                  aria-hidden="true"
                ></i>
                <p className="text-sm text-red-800 m-0">{errorMessage}</p>
              </div>
            </div>
          )}

          <form onSubmit={handleSubmit} noValidate>
            <div className="mb-5">
              <label htmlFor="name" className="block text-sm font-semibold text-gray-700 mb-2">
                Full Name
              </label>
              <IconInput
                id="name"
                type="text"
                icon="fa-user"
                placeholder="John Doe"
                error={nameError}
                {...form.register('name')}
              />
              <FieldError id="name-error" message={nameError} />
            </div>

            <div className="mb-5">
              <label htmlFor="email" className="block text-sm font-semibold text-gray-700 mb-2">
                Email Address
              </label>
              <IconInput
                id="email"
                type="email"
                icon="fa-envelope"
                placeholder="you@example.com"
                error={emailError}
                {...form.register('email')}
              />
              <FieldError id="email-error" message={emailError} />
            </div>

            <div className="mb-5">
              <label htmlFor="password" className="block text-sm font-semibold text-gray-700 mb-2">
                Password
              </label>
              <PasswordInput
                id="password"
                placeholder="Create a password"
                error={passwordError}
                {...form.register('password')}
              />
              <FieldError id="password-error" message={passwordError} />
            </div>

            <div className="mb-6">
              <label
                htmlFor="confirmPassword"
                className="block text-sm font-semibold text-gray-700 mb-2"
              >
                Confirm Password
              </label>
              <PasswordInput
                id="confirmPassword"
                placeholder="Confirm your password"
                error={confirmPasswordError}
                {...form.register('confirmPassword')}
              />
              <FieldError id="confirmPassword-error" message={confirmPasswordError} />
            </div>

            <button
              type="submit"
              disabled={isLoading}
              aria-busy={isLoading}
              className="w-full bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold py-3 px-4 rounded-lg transition-all duration-300 hover:from-indigo-700 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
            >
              {isLoading ? (
                <span className="flex items-center justify-center">
                  <i className="fas fa-spinner fa-spin mr-2" aria-hidden="true"></i>
                  Creating account...
                </span>
              ) : (
                <span>Create Account</span>
              )}
            </button>
          </form>

          <div className="mt-8 text-center">
            <p className="text-sm text-gray-600">
              Already have an account?
              <Link
                to="/login"
                className="font-semibold text-indigo-600 hover:text-indigo-500 transition-colors no-underline ml-1"
              >
                Sign in
              </Link>
            </p>
          </div>
        </div>

        <div className="mt-6 text-center">
          <Link
            to="/"
            className="text-white hover:text-gray-200 transition-colors no-underline flex items-center justify-center gap-2"
          >
            <i className="fas fa-arrow-left" aria-hidden="true"></i>
            Back to Home
          </Link>
        </div>
      </div>
    </div>
  );
}
