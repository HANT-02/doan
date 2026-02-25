import React, { createContext, useContext, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store';
import { useLazyGetMeQuery } from '@/api/authApi';
import { setLoadingAuth, updateUser } from '@/store/authSlice';

interface AuthContextType {
    user: any;
    isAuthenticated: boolean;
    isLoadingAuth: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const dispatch = useAppDispatch();
    const { user, isAuthenticated, isLoadingAuth } = useAppSelector((state) => state.auth);
    const [getMe] = useLazyGetMeQuery();

    useEffect(() => {
        const bootstrapAuth = async () => {
            const token = localStorage.getItem('accessToken');

            if (token) {
                try {
                    // Fetch fresh user profile
                    const result = await getMe(undefined);
                    if (result.data?.success) {
                        dispatch(updateUser(result.data.data));
                    }
                    // If result.error or 401, baseQuery's auto-logout will handle it IF it's a 401
                    // Otherwise keep current local state
                } catch (err) {
                    console.error('Auth bootstrap failed:', err);
                }
            }
            dispatch(setLoadingAuth(false));
        };

        bootstrapAuth();
    }, [dispatch, getMe]);

    return (
        <AuthContext.Provider value={{
            user,
            isAuthenticated,
            isLoadingAuth
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
