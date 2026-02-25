import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

interface User {
    id: number;
    code: string;
    full_name: string;
    email: string;
    role: string;
    is_active: boolean;
}

interface AuthState {
    user: User | null;
    role: string | null;
    accessToken: string | null;
    isAuthenticated: boolean;
    isLoadingAuth: boolean;
}

const getStoredUser = () => {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
};

const initialState: AuthState = {
    user: getStoredUser(),
    role: getStoredUser()?.role || null, // Should be uppercase as per BE
    accessToken: localStorage.getItem('accessToken') || null,
    isAuthenticated: !!localStorage.getItem('accessToken'),
    isLoadingAuth: true,
};

const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setCredentials: (
            state,
            action: PayloadAction<{ user: User; accessToken: string; refreshToken?: string }>
        ) => {
            const { user, accessToken, refreshToken } = action.payload;
            state.user = user;
            state.role = user.role; // Keep as BE sends (ADMIN, TEACHER, etc)
            state.accessToken = accessToken;
            state.isAuthenticated = true;

            localStorage.setItem('accessToken', accessToken);
            localStorage.setItem('user', JSON.stringify(user));
            if (refreshToken) {
                localStorage.setItem('refreshToken', refreshToken);
            }
        },
        logout: (state) => {
            state.user = null;
            state.role = null;
            state.accessToken = null;
            state.isAuthenticated = false;
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
            localStorage.removeItem('user');
        },
        setLoadingAuth: (state, action: PayloadAction<boolean>) => {
            state.isLoadingAuth = action.payload;
        },
        updateUser: (state, action: PayloadAction<User>) => {
            state.user = action.payload;
            state.role = action.payload.role;
            state.isAuthenticated = true;
            localStorage.setItem('user', JSON.stringify(action.payload));
        },
    },
});

export const { setCredentials, logout, setLoadingAuth, updateUser } = authSlice.actions;

export default authSlice.reducer;
