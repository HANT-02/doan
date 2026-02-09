import axios from 'axios';

// Base URL from env or default to proxy
const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

export const apiClient = axios.create({
    baseURL: BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const setAuthToken = (token: string | null) => {
    if (token) {
        apiClient.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    } else {
        delete apiClient.defaults.headers.common['Authorization'];
    }
};

// Response interceptor for handling 401 and errors
apiClient.interceptors.response.use(
    (response) => response.data, // Unwrap data
    async (error) => {
        // Handle 401 if needed (refresh token logic will be in Context or here)
        // For now just reject
        const message = error.response?.data?.message || error.message || 'Unknown error';
        return Promise.reject(new Error(message));
    }
);
