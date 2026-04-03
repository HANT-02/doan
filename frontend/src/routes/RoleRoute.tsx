import { Navigate, Outlet } from 'react-router-dom';
import { useAppSelector } from '@/store';
import FullScreenLoader from '@/components/common/FullScreenLoader';
import { hasAnyRole } from '@/utils/roles';

interface RoleRouteProps {
    allowedRoles: string[];
}

export const RoleRoute: React.FC<RoleRouteProps> = ({ allowedRoles }) => {
    const { user, isLoadingAuth } = useAppSelector((state) => state.auth);

    if (isLoadingAuth) {
        return <FullScreenLoader />;
    }

    if (!user) {
        return <Navigate to="/login" replace />;
    }

    if (!hasAnyRole(user.role, allowedRoles)) {
        return <Navigate to="/403" replace />;
    }

    return <Outlet />;
};
