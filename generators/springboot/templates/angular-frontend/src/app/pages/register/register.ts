import { Component, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule, AbstractControl, ValidationErrors } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { FormErrorService } from '../../services/form-error.service';
import {RegisterRequest} from '../../models/auth.model';

@Component({
  selector: 'app-register',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './register.html'
})
export class Register {
  private readonly fb = inject(FormBuilder);
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  protected readonly formError = inject(FormErrorService);

  protected readonly registerForm: FormGroup;
  protected readonly isLoading = signal(false);
  protected readonly errorMessage = signal<string | null>(null);
  protected readonly showPassword = signal(false);
  protected readonly showConfirmPassword = signal(false);

  constructor() {
    this.registerForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(4)]],
      confirmPassword: ['', [Validators.required]]
    }, { validators: this.passwordMatchValidator });
  }

  private passwordMatchValidator(control: AbstractControl): ValidationErrors | null {
    const password = control.get('password');
    const confirmPassword = control.get('confirmPassword');

    if (!password || !confirmPassword) {
      return null;
    }

    return password.value === confirmPassword.value ? null : { passwordMismatch: true };
  }

  protected togglePasswordVisibility(): void {
    this.showPassword.update(value => !value);
  }

  protected toggleConfirmPasswordVisibility(): void {
    this.showConfirmPassword.update(value => !value);
  }

  protected onSubmit(): void {
    if (this.registerForm.invalid) {
      this.registerForm.markAllAsTouched();
      return;
    }

    this.isLoading.set(true);
    this.errorMessage.set(null);

    const { name, email, password } = this.registerForm.value;
    const registerData: RegisterRequest = { name, email, password };

    this.authService.register(registerData).subscribe({
      next: () => {
        this.isLoading.set(false);
        // Auto-login after registration
        this.authService.login(email, password).subscribe({
          next: () => {
            this.router.navigate(['/']);
          },
          error: () => {
            // If auto-login fails, redirect to login page
            this.router.navigate(['/login']);
          }
        });
      },
      error: (error) => {
        this.isLoading.set(false);
        if (error.status === 400) {
          this.errorMessage.set('Invalid registration data. Please check your inputs.');
        } else if (error.status === 409) {
          this.errorMessage.set('Email already exists. Please try another.');
        } else if (error.status === 0) {
          this.errorMessage.set('Unable to connect to the server. Please check your connection.');
        } else {
          this.errorMessage.set('An error occurred during registration. Please try again later.');
        }
        console.error('Registration error:', error);
      }
    });
  }

  protected getFieldError(fieldName: string): string | null {
    return this.formError.getFieldError(this.registerForm.get(fieldName), fieldName, {
      pattern: 'Username can only contain letters, numbers, hyphens, and underscores.',
    });
  }

  protected getPasswordMismatchError(): string | null {
    return this.formError.getFormError(this.registerForm, 'passwordMismatch');
  }
}
