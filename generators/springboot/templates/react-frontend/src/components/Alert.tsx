export type AlertType = 'success' | 'error' | 'warning' | 'info';

interface AlertProps {
  type: AlertType;
  message: string;
  dismissible?: boolean;
  onDismiss?: () => void;
}

const alertClasses: Record<AlertType, string> = {
  success: 'bg-green-50 border-green-200',
  error: 'bg-red-50 border-red-200',
  warning: 'bg-yellow-50 border-yellow-200',
  info: 'bg-blue-50 border-blue-200',
};

const iconClasses: Record<AlertType, string> = {
  success: 'fa-check-circle text-green-500',
  error: 'fa-exclamation-circle text-red-500',
  warning: 'fa-exclamation-triangle text-yellow-500',
  info: 'fa-info-circle text-blue-500',
};

const textClasses: Record<AlertType, string> = {
  success: 'text-green-800',
  error: 'text-red-800',
  warning: 'text-yellow-800',
  info: 'text-blue-800',
};

export default function Alert({ type, message, dismissible = false, onDismiss }: AlertProps) {
  return (
    <div className={`block border-2 rounded-lg p-4 ${alertClasses[type]}`} role="alert">
      <div className="flex items-center justify-between">
        <div className="flex items-center flex-1">
          <i className={`fas ${iconClasses[type]}`} aria-hidden="true"></i>
          <p className={`text-sm ${textClasses[type]} m-0 ml-3`}>{message}</p>
        </div>
        {dismissible && (
          <button
            type="button"
            onClick={onDismiss}
            className="ml-4 text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 rounded"
            aria-label="Dismiss alert"
          >
            <i className="fas fa-times" aria-hidden="true"></i>
          </button>
        )}
      </div>
    </div>
  );
}
