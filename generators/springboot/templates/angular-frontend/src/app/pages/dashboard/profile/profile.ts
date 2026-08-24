import { Component, ChangeDetectionStrategy, signal, inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from '../../../services/auth.service';
import { UserService } from '../../../services/user.service';
import { FormErrorService } from '../../../services/form-error.service';
import { UserAvatar } from '../../../components/user-avatar/user-avatar';
import {UpdateUserRequest, UserProfile} from '../../../models/user.model';

@Component({
  selector: 'app-profile',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [ReactiveFormsModule, UserAvatar],
  templateUrl: './profile.html'
})
export class Profile implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly authService = inject(AuthService);
  private readonly userService = inject(UserService);
  private readonly router = inject(Router);
  protected readonly formError = inject(FormErrorService);

  protected readonly userProfile = signal<UserProfile | null>(null);
  protected readonly isLoading = signal(true);
  protected readonly isSaving = signal(false);
  protected readonly isEditMode = signal(false);
  protected readonly errorMessage = signal<string | null>(null);
  protected readonly successMessage = signal<string | null>(null);

  protected readonly profileForm: FormGroup;

  constructor() {
    this.profileForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
    });

    // Disable form initially
    this.profileForm.disable();
  }

  ngOnInit(): void {
    const currentUser = this.authService.user();
    if (!currentUser) {
      this.router.navigate(['/login']);
      return;
    }

    this.loadProfile();
  }

  private loadProfile(): void {
    this.isLoading.set(true);
    this.errorMessage.set(null);

    this.userService.getUserProfile().subscribe({
      next: (profile) => {
        this.userProfile.set(profile);
        this.populateForm(profile);
        this.isLoading.set(false);
      },
      error: (error) => {
        this.isLoading.set(false);
        this.errorMessage.set('Failed to load profile. Please try again.');
        console.error('Profile load error:', error);
      }
    });
  }

  private populateForm(profile: UserProfile): void {
    this.profileForm.patchValue({
      name: profile.name,
    });
  }

  protected enableEditMode(): void {
    this.isEditMode.set(true);
    this.profileForm.enable();
    this.successMessage.set(null);
    this.errorMessage.set(null);
  }

  protected cancelEdit(): void {
    this.isEditMode.set(false);
    this.profileForm.disable();
    const profile = this.userProfile();
    if (profile) {
      this.populateForm(profile);
    }
    this.errorMessage.set(null);
  }

  protected onSubmit(): void {
    if (this.profileForm.invalid) {
      this.profileForm.markAllAsTouched();
      return;
    }

    const profile = this.userProfile();
    if (!profile) return;

    this.isSaving.set(true);
    this.errorMessage.set(null);
    this.successMessage.set(null);

    const formValue = this.profileForm.getRawValue();
    const updateData: UpdateUserRequest = {
      name: formValue.name
    };

    this.userService.updateUserProfile(updateData).subscribe({
      next: () => {
        this.isSaving.set(false);
        this.isEditMode.set(false);
        this.profileForm.disable();
        this.successMessage.set('Profile updated successfully!');

        // Reload profile to get fresh data
        this.loadProfile();

        // Clear success message after 3 seconds
        setTimeout(() => {
          this.successMessage.set(null);
        }, 3000);
      },
      error: (error) => {
        this.isSaving.set(false);
        if (error.status === 403) {
          this.errorMessage.set('You do not have permission to update this profile.');
        } else if (error.status === 400) {
          this.errorMessage.set('Invalid profile data. Please check your inputs.');
        } else {
          this.errorMessage.set('Failed to update profile. Please try again.');
        }
        console.error('Profile update error:', error);
      }
    });
  }

  protected getFieldError(fieldName: string): string | null {
    return this.formError.getFieldError(this.profileForm.get(fieldName), fieldName, {
      pattern: 'Please enter a valid phone number.',
    });
  }

  protected formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }
}
