import React, { createContext, useContext, useState, useEffect } from 'react';
import type { User } from '@/types/auth';
import { authApi } from '@/api/authApi';
import { setAuthToken } from '@/api/client';

interface AuthContextType {
    user: User | null;
    accessToken: string | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (data: any) => Promise<void>;
    register: (data: any) => Promise<void>;
    logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    // Initialize auth state from local storage
    useEffect(() => {
        const initAuth = async () => {
            const storedRefreshToken = localStorage.getItem('refresh_token');

            if (storedRefreshToken) {
                try {
                    // 1. Refresh token immediately to ensure validity and get fresh access token
                    const { data: refreshData } = await authApi.refreshToken(storedRefreshToken);
                    const newAccessToken = refreshData.access_token;

                    setAccessToken(newAccessToken);
                    setAuthToken(newAccessToken);

                    // 2. Fetch fresh user data from server
                    const { data: userData } = await authApi.getMe();
                    const normalizedUser = { ...userData, role: userData.role.toLowerCase() };
                    setUser(normalizedUser);

                    localStorage.setItem('user', JSON.stringify(normalizedUser));
                } catch (error) {
                    console.error("Auth bootstrap failed:", error);
                    handleLogoutCleanup();
                }
            } else {
                handleLogoutCleanup();
            }
            setIsLoading(false);
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
            throw error;
        }
    };

    const register = async (data: any) => {
        await authApi.register(data);
    };

    const logout = async () => {
        try {
            if (accessToken) {
                await authApi.logout(accessToken);
            }
        } catch (error) {
            console.error("Logout error", error);
        } finally {
            handleLogoutCleanup();
        }
    };

    return (
        <AuthContext.Provider value={{
            user,
            accessToken,
            isAuthenticated: !!user,
            isLoading,
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
