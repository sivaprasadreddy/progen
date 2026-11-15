import { useState, type FormEvent } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router';
import { useAuth } from '../hooks/useAuth';
import { useForm } from '../hooks/useForm';
import { Validators } from '../utils/validators';
import { HttpError } from '../api/httpClient';
import { safeReturnUrl } from '../utils/navigation';
import Alert from '../components/Alert';
import IconInput from '../components/IconInput';
import PasswordInput from '../components/PasswordInput';
import FieldError from '../components/FieldError';

export default function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [isLoading, setIsLoading] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const form = useForm({
    initialValues: { email: '', password: '' },
    validators: {
      email: [Validators.required, Validators.email],
      password: [Validators.required, Validators.minLength(4)],
    },
  });

  const emailError = form.getFieldError('email');
  const passwordError = form.getFieldError('password');

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!form.isValid) {
      form.markAllAsTouched();
      return;
    }

    setIsLoading(true);
    setErrorMessage(null);
    try {
      await login(form.values.email, form.values.password);
      navigate(safeReturnUrl(searchParams.get('returnUrl')), { replace: true });
    } catch (error) {
      const status = error instanceof HttpError ? error.status : -1;
      if (status === 401) {
        setErrorMessage('Invalid email or password. Please try again.');
      } else if (status === 0) {
        setErrorMessage('Unable to connect to the server. Please check your connection.');
      } else {
        setErrorMessage('An error occurred during login. Please try again later.');
      }
      console.error('Login error:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-2xl shadow-2xl p-8 md:p-10">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">Welcome Back</h1>
            <p className="text-gray-600">Sign in to your account</p>
          </div>

          {errorMessage && (
            <div className="mb-6">
              <Alert type="error" message={errorMessage} />
            </div>
          )}

          <form onSubmit={handleSubmit} noValidate>
            <div className="mb-6">
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

            <div className="mb-6">
              <label htmlFor="password" className="block text-sm font-semibold text-gray-700 mb-2">
                Password
              </label>
              <PasswordInput
                id="password"
                placeholder="Enter your password"
                error={passwordError}
                {...form.register('password')}
              />
              <FieldError id="password-error" message={passwordError} />
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
                  Signing in...
                </span>
              ) : (
                <span>Sign In</span>
              )}
            </button>
          </form>

          <div className="mt-8 text-center">
            <p className="text-sm text-gray-600">
              Don't have an account?
              <Link
                to="/register"
                className="font-semibold text-indigo-600 hover:text-indigo-500 transition-colors no-underline ml-1"
              >
                Sign up now
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
