import { baseApi } from './baseApi';
import { normalizeDetailEnvelope, normalizeListEnvelope } from './responseUtils';

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
    message?: string;
    data: {
        courses: Course[];
        pagination: {
            total_pages: number;
            total_items: number;
            current_page: number;
            items_per_page: number;
        };
    };
}

export interface CourseResponse {
    success: boolean;
    message?: string;
    data: Course | null;
}

interface RawCourseListResponse {
    success?: boolean;
    message?: string;
    data?: {
        courses?: Course[];
        Courses?: Course[];
        pagination?: {
            total_pages?: number;
            total_items?: number;
            current_page?: number;
            items_per_page?: number;
            limit?: number;
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

interface RawCourseResponse {
    success?: boolean;
    message?: string;
    data?: {
        Course?: Course;
    } | Course;
}

export const courseApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getCourses: builder.query<ListCoursesResponse, ListCoursesParams | void>({
            query: (params) => ({
                url: '/v1/courses',
                params: params || undefined,
            }),
            transformResponse: (response: RawCourseListResponse) => {
                const normalized = normalizeListEnvelope<Course>(response, ['courses', 'Courses']);
                return {
                    success: normalized.success,
                    message: normalized.message,
                    data: {
                        courses: normalized.items,
                        pagination: normalized.pagination,
                    },
                };
            },
            providesTags: ['Course'],
        }),
        getCourseById: builder.query<CourseResponse, string>({
            query: (id) => `/v1/courses/${id}`,
            transformResponse: (response: RawCourseResponse) => normalizeDetailEnvelope<Course>(response, ['Course']),
            providesTags: (_result, _error, id) => [{ type: 'Course', id }],
        }),
        createCourse: builder.mutation<CourseResponse, Partial<Course>>({
            query: (body) => ({
                url: '/v1/courses',
                method: 'POST',
                body,
            }),
            invalidatesTags: ['Course'],
        }),
        updateCourse: builder.mutation<CourseResponse, Partial<Course> & { id: string }>({
            query: ({ id, ...body }) => ({
                url: `/v1/courses/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Course', id }, 'Course'],
        }),
        deleteCourse: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/courses/${id}`,
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
