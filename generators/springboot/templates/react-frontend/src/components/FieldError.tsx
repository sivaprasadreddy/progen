interface FieldErrorProps {
  id: string;
  message: string | null;
}

export default function FieldError({ id, message }: FieldErrorProps) {
  if (!message) {
    return null;
  }
  return (
    <p id={id} className="mt-2 text-sm text-red-600" role="alert">
      {message}
    </p>
  );
}
