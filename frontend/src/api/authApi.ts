import { baseApi } from './baseApi';
import { setCredentials, logout } from '../store/authSlice';

export const authApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        login: builder.mutation({
            query: (credentials) => ({
                url: '/v1/auth/login',
                method: 'POST',
                body: credentials,
            }),
            async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
                try {
                    const { data } = await queryFulfilled;
                    if (data && data.success && data.data) {
                        dispatch(setCredentials({
                            user: data.data.user,
                            accessToken: data.data.access_token,
                            refreshToken: data.data.refresh_token,
                        }));
                    }
                } catch (error) {
                    // Handled in component
                }
            },
        }),
        logoutAccount: builder.mutation({
            query: (token) => ({
                url: '/v1/auth/logout',
                method: 'POST',
                body: { token },
            }),
            async onQueryStarted(_arg, { dispatch, queryFulfilled }) {
                try {
                    await queryFulfilled;
                } finally {
                    dispatch(logout());
                }
            }
        }),
        refreshToken: builder.mutation({
            query: (refreshToken) => ({
                url: '/v1/auth/refresh',
                method: 'POST',
                body: { refresh_token: refreshToken },
            }),
        }),
        register: builder.mutation({
            query: (body) => ({
                url: '/v1/auth/register',
                method: 'POST',
                body,
            }),
        }),
        forgotPassword: builder.mutation({
            query: (email) => ({
                url: '/v1/auth/forgot-password',
                method: 'POST',
                body: { email },
            }),
        }),
        resetPassword: builder.mutation({
            query: (body) => ({
                url: '/v1/auth/reset-password',
                method: 'POST',
                body,
            }),
        }),
        changePassword: builder.mutation({
            query: (body) => ({
                url: '/v1/auth/change-password',
                method: 'POST',
                body,
            }),
        }),
        getMe: builder.query({
            query: () => '/v1/auth/me',
            providesTags: ['User']
        }),
    }),
});

export const {
    useLoginMutation,
    useLogoutAccountMutation,
    useRefreshTokenMutation,
    useRegisterMutation,
    useForgotPasswordMutation,
    useResetPasswordMutation,
    useChangePasswordMutation,
    useGetMeQuery,
    useLazyGetMeQuery
} = authApi;
