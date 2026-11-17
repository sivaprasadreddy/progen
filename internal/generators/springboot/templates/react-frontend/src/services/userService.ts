import { http } from '../api/httpClient';
import type { UpdateUserRequest, UserProfile } from '../models/user';

export const userService = {
  getUserProfile(): Promise<UserProfile> {
    return http.get<UserProfile>('/api/auth/me');
  },

  updateUserProfile(data: UpdateUserRequest): Promise<void> {
    return http.put<void>('/api/users/me', data);
  },
};
