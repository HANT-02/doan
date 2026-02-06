import { apiClient } from './client';
import type { LoginResponse, RefreshTokenResponse, MessageResponse } from '@/types/auth';

export const authApi = {
    login: (data: any) => apiClient.post<any, { data: LoginResponse }>('/login', data),
    logout: (token: string) => apiClient.post<any, { data: MessageResponse }>('/logout', { token }),
    refreshToken: (token: string) => apiClient.post<any, { data: RefreshTokenResponse }>('/refresh', { refresh_token: token }),
    register: (data: any) => apiClient.post('/register', data),
    forgotPassword: (email: string) => apiClient.post<any, { data: MessageResponse }>('/forgot-password', { email }),
    resetPassword: (data: any) => apiClient.post<any, { data: MessageResponse }>('/reset-password', data),
    changePassword: (data: any) => apiClient.post<any, { data: MessageResponse }>('/change-password', data),
};
