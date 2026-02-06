export interface User {
    id: string;
    code: string;
    full_name: string;
    email: string;
    role: string;
    is_active: boolean;
}

export interface LoginResponse {
    access_token: string;
    refresh_token: string;
    user: User;
}

export interface RefreshTokenResponse {
    access_token: string;
}

export interface MessageResponse {
    message: string;
}

export interface ApiError {
    message: string;
    code?: string;
    details?: unknown;
}
