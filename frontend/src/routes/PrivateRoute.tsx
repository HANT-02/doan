import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';

export const PrivateRoute: React.FC = () => {
    const { isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return <div>Loading...</div>; // Or a proper loading spinner
    }

    return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};
