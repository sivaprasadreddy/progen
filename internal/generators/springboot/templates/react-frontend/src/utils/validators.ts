/**
 * Validation errors keyed by validator name, e.g. `{ required: true }` or
 * `{ minlength: { requiredLength: 4, actualLength: 2 } }`.
 */
export type ValidationErrors = Record<string, unknown>;

export type Validator = (value: unknown) => ValidationErrors | null;

const EMAIL_REGEXP =
  /^(?=.{1,254}$)(?=.{1,64}@)[a-zA-Z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-zA-Z0-9!#$%&'*+/=?^_`{|}~-]+)*@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$/;

function isEmpty(value: unknown): boolean {
  return value === null || value === undefined || (typeof value === 'string' && value.length === 0);
}

/**
 * Common field validators, modelled after Angular's built-in `Validators`.
 */
export const Validators = {
  required: ((value) =>
    isEmpty(value) || (typeof value === 'string' && value.trim().length === 0)
      ? { required: true }
      : null) as Validator,

  email: ((value) =>
    isEmpty(value) || EMAIL_REGEXP.test(String(value)) ? null : { email: true }) as Validator,

  minLength(requiredLength: number): Validator {
    return (value) => {
      if (isEmpty(value)) return null;
      const actualLength = String(value).length;
      return actualLength < requiredLength ? { minlength: { requiredLength, actualLength } } : null;
    };
  },

  maxLength(requiredLength: number): Validator {
    return (value) => {
      if (isEmpty(value)) return null;
      const actualLength = String(value).length;
      return actualLength > requiredLength ? { maxlength: { requiredLength, actualLength } } : null;
    };
  },

  min(min: number): Validator {
    return (value) => {
      if (isEmpty(value)) return null;
      const actual = Number(value);
      return !Number.isNaN(actual) && actual < min ? { min: { min, actual } } : null;
    };
  },

  max(max: number): Validator {
    return (value) => {
      if (isEmpty(value)) return null;
      const actual = Number(value);
      return !Number.isNaN(actual) && actual > max ? { max: { max, actual } } : null;
    };
  },

  pattern(regexp: RegExp): Validator {
    return (value) => {
      if (isEmpty(value)) return null;
      return regexp.test(String(value))
        ? null
        : { pattern: { requiredPattern: regexp.toString(), actualValue: value } };
    };
  },
};

/**
 * Runs validators in order and merges their errors. Returns null when valid.
 */
export function runValidators(
  value: unknown,
  validators: Validator[] = [],
): ValidationErrors | null {
  const errors = validators.reduce<ValidationErrors>(
    (acc, validator) => ({ ...acc, ...(validator(value) ?? {}) }),
    {},
  );
  return Object.keys(errors).length > 0 ? errors : null;
}
