import React, { createContext, useContext, useState, useEffect } from 'react';
import type { User } from '@/types/auth';
import { authApi } from '@/api/authApi';
import { setAuthToken } from '@/api/client';

interface AuthContextType {
    user: User | null;
    accessToken: string | null;
    isAuthenticated: boolean;
    isLoadingAuth: boolean; // Renamed from isLoading
    login: (data: any) => Promise<void>;
    register: (data: any) => Promise<void>;
    logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const [isLoadingAuth, setIsLoadingAuth] = useState(true); // Renamed state variable

    // Initialize auth state from local storage and bootstrap session
    useEffect(() => {
        const initAuth = async () => {
            const storedRefreshToken = localStorage.getItem('refresh_token');
            const storedUser = localStorage.getItem('user');

            if (storedRefreshToken) {
                try {
                    // Try to refresh the access token first
                    const { data: refreshData } = await authApi.refreshToken(storedRefreshToken);
                    const newAccessToken = refreshData.access_token;
                    setAccessToken(newAccessToken);
                    setAuthToken(newAccessToken);

                    // If refresh is successful, fetch the latest user data to ensure sync
                    const { data: userData } = await authApi.getMe();
                    const normalizedUser = { ...userData, role: userData.role.toLowerCase() };
                    setUser(normalizedUser);
                    localStorage.setItem('user', JSON.stringify(normalizedUser));
                    localStorage.setItem('refresh_token', storedRefreshToken);
                } catch (error: any) {
                    console.error("Auth bootstrap failed:", error);

                    // If refresh token is expired or invalid (401/403), we must log out
                    if (error.response && (error.response.status === 401 || error.response.status === 403)) {
                        handleLogoutCleanup();
                    } else if (storedUser) {
                        // For transient network errors, if we have a stored user, we could potentially
                        // "optimistically" keep them logged in but marked as "stale" or just null out.
                        // Best practice for strict apps: null out if we can't verify session.
                        setUser(null);
                    }
                }
            } else {
                handleLogoutCleanup();
            }
            setIsLoadingAuth(false);
        };

        initAuth();
    }, []);

    const handleLogoutCleanup = () => {
        setUser(null);
        setAccessToken(null);
        setAuthToken(null);
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
    };

    const login = async (credentials: any) => {
        try {
            const { data } = await authApi.login(credentials);
            const normalizedUser = { ...data.user, role: data.user.role.toLowerCase() };

            setUser(normalizedUser);
            setAccessToken(data.access_token);
            setAuthToken(data.access_token);

            localStorage.setItem('refresh_token', data.refresh_token);
            localStorage.setItem('user', JSON.stringify(normalizedUser));
        } catch (error) {
            console.error("Login failed:", error);
            throw error;
        }
    };

    const register = async (data: any) => {
        try {
            await authApi.register(data);
        } catch (error) {
            console.error("Registration failed:", error);
            throw error;
        }
    };

    const logout = async () => {
        try {
            // Best effort logout on server
            const rt = localStorage.getItem('refresh_token');
            if (rt) {
                await authApi.logout(rt);
            }
        } catch (error) {
            console.error("Logout server-side error", error);
        } finally {
            handleLogoutCleanup();
        }
    };

    return (
        <AuthContext.Provider value={{
            user,
            accessToken,
            isAuthenticated: !!user,
            isLoadingAuth,
            login,
            register,
            logout
        }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
