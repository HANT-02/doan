
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import FullScreenLoader from '@/components/common/FullScreenLoader'; // Import the new component

export const PrivateRoute = () => {
    const { isAuthenticated, isLoadingAuth } = useAuth();

    if (isLoadingAuth) {
        return <FullScreenLoader />;
    }

    // Attempt to keep the user on the requested page if feasible, otherwise redirect
    // Use 'replace' to avoid history stack issues
    return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};
