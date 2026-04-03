const ADMIN_ROLES = ['ADMIN', 'SUPER_ADMIN'] as const;

export const normalizeRole = (role?: string | null) => role?.trim().toUpperCase() || '';

export const hasAnyRole = (role: string | null | undefined, allowedRoles: readonly string[]) => {
    const normalizedRole = normalizeRole(role);
    return allowedRoles.some((allowedRole) => normalizeRole(allowedRole) === normalizedRole);
};

export const isAdminRole = (role?: string | null) => hasAnyRole(role, ADMIN_ROLES);

