import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type { BaseQueryFn, FetchArgs, FetchBaseQueryError } from '@reduxjs/toolkit/query';
import type { RootState } from '../store';
import { logout, setAccessToken } from '../store/authSlice';

// Create base query with auth header injection
const baseQuery = fetchBaseQuery({
    baseUrl: 'http://localhost:9000/api', // Adjust base URL as needed based on your vite proxy or BE port
    prepareHeaders: (headers, { getState }) => {
        // By default, if we have a token in the store, let's use that for authenticated requests
        const token = (getState() as RootState).auth.accessToken;
        if (token) {
            headers.set('authorization', `Bearer ${token}`);
        }
        return headers;
    },
});

const getRequestUrl = (args: string | FetchArgs) => (typeof args === 'string' ? args : args.url);

const isRefreshRequest = (args: string | FetchArgs) => getRequestUrl(args).includes('/v1/auth/refresh');

// Create custom baseQuery that handles 401s auto-logout
const baseQueryWithReauth: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
    args,
    api,
    extraOptions
) => {
    let result = await baseQuery(args, api, extraOptions);

    if (result.error && result.error.status === 401 && !isRefreshRequest(args)) {
        const refreshToken = localStorage.getItem('refreshToken');

        if (!refreshToken) {
            api.dispatch(logout());
            return result;
        }

        const refreshResult = await baseQuery(
            {
                url: '/v1/auth/refresh',
                method: 'POST',
                body: { refresh_token: refreshToken },
            },
            api,
            extraOptions,
        );

        const refreshedAccessToken =
            typeof refreshResult.data === 'object' && refreshResult.data && 'data' in refreshResult.data
                ? (refreshResult.data as { data?: { access_token?: string } }).data?.access_token
                : undefined;

        if (!refreshResult.error && refreshedAccessToken) {
            api.dispatch(setAccessToken(refreshedAccessToken));
            result = await baseQuery(args, api, extraOptions);
        } else {
            api.dispatch(logout());
        }
    }

    return result;
};

// Initialize an empty api service that we'll inject endpoints into later as needed
export const baseApi = createApi({
    reducerPath: 'api',
    baseQuery: baseQueryWithReauth,
    endpoints: () => ({}),
    tagTypes: ['Teacher', 'User', 'Class', 'Room', 'Shift', 'Student', 'Program', 'Course', 'Scheduling', 'Audit', 'Material', 'Dashboard']
});
