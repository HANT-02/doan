import { Navigate, Outlet } from 'react-router-dom';
import { useAppSelector } from '@/store';
import FullScreenLoader from '@/components/common/FullScreenLoader';

export const PrivateRoute = () => {
    const { isAuthenticated, isLoadingAuth } = useAppSelector((state) => state.auth);

    if (isLoadingAuth) {
        return <FullScreenLoader />;
    }

    return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};
