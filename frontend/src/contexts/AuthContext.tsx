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
            const storedUser = localStorage.getItem('user');

            if (storedRefreshToken && storedUser) {
                try {
                    setUser(JSON.parse(storedUser));
                    // Optional: Refresh token immediately to ensure validity and get fresh access token
                    const { data } = await authApi.refreshToken(storedRefreshToken);
                    setAccessToken(data.access_token);
                    setAuthToken(data.access_token);
                } catch (error) {
                    console.error("Session expired:", error);
                    handleLogoutCleanup();
                }
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
            setUser(data.user);
            setAccessToken(data.access_token);
            setAuthToken(data.access_token);

            localStorage.setItem('refresh_token', data.refresh_token);
            localStorage.setItem('user', JSON.stringify(data.user));
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
