import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import type { RootState } from '../store';
import { logout } from '../store/authSlice';

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

// Create custom baseQuery that handles 401s auto-logout
const baseQueryWithReauth = async (args: any, api: any, extraOptions: any) => {
    let result = await baseQuery(args, api, extraOptions);

    if (result.error && result.error.status === 401) {
        // Dispatch logout on 401
        // Optional: implement silent refresh logic here if needed
        api.dispatch(logout());
    }
    return result;
};

// Initialize an empty api service that we'll inject endpoints into later as needed
export const baseApi = createApi({
    reducerPath: 'api',
    baseQuery: baseQueryWithReauth,
    endpoints: () => ({}),
    tagTypes: ['Teacher', 'User', 'Class', 'Room', 'Student']
});
