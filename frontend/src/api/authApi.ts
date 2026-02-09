import { apiClient } from './client';
import type { LoginResponse as AuthResponse, RefreshTokenResponse, MessageResponse, User } from '@/types/auth';

export const authApi = {
    getMe: () => apiClient.get<any, { data: User }>('/auth/me'),
    login: (data: any) => apiClient.post<any, { data: AuthResponse }>('/auth/login', data),
    logout: (token: string) => apiClient.post<any, { data: MessageResponse }>('/auth/logout', { token }),
    refreshToken: (token: string) => apiClient.post<any, { data: RefreshTokenResponse }>('/auth/refresh', { refresh_token: token }),
    register: (data: any) => apiClient.post('/auth/register', data),
    forgotPassword: (email: string) => apiClient.post<any, { data: MessageResponse }>('/auth/forgot-password', { email }),
    resetPassword: (data: any) => apiClient.post<any, { data: MessageResponse }>('/auth/reset-password', data),
    changePassword: (data: any) => apiClient.post<any, { data: MessageResponse }>('/auth/change-password', data),
};
