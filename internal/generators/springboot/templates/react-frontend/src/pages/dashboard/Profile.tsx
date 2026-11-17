import { useCallback, useEffect, useRef, useState } from 'react';
import { userService } from '../../services/userService';
import { useForm } from '../../hooks/useForm';
import { Validators } from '../../utils/validators';
import { HttpError } from '../../api/httpClient';
import { cn } from '../../utils/classNames';
import type { UpdateUserRequest, UserProfile } from '../../models/user';
import UserAvatar from '../../components/UserAvatar';

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

export default function Profile() {
  const [userProfile, setUserProfile] = useState<UserProfile | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isEditMode, setIsEditMode] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const successTimer = useRef<ReturnType<typeof setTimeout>>(undefined);

  const form = useForm({
    initialValues: { name: '' },
    validators: { name: [Validators.required, Validators.minLength(2)] },
  });
  const { reset } = form;

  const loadProfile = useCallback(async () => {
    setIsLoading(true);
    setErrorMessage(null);
    try {
      const profile = await userService.getUserProfile();
      setUserProfile(profile);
      reset({ name: profile.name });
    } catch (error) {
      setErrorMessage('Failed to load profile. Please try again.');
      console.error('Profile load error:', error);
    } finally {
      setIsLoading(false);
    }
  }, [reset]);

  useEffect(() => {
    void loadProfile();
    return () => clearTimeout(successTimer.current);
  }, [loadProfile]);

  const enableEditMode = () => {
    setIsEditMode(true);
    setSuccessMessage(null);
    setErrorMessage(null);
  };

  const cancelEdit = () => {
    setIsEditMode(false);
    if (userProfile) {
      reset({ name: userProfile.name });
    }
    setErrorMessage(null);
  };

  const handleSubmit = async () => {
    if (!form.isValid) {
      form.markAllAsTouched();
      return;
    }
    if (!userProfile) return;

    setIsSaving(true);
    setErrorMessage(null);
    setSuccessMessage(null);

    const updateData: UpdateUserRequest = { name: form.values.name };
    try {
      await userService.updateUserProfile(updateData);
      setIsEditMode(false);
      setSuccessMessage('Profile updated successfully!');
      await loadProfile();
      clearTimeout(successTimer.current);
      successTimer.current = setTimeout(() => setSuccessMessage(null), 3000);
    } catch (error) {
      const status = error instanceof HttpError ? error.status : -1;
      if (status === 403) {
        setErrorMessage('You do not have permission to update this profile.');
      } else if (status === 400) {
        setErrorMessage('Invalid profile data. Please check your inputs.');
      } else {
        setErrorMessage('Failed to update profile. Please try again.');
      }
      console.error('Profile update error:', error);
    } finally {
      setIsSaving(false);
    }
  };

  const nameError = form.getFieldError('name');

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        {isLoading && !userProfile && (
          <div className="flex justify-center items-center py-20" role="status">
            <div
              className="animate-spin rounded-full h-16 w-16 border-b-4 border-indigo-500"
              aria-label="Loading profile"
            ></div>
            <span className="sr-only">Loading profile...</span>
          </div>
        )}

        {!isLoading && errorMessage && !userProfile && (
          <div
            className="bg-red-50 border-2 border-red-200 rounded-lg p-6 text-center max-w-2xl mx-auto"
            role="alert"
          >
            <i
              className="fas fa-exclamation-circle text-red-500 text-3xl mb-4"
              aria-hidden="true"
            ></i>
            <p className="text-red-800 text-lg font-medium m-0">{errorMessage}</p>
          </div>
        )}

        {userProfile && (
          <div className="bg-white rounded-2xl shadow-lg overflow-hidden">
            <div className="bg-gradient-to-r from-indigo-600 to-purple-600 px-8 py-12 text-center relative">
              <div className="relative inline-block">
                <UserAvatar
                  fullName={userProfile.name}
                  size="2xl"
                  colorVariant="white"
                  ariaLabel="Profile avatar"
                />
              </div>
              <h1 className="text-3xl font-bold text-white mt-4 mb-2">{userProfile.name}</h1>
              <p className="text-indigo-100 text-lg">{userProfile.email}</p>
              <div className="mt-4">
                <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-white/20 text-white backdrop-blur-sm">
                  <i className="fas fa-shield-alt mr-2" aria-hidden="true"></i>
                  {userProfile.role}
                </span>
              </div>
            </div>

            {successMessage && (
              <div
                className="mx-8 mt-6 bg-green-50 border-2 border-green-200 rounded-lg p-4"
                role="alert"
                aria-live="polite"
              >
                <div className="flex items-center">
                  <i className="fas fa-check-circle text-green-500 mr-3" aria-hidden="true"></i>
                  <p className="text-sm text-green-800 m-0">{successMessage}</p>
                </div>
              </div>
            )}

            {errorMessage && (
              <div
                className="mx-8 mt-6 bg-red-50 border-2 border-red-200 rounded-lg p-4"
                role="alert"
                aria-live="polite"
              >
                <div className="flex items-center">
                  <i className="fas fa-exclamation-circle text-red-500 mr-3" aria-hidden="true"></i>
                  <p className="text-sm text-red-800 m-0">{errorMessage}</p>
                </div>
              </div>
            )}

            <div className="px-8 py-6 border-b border-gray-200">
              {!isEditMode ? (
                <button
                  type="button"
                  onClick={enableEditMode}
                  className="bg-indigo-600 text-white font-semibold px-6 py-2.5 rounded-lg transition-all duration-300 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 shadow-md hover:shadow-lg"
                >
                  <i className="fas fa-edit mr-2" aria-hidden="true"></i>
                  Edit Profile
                </button>
              ) : (
                <div className="flex gap-3">
                  <button
                    type="button"
                    onClick={handleSubmit}
                    disabled={isSaving}
                    className="bg-indigo-600 text-white font-semibold px-6 py-2.5 rounded-lg transition-all duration-300 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed shadow-md hover:shadow-lg"
                  >
                    {isSaving ? (
                      <span className="flex items-center">
                        <i className="fas fa-spinner fa-spin mr-2" aria-hidden="true"></i>
                        Saving...
                      </span>
                    ) : (
                      <span>
                        <i className="fas fa-save mr-2" aria-hidden="true"></i>
                        Save Changes
                      </span>
                    )}
                  </button>
                  <button
                    type="button"
                    onClick={cancelEdit}
                    disabled={isSaving}
                    className="bg-gray-200 text-gray-700 font-semibold px-6 py-2.5 rounded-lg transition-all duration-300 hover:bg-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-400 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Cancel
                  </button>
                </div>
              )}
            </div>

            <form
              className="px-8 py-6"
              noValidate
              onSubmit={(event) => {
                event.preventDefault();
                void handleSubmit();
              }}
            >
              <div className="space-y-6">
                <div>
                  <h2 className="text-xl font-semibold text-gray-900 mb-4">Basic Information</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                      <label
                        htmlFor="name"
                        className="block text-sm font-medium text-gray-700 mb-2"
                      >
                        Full Name
                      </label>
                      {isEditMode ? (
                        <>
                          <input
                            id="name"
                            type="text"
                            disabled={isSaving}
                            className={cn(
                              'block w-full px-4 py-2.5 border-2 rounded-lg text-gray-900 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all',
                              nameError ? 'border-red-300' : 'border-gray-300',
                            )}
                            aria-invalid={nameError ? 'true' : 'false'}
                            {...form.register('name')}
                          />
                          {nameError && <p className="mt-1 text-sm text-red-600">{nameError}</p>}
                        </>
                      ) : (
                        <p className="text-gray-900 py-2.5">{userProfile.name}</p>
                      )}
                    </div>

                    <div>
                      <span className="block text-sm font-medium text-gray-700 mb-2">
                        Email Address
                      </span>
                      <p className="text-gray-900 py-2.5">{userProfile.email}</p>
                      <p className="text-xs text-gray-500 mt-1">Email cannot be changed</p>
                    </div>
                  </div>
                </div>

                <div className="pt-6 border-t border-gray-200">
                  <h2 className="text-xl font-semibold text-gray-900 mb-4">Account Information</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-gray-600">Member Since:</span>
                      <span className="ml-2 font-medium text-gray-900">
                        {formatDate(userProfile.createdAt)}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </form>
          </div>
        )}
      </div>
    </div>
  );
}
