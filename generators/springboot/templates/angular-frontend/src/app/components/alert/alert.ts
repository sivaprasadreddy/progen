import { Component, input, output, ChangeDetectionStrategy } from '@angular/core';

export type AlertType = 'success' | 'error' | 'warning' | 'info';

@Component({
  selector: 'app-alert',
  templateUrl: './alert.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
  host: {
    'class': 'block'
  }
})
export class Alert {
  type = input.required<AlertType>();
  message = input.required<string>();
  dismissible = input<boolean>(false);

  dismissed = output<void>();

  protected getAlertClasses(): string {
    const baseClasses = 'border-2 rounded-lg p-4';
    const typeClasses = {
      success: 'bg-green-50 border-green-200',
      error: 'bg-red-50 border-red-200',
      warning: 'bg-yellow-50 border-yellow-200',
      info: 'bg-blue-50 border-blue-200'
    };
    return `${baseClasses} ${typeClasses[this.type()]}`;
  }

  protected getIconClass(): string {
    const iconClasses = {
      success: 'fa-check-circle text-green-500',
      error: 'fa-exclamation-circle text-red-500',
      warning: 'fa-exclamation-triangle text-yellow-500',
      info: 'fa-info-circle text-blue-500'
    };
    return `fas ${iconClasses[this.type()]}`;
  }

  protected getTextClass(): string {
    const textClasses = {
      success: 'text-green-800',
      error: 'text-red-800',
      warning: 'text-yellow-800',
      info: 'text-blue-800'
    };
    return `text-sm ${textClasses[this.type()]} m-0`;
  }

  protected onDismiss(): void {
    this.dismissed.emit();
  }
}
