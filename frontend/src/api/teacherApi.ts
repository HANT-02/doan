import { baseApi } from './baseApi';

export interface Teacher {
    id: string;
    code: string;
    full_name: string;
    email: string;
    phone?: string;
    status: string;
    employment_type: string;
    is_school_teacher?: boolean;
    school_name?: string;
    skills?: string[];
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
    message?: string;
    data: {
        teachers: Teacher[];
        pagination: {
            total_pages: number;
            total_items: number;
            current_page: number;
            items_per_page: number;
        };
    };
}

interface RawTeacherListResponse {
    success?: boolean;
    message?: string;
    data?: {
        teachers?: Teacher[];
        Teachers?: Teacher[];
        pagination?: {
            total_pages?: number;
            total_items?: number;
            current_page?: number;
            items_per_page?: number;
        };
        Pagination?: {
            total_pages?: number;
            total_items?: number;
            current_page?: number;
            items_per_page?: number;
            TotalPages?: number;
            TotalItems?: number;
            CurrentPage?: number;
            ItemsPerPage?: number;
        };
    };
}

export const teacherApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getTeachers: builder.query<ListTeachersResponse, ListTeachersParams>({
            query: (params) => ({
                url: 'v1/teachers',
                params,
            }),
            transformResponse: (response: RawTeacherListResponse) => ({
                success: response.success ?? false,
                message: response.message,
                data: {
                    teachers: response.data?.teachers || response.data?.Teachers || [],
                    pagination: {
                        total_pages: response.data?.pagination?.total_pages || response.data?.Pagination?.total_pages || response.data?.Pagination?.TotalPages || 1,
                        total_items: response.data?.pagination?.total_items || response.data?.Pagination?.total_items || response.data?.Pagination?.TotalItems || 0,
                        current_page: response.data?.pagination?.current_page || response.data?.Pagination?.current_page || response.data?.Pagination?.CurrentPage || 1,
                        items_per_page: response.data?.pagination?.items_per_page || response.data?.Pagination?.items_per_page || response.data?.Pagination?.ItemsPerPage || 10,
                    },
                },
            }),
            providesTags: (result) =>
                result?.data?.teachers?.length
                    ? [
                        ...result.data.teachers.map(({ id }) => ({ type: 'Teacher' as const, id })),
                        { type: 'Teacher', id: 'LIST' },
                    ]
                    : [{ type: 'Teacher', id: 'LIST' }],
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
            invalidatesTags: [{ type: 'Teacher', id: 'LIST' }],
        }),
        updateTeacher: builder.mutation({
            query: ({ id, ...body }) => ({
                url: `v1/teachers/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Teacher', id }, { type: 'Teacher', id: 'LIST' }],
        }),
        deleteTeacher: builder.mutation({
            query: (id) => ({
                url: `v1/teachers/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: [{ type: 'Teacher', id: 'LIST' }],
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
