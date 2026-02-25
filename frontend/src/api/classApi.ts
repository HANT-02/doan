import { baseApi } from './baseApi';

export interface Class {
    id: string;
    code: string;
    name: string;
    notes?: string;
    start_date: string;
    end_date?: string;
    max_students: number;
    status: 'OPEN' | 'CLOSED' | 'CANCELLED';
    price: number;
    program_id?: string;
    course_id?: string;
    teacher_id?: string;
    created_at: string;
    updated_at: string;
}

export interface ListClassesResponse {
    success: boolean;
    message: string;
    data: {
        classes: Class[];
        pagination: {
            items_per_page: number;
            total_items: number;
            current_page: number;
            total_pages: number;
        };
    };
}

export interface ClassResponse {
    success: boolean;
    message: string;
    data: Class;
}

export const classApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getClasses: builder.query<ListClassesResponse, { page?: number; limit?: number; search?: string; status?: string; sortBy?: string; sortOrder?: string }>({
            query: (params) => ({
                url: '/v1/classes',
                params,
            }),
            transformResponse: (response: any) => {
                return {
                    success: response.success,
                    message: response.message,
                    data: {
                        classes: response.data?.Classes || [],
                        pagination: {
                            items_per_page: response.data?.Pagination?.ItemsPerPage || 10,
                            total_items: response.data?.Pagination?.TotalItems || 0,
                            current_page: response.data?.Pagination?.CurrentPage || 1,
                            total_pages: response.data?.Pagination?.TotalPages || 1,
                        }
                    }
                };
            },
            providesTags: (result) =>
                result?.data?.classes
                    ? [
                        ...result.data.classes.map(({ id }) => ({ type: 'Class' as const, id })),
                        { type: 'Class', id: 'LIST' },
                    ]
                    : [{ type: 'Class', id: 'LIST' }],
        }),
        getClassById: builder.query<ClassResponse, string>({
            query: (id) => `/v1/classes/${id}`,
            transformResponse: (response: any) => ({
                ...response,
                data: response.data?.Class || response.data
            }),
            providesTags: (_result, _error, id) => [{ type: 'Class', id }],
        }),
        createClass: builder.mutation<ClassResponse, Partial<Class>>({
            query: (body) => ({
                url: '/v1/classes',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'Class', id: 'LIST' }],
        }),
        updateClass: builder.mutation<ClassResponse, { id: string; body: Partial<Class> }>({
            query: ({ id, body }) => ({
                url: `/v1/classes/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Class', id }, { type: 'Class', id: 'LIST' }],
        }),
        deleteClass: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/classes/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: [{ type: 'Class', id: 'LIST' }],
        }),
    }),
});

export const {
    useGetClassesQuery,
    useGetClassByIdQuery,
    useCreateClassMutation,
    useUpdateClassMutation,
    useDeleteClassMutation,
} = classApi;
