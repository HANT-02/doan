import { baseApi } from './baseApi';

export interface Course {
    id: string;
    code: string;
    name: string;
    description?: string;
    grade_level?: string;
    subject?: string;
    session_count?: number;
    session_duration_minutes?: number;
    total_hours?: number;
    price?: number;
    status: string;
    created_at: string;
    updated_at: string;
}

export interface ListCoursesParams {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    grade_level?: string;
    subject?: string;
}

export interface ListCoursesResponse {
    success: boolean;
    data: {
        courses: Course[];
        pagination: {
            total_pages: number;
            total_items: number;
            current_page: number;
            limit: number;
        };
    };
    message?: string;
}

export const courseApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getCourses: builder.query<ListCoursesResponse, ListCoursesParams | void>({
            query: (params) => ({
                url: 'v1/courses',
                params: params || undefined,
            }),
            providesTags: ['Course'],
        }),
        getCourseById: builder.query<{ success: boolean; data: Course }, string>({
            query: (id) => `v1/courses/${id}`,
            providesTags: (_result, _error, id) => [{ type: 'Course', id }],
        }),
        createCourse: builder.mutation<{ success: boolean; data: Course }, Partial<Course>>({
            query: (body) => ({
                url: 'v1/courses',
                method: 'POST',
                body,
            }),
            invalidatesTags: ['Course'],
        }),
        updateCourse: builder.mutation<{ success: boolean; data: Course }, Partial<Course> & { id: string }>({
            query: ({ id, ...body }) => ({
                url: `v1/courses/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Course', id }, 'Course'],
        }),
        deleteCourse: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `v1/courses/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: ['Course'],
        }),
    }),
});

export const {
    useGetCoursesQuery,
    useGetCourseByIdQuery,
    useCreateCourseMutation,
    useUpdateCourseMutation,
    useDeleteCourseMutation,
} = courseApi;
