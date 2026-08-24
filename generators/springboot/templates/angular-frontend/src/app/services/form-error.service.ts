import { Injectable } from '@angular/core';
import { AbstractControl, FormGroup, ValidationErrors } from '@angular/forms';

/**
 * Maps error keys to human-readable error messages.
 */
export interface ErrorMessageMap {
  [errorKey: string]: string | ((error: ValidationErrors) => string);
}

/**
 * Configuration for field-specific error messages.
 */
export interface FieldErrorConfig {
  [fieldName: string]: ErrorMessageMap;
}

/**
 * Service for handling form validation error messages.
 * Provides centralized error message generation with support for custom messages.
 */
@Injectable({
  providedIn: 'root',
})
export class FormErrorService {
  /**
   * Default error messages for common validators.
   * Can be overridden with custom messages per field.
   */
  private readonly defaultErrorMessages: ErrorMessageMap = {
    required: () => 'This field is required.',
    email: () => 'Please enter a valid email address.',
    minlength: (error: ValidationErrors) =>
      `Must be at least ${error['requiredLength']} characters.`,
    maxlength: (error: ValidationErrors) =>
      `Must not exceed ${error['requiredLength']} characters.`,
    min: (error: ValidationErrors) => `Value must be at least ${error['min']}.`,
    max: (error: ValidationErrors) => `Value must not exceed ${error['max']}.`,
    pattern: () => 'Please enter a valid format.',
    passwordMismatch: () => 'Passwords do not match.',
  };

  /**
   * Gets the first error message for a form field.
   * Only returns errors for touched fields.
   *
   * @param field - The form control to check
   * @param fieldName - Optional field name to include in error messages
   * @param customMessages - Optional custom error messages for this field
   * @returns Error message string or null if no errors
   *
   * @example
   * const error = formErrorService.getFieldError(
   *   form.get('email'),
   *   'Email',
   *   { required: 'Email address is required' }
   * );
   */
  getFieldError(
    field: AbstractControl | null,
    fieldName?: string,
    customMessages?: ErrorMessageMap,
  ): string | null {
    if (!field || !field.touched || !field.errors) {
      return null;
    }

    const errorKey = Object.keys(field.errors)[0];
    const errorValue = field.errors[errorKey];

    // Check for custom message first
    if (customMessages && customMessages[errorKey]) {
      return this.formatErrorMessage(
        customMessages[errorKey],
        errorValue,
        fieldName,
      );
    }

    // Fall back to default message
    if (this.defaultErrorMessages[errorKey]) {
      return this.formatErrorMessage(
        this.defaultErrorMessages[errorKey],
        errorValue,
        fieldName,
      );
    }

    // Unknown error type
    return 'Invalid value.';
  }

  /**
   * Gets error message for form-level validators (e.g., password match).
   * Only returns errors for forms where at least one field is touched.
   *
   * @param form - The form group to check
   * @param errorKey - The form-level error key to check for
   * @param customMessage - Optional custom error message
   * @returns Error message string or null if no errors
   *
   * @example
   * const error = formErrorService.getFormError(
   *   form,
   *   'passwordMismatch',
   *   'The passwords you entered do not match'
   * );
   */
  getFormError(
    form: FormGroup,
    errorKey: string,
    customMessage?: string,
  ): string | null {
    if (!form.errors || !form.errors[errorKey]) {
      return null;
    }

    // Check if any field is touched before showing form-level error
    const hasAnyTouchedField = Object.keys(form.controls).some(
      (key) => form.get(key)?.touched,
    );

    if (!hasAnyTouchedField) {
      return null;
    }

    if (customMessage) {
      return customMessage;
    }

    const errorValue = form.errors[errorKey];
    if (this.defaultErrorMessages[errorKey]) {
      return this.formatErrorMessage(
        this.defaultErrorMessages[errorKey],
        errorValue,
      );
    }

    return 'Form validation failed.';
  }

  /**
   * Gets multiple field errors from a form.
   * Useful for displaying all field errors at once.
   *
   * @param form - The form group to check
   * @param fieldConfigs - Map of field names to their custom error messages
   * @returns Map of field names to error messages
   *
   * @example
   * const errors = formErrorService.getAllFieldErrors(form, {
   *   email: { required: 'Email is required' },
   *   password: { minlength: 'Password too short' }
   * });
   */
  getAllFieldErrors(
    form: FormGroup,
    fieldConfigs?: FieldErrorConfig,
  ): { [fieldName: string]: string } {
    const errors: { [fieldName: string]: string } = {};

    Object.keys(form.controls).forEach((fieldName) => {
      const field = form.get(fieldName);
      const customMessages = fieldConfigs ? fieldConfigs[fieldName] : undefined;
      const error = this.getFieldError(field, fieldName, customMessages);

      if (error) {
        errors[fieldName] = error;
      }
    });

    return errors;
  }

  /**
   * Checks if a specific field has an error.
   * Only considers touched fields.
   *
   * @param field - The form control to check
   * @returns True if field has errors and is touched
   */
  hasFieldError(field: AbstractControl | null): boolean {
    return !!(field && field.touched && field.errors);
  }

  /**
   * Capitalizes the first letter of a string and adds spaces before capitals.
   * Useful for converting field names to human-readable format.
   *
   * @param str - String to capitalize
   * @returns Capitalized string
   *
   * @example
   * capitalize('fullName') // returns 'Full Name'
   * capitalize('email') // returns 'Email'
   */
  capitalize(str: string): string {
    const withSpaces = str.replace(/([A-Z])/g, ' $1').trim();
    return withSpaces.charAt(0).toUpperCase() + withSpaces.slice(1);
  }

  /**
   * Formats an error message by evaluating functions and optionally prefixing with field name.
   */
  private formatErrorMessage(
    message: string | ((error: ValidationErrors) => string),
    errorValue: ValidationErrors,
    fieldName?: string,
  ): string {
    let errorMessage: string;

    if (typeof message === 'function') {
      errorMessage = message(errorValue);
    } else {
      errorMessage = message;
    }

    // If message doesn't already include field name and we have one, prefix it
    if (
      fieldName &&
      !errorMessage.toLowerCase().includes(fieldName.toLowerCase())
    ) {
      const capitalizedFieldName = this.capitalize(fieldName);
      // Check if message starts with "This field" or similar generic terms
      if (errorMessage.startsWith('This field')) {
        return errorMessage.replace('This field', capitalizedFieldName);
      }
      if (errorMessage.startsWith('Must') || errorMessage.startsWith('Value')) {
        return `${capitalizedFieldName} ${errorMessage.charAt(0).toLowerCase()}${errorMessage.slice(1)}`;
      }
      return `${capitalizedFieldName} ${errorMessage.charAt(0).toLowerCase()}${errorMessage.slice(1)}`;
    }

    return errorMessage;
  }
}
