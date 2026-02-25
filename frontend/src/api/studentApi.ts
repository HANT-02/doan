import { baseApi } from './baseApi';

export interface Student {
    id: string;
    code: string;
    full_name: string;
    email: string;
    phone: string;
    guardian_phone: string;
    grade_level: string;
    status: 'ACTIVE' | 'INACTIVE';
    date_of_birth?: string;
    gender?: string;
    address?: string;
    created_at: string;
    updated_at: string;
}

export interface ListStudentsResponse {
    success: boolean;
    message: string;
    data: {
        students: Student[];
        pagination: {
            items_per_page: number;
            total_items: number;
            current_page: number;
            total_pages: number;
        };
    };
}

export interface StudentResponse {
    success: boolean;
    message: string;
    data: Student;
}

export const studentApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getStudents: builder.query<ListStudentsResponse, { page?: number; limit?: number; search?: string; status?: string; sortBy?: string; sortOrder?: string }>({
            query: (params) => ({
                url: '/v1/students',
                params,
            }),
            transformResponse: (response: any) => {
                return {
                    success: response.success,
                    message: response.message,
                    data: {
                        students: response.data?.Students || [],
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
                result?.data?.students
                    ? [
                        ...result.data.students.map(({ id }) => ({ type: 'Student' as const, id })),
                        { type: 'Student', id: 'LIST' },
                    ]
                    : [{ type: 'Student', id: 'LIST' }],
        }),
        getStudentById: builder.query<StudentResponse, string>({
            query: (id) => `/v1/students/${id}`,
            transformResponse: (response: any) => ({
                ...response,
                data: response.data?.Student || response.data
            }),
            providesTags: (_result, _error, id) => [{ type: 'Student', id }],
        }),
        createStudent: builder.mutation<StudentResponse, Partial<Student>>({
            query: (body) => ({
                url: '/v1/students',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'Student', id: 'LIST' }],
        }),
        updateStudent: builder.mutation<StudentResponse, { id: string; body: Partial<Student> }>({
            query: ({ id, body }) => ({
                url: `/v1/students/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Student', id }, { type: 'Student', id: 'LIST' }],
        }),
        deleteStudent: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/students/${id}`,
                method: 'DELETE',
                body: {},
            }),
            invalidatesTags: [{ type: 'Student', id: 'LIST' }],
        }),
    }),
});

export const {
    useGetStudentsQuery,
    useGetStudentByIdQuery,
    useCreateStudentMutation,
    useUpdateStudentMutation,
    useDeleteStudentMutation,
} = studentApi;
