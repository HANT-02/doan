import { baseApi } from './baseApi';

export interface Teacher {
    id: string;
    code: string;
    full_name: string;
    email: string;
    phone?: string;
    status: string;
    employment_type: string;
    specialization?: string;
    bio?: string;
    is_school_teacher?: boolean;
    school_name?: string;
    notes?: string;
    created_at: string;
    updated_at: string;
}

export interface ListTeachersParams {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    employment_type?: string;
}

export interface ListTeachersResponse {
    success: boolean;
    data: {
        teachers: Teacher[];
        pagination: {
            total_pages: number;
            total_items: number;
            current_page: number;
            limit: number;
        };
    };
    message?: string;
}

export const teacherApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getTeachers: builder.query({
            query: (params) => ({
                url: 'v1/teachers',
                params, // Should pass page, limit, search as standard params
            }),
            providesTags: ['Teacher'],
        }),
        getTeacherById: builder.query({
            query: (id) => `v1/teachers/${id}`,
            providesTags: (_result, _error, id) => [{ type: 'Teacher', id }],
        }),
        createTeacher: builder.mutation({
            query: (body) => ({
                url: 'v1/teachers',
                method: 'POST',
                body,
            }),
            invalidatesTags: ['Teacher'],
        }),
        updateTeacher: builder.mutation({
            query: ({ id, ...body }) => ({
                url: `v1/teachers/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Teacher', id }, 'Teacher'],
        }),
        deleteTeacher: builder.mutation({
            query: (id) => ({
                url: `v1/teachers/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: ['Teacher'],
        }),
    }),
});

export const {
    useGetTeachersQuery,
    useGetTeacherByIdQuery,
    useCreateTeacherMutation,
    useUpdateTeacherMutation,
    useDeleteTeacherMutation,
} = teacherApi;
