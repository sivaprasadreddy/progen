import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import {UpdateUserRequest, UserProfile} from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private readonly http = inject(HttpClient);

  getUserProfile(): Observable<UserProfile> {
    return this.http.get<UserProfile>(`/api/auth/me`);
  }

  updateUserProfile(data: UpdateUserRequest): Observable<void> {
    return this.http.put<void>(`/api/users/me`, data);
  }
}
