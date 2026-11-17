import { useCallback, useState, type ChangeEvent } from 'react';
import { runValidators, type ValidationErrors, type Validator } from '../utils/validators';
import {
  getFieldErrorMessage,
  getFormErrorMessage,
  type ErrorMessageMap,
} from '../utils/formErrors';

type FieldValues = Record<string, string>;

export interface UseFormOptions<T extends FieldValues> {
  initialValues: T;
  validators?: Partial<Record<keyof T, Validator[]>>;
  /** Form-level validator, e.g. to check that two fields match. */
  formValidator?: (values: T) => ValidationErrors | null;
}

/**
 * Lightweight form state management with field and form-level validation.
 * Errors are only surfaced for touched fields, mirroring Angular reactive forms.
 */
export function useForm<T extends FieldValues>({
  initialValues,
  validators = {},
  formValidator,
}: UseFormOptions<T>) {
  const [values, setValues] = useState<T>(initialValues);
  const [touched, setTouched] = useState<Partial<Record<keyof T, boolean>>>({});

  // Validation is cheap, so it is recomputed on every render.
  const fieldErrors: Partial<Record<keyof T, ValidationErrors | null>> = {};
  (Object.keys(values) as (keyof T)[]).forEach((name) => {
    fieldErrors[name] = runValidators(values[name], validators[name]);
  });
  const formErrors = formValidator?.(values) ?? null;

  const isValid =
    Object.values(fieldErrors).every((errors) => errors === null) && formErrors === null;

  const anyFieldTouched = Object.values(touched).some(Boolean);

  const setValue = useCallback((name: keyof T, value: string) => {
    setValues((prev) => ({ ...prev, [name]: value }));
  }, []);

  const handleChange = useCallback(
    (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
      const { name, value } = event.target;
      setValue(name as keyof T, value);
    },
    [setValue],
  );

  const handleBlur = useCallback((event: { target: { name: string } }) => {
    setTouched((prev) => ({ ...prev, [event.target.name]: true }));
  }, []);

  const markAllAsTouched = useCallback(() => {
    const all: Partial<Record<keyof T, boolean>> = {};
    (Object.keys(values) as (keyof T)[]).forEach((name) => {
      all[name] = true;
    });
    setTouched(all);
  }, [values]);

  const reset = useCallback((nextValues: T) => {
    setValues(nextValues);
    setTouched({});
  }, []);

  const getFieldError = (name: keyof T, customMessages?: ErrorMessageMap): string | null =>
    getFieldErrorMessage(fieldErrors[name], !!touched[name], String(name), customMessages);

  const getFormError = (errorKey: string, customMessage?: string): string | null =>
    getFormErrorMessage(formErrors, errorKey, anyFieldTouched, customMessage);

  /** Props to spread onto an input element bound to `name`. */
  const register = (name: keyof T) => ({
    name: String(name),
    value: values[name],
    onChange: handleChange,
    onBlur: handleBlur,
  });

  return {
    values,
    touched,
    isValid,
    setValue,
    handleChange,
    handleBlur,
    markAllAsTouched,
    reset,
    getFieldError,
    getFormError,
    register,
  };
}
