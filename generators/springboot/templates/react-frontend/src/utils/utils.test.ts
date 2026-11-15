import { getPageNumbers } from './pagination';
import { getUserInitials } from './user';
import { capitalize, getFieldErrorMessage } from './formErrors';
import { runValidators, Validators } from './validators';
import { safeReturnUrl } from './navigation';

describe('getPageNumbers', () => {
  it.each([
    [1, 3, [1, 2, 3]],
    [3, 10, [1, 2, 3, 4, 5, -1, 10]],
    [9, 10, [1, -1, 6, 7, 8, 9, 10]],
    [5, 10, [1, -1, 4, 5, 6, -1, 10]],
  ])('page %i of %i', (current, total, expected) => {
    expect(getPageNumbers(current, total)).toEqual(expected);
  });
});

describe('getUserInitials', () => {
  it('derives initials from the full name', () => {
    expect(getUserInitials('John Doe')).toBe('JD');
    expect(getUserInitials('John')).toBe('J');
    expect(getUserInitials('', 'johndoe')).toBe('J');
    expect(getUserInitials('')).toBe('');
  });
});

describe('form errors', () => {
  it('capitalizes field names', () => {
    expect(capitalize('fullName')).toBe('Full Name');
  });

  it('returns null for untouched fields', () => {
    expect(getFieldErrorMessage({ required: true }, false, 'email')).toBeNull();
  });

  it('formats validator errors with the field name', () => {
    const errors = runValidators('ab', [Validators.required, Validators.minLength(4)]);
    expect(getFieldErrorMessage(errors, true, 'password')).toBe(
      'Password must be at least 4 characters.',
    );
    expect(getFieldErrorMessage({ required: true }, true, 'email')).toBe('Email is required.');
  });

  it('prefers custom messages', () => {
    expect(getFieldErrorMessage({ pattern: {} }, true, 'phone', { pattern: 'Bad phone' })).toBe(
      'Bad phone',
    );
  });
});

describe('safeReturnUrl', () => {
  it('only allows same-origin relative paths', () => {
    expect(safeReturnUrl('/dashboard/profile')).toBe('/dashboard/profile');
    expect(safeReturnUrl('//evil.example.com')).toBe('/');
    expect(safeReturnUrl('https://evil.example.com')).toBe('/');
    expect(safeReturnUrl(null)).toBe('/');
  });
});
