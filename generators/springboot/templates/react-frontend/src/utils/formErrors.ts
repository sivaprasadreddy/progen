import type { ValidationErrors } from './validators';

/**
 * Maps error keys to human-readable error messages.
 */
export type ErrorMessage = string | ((error: ValidationErrors) => string);

export type ErrorMessageMap = Record<string, ErrorMessage>;

/**
 * Default error messages for common validators.
 * Can be overridden with custom messages per field.
 */
const defaultErrorMessages: ErrorMessageMap = {
  required: () => 'This field is required.',
  email: () => 'Please enter a valid email address.',
  minlength: (error) => `Must be at least ${error['requiredLength']} characters.`,
  maxlength: (error) => `Must not exceed ${error['requiredLength']} characters.`,
  min: (error) => `Value must be at least ${error['min']}.`,
  max: (error) => `Value must not exceed ${error['max']}.`,
  pattern: () => 'Please enter a valid format.',
  passwordMismatch: () => 'Passwords do not match.',
};

/**
 * Capitalizes the first letter of a string and adds spaces before capitals.
 *
 * @example
 * capitalize('fullName') // returns 'Full Name'
 * capitalize('email') // returns 'Email'
 */
export function capitalize(str: string): string {
  const withSpaces = str.replace(/([A-Z])/g, ' $1').trim();
  return withSpaces.charAt(0).toUpperCase() + withSpaces.slice(1);
}

/**
 * Formats an error message by evaluating functions and optionally prefixing with field name.
 */
function formatErrorMessage(
  message: ErrorMessage,
  errorValue: unknown,
  fieldName?: string,
): string {
  const errorDetails = (
    typeof errorValue === 'object' && errorValue !== null ? errorValue : {}
  ) as ValidationErrors;
  const errorMessage = typeof message === 'function' ? message(errorDetails) : message;

  // If message doesn't already include field name and we have one, prefix it
  if (fieldName && !errorMessage.toLowerCase().includes(fieldName.toLowerCase())) {
    const capitalizedFieldName = capitalize(fieldName);
    if (errorMessage.startsWith('This field')) {
      return errorMessage.replace('This field', capitalizedFieldName);
    }
    return `${capitalizedFieldName} ${errorMessage.charAt(0).toLowerCase()}${errorMessage.slice(1)}`;
  }

  return errorMessage;
}

/**
 * Gets the first error message for a field. Only returns errors for touched fields.
 *
 * @example
 * getFieldErrorMessage(errors.email, touched.email, 'email', { required: 'Email address is required' });
 */
export function getFieldErrorMessage(
  errors: ValidationErrors | null | undefined,
  touched: boolean,
  fieldName?: string,
  customMessages?: ErrorMessageMap,
): string | null {
  if (!touched || !errors) {
    return null;
  }
  const errorKey = Object.keys(errors)[0];
  if (!errorKey) {
    return null;
  }
  const errorValue = errors[errorKey];

  if (customMessages?.[errorKey]) {
    return formatErrorMessage(customMessages[errorKey], errorValue, fieldName);
  }
  if (defaultErrorMessages[errorKey]) {
    return formatErrorMessage(defaultErrorMessages[errorKey], errorValue, fieldName);
  }
  return 'Invalid value.';
}

/**
 * Gets the error message for a form-level validator (e.g. password match).
 * Only returns an error when at least one field is touched.
 */
export function getFormErrorMessage(
  formErrors: ValidationErrors | null,
  errorKey: string,
  anyFieldTouched: boolean,
  customMessage?: string,
): string | null {
  if (!formErrors || !formErrors[errorKey] || !anyFieldTouched) {
    return null;
  }
  if (customMessage) {
    return customMessage;
  }
  if (defaultErrorMessages[errorKey]) {
    return formatErrorMessage(defaultErrorMessages[errorKey], formErrors[errorKey]);
  }
  return 'Form validation failed.';
}
