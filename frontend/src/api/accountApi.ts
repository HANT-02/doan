import { baseApi } from './baseApi';
import { normalizeDetailEnvelope, normalizeListEnvelope } from './responseUtils';

export interface AccountUser {
    id: string;
    code: string;
    fullName?: string;
    full_name?: string;
    email: string;
    role: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface NormalizedAccountUser {
    id: string;
    code: string;
    full_name: string;
    email: string;
    role: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface ListUsersResponse {
    success: boolean;
    message: string;
    data: {
        users: NormalizedAccountUser[];
        pagination: {
            items_per_page: number;
            total_items: number;
            current_page: number;
            total_pages: number;
        };
    };
}

export interface UserDetailResponse {
    success: boolean;
    message: string;
    data: {
        user: NormalizedAccountUser | null;
    };
}

const normalizeUser = (user: AccountUser): NormalizedAccountUser => ({
    id: user.id,
    code: user.code,
    full_name: user.full_name ?? user.fullName ?? '',
    email: user.email,
    role: user.role,
    is_active: user.is_active,
    created_at: user.created_at,
    updated_at: user.updated_at,
});

export const accountApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getUsers: builder.query<ListUsersResponse, {
            page?: number;
            limit?: number;
            search?: string;
            role?: string;
            is_active?: boolean;
        }>({
            query: (params) => ({
                url: '/v1/users',
                params,
            }),
            transformResponse: (response: any) => {
                const normalized = normalizeListEnvelope<AccountUser>(response, ['users', 'Users']);
                return {
                    success: normalized.success,
                    message: normalized.message,
                    data: {
                        users: normalized.items.map(normalizeUser),
                        pagination: normalized.pagination,
                    },
                };
            },
            providesTags: (result) =>
                result?.data?.users
                    ? [
                        ...result.data.users.map(({ id }) => ({ type: 'User' as const, id })),
                        { type: 'User', id: 'LIST' },
                    ]
                    : [{ type: 'User', id: 'LIST' }],
        }),
        getUserById: builder.query<UserDetailResponse, string>({
            query: (id) => `/v1/users/${id}`,
            transformResponse: (response: any) => {
                const normalized = normalizeDetailEnvelope<AccountUser>(response, ['user', 'User']);
                return {
                    success: normalized.success,
                    message: normalized.message,
                    data: {
                        user: normalized.data ? normalizeUser(normalized.data) : null,
                    },
                };
            },
            providesTags: (_result, _error, id) => [{ type: 'User', id }],
        }),
        createUser: builder.mutation<UserDetailResponse, {
            code?: string;
            full_name: string;
            email: string;
            role: string;
            is_active: boolean;
            password: string;
        }>({
            query: (body) => ({
                url: '/v1/users',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'User', id: 'LIST' }],
        }),
        updateUser: builder.mutation<{ success: boolean; message: string }, {
            id: string;
            body: {
                full_name?: string;
                role?: string;
                is_active?: boolean;
            };
        }>({
            query: ({ id, body }) => ({
                url: `/v1/users/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'User', id }, { type: 'User', id: 'LIST' }],
        }),
        resetUserPassword: builder.mutation<{ success: boolean; message: string }, {
            id: string;
            new_password: string;
        }>({
            query: ({ id, new_password }) => ({
                url: `/v1/users/${id}/reset-password`,
                method: 'POST',
                body: { new_password },
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'User', id }],
        }),
    }),
});

export const {
    useGetUsersQuery,
    useLazyGetUserByIdQuery,
    useCreateUserMutation,
    useUpdateUserMutation,
    useResetUserPasswordMutation,
} = accountApi;
