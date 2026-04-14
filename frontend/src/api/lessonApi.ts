import { baseApi } from './baseApi';

export interface Lesson {
    id: string;
    class_id: string;
    class?: {
        id: string;
        code: string;
        name: string;
    };
    teacher_id?: string;
    teacher?: {
        id: string;
        code: string;
        full_name: string;
    };
    date_start: string;
    date_end: string;
    room_id?: string;
    room?: {
        id: string;
        code?: string;
        name: string;
    };
    notes: string;
    created_at: string;
    updated_at: string;
}

export interface ListLessonsResponse {
    success: boolean;
    message: string;
    data: {
        lessons: Lesson[];
        pagination: {
            items_per_page: number;
            total_items: number;
            current_page: number;
            total_pages: number;
        };
    };
}

export interface LessonResponse {
    success: boolean;
    message: string;
    data: Lesson;
}

export interface ListLessonsParams {
    page?: number;
    limit?: number;
    class_id?: string;
    teacher_id?: string;
    date_from?: string;
    date_to?: string;
    sortBy?: string;
    sortOrder?: string;
}

export const lessonApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getLessons: builder.query<ListLessonsResponse, ListLessonsParams>({
            query: (params) => ({
                url: '/v1/lessons',
                params,
            }),
            transformResponse: (response: any) => {
                const data = response?.data || {};
                const lessons = data?.Lessons || data?.lessons || [];
                const pagination = data?.Pagination || data?.pagination || {};
                return {
                    success: response?.success ?? false,
                    message: response?.message ?? '',
                    data: {
                        lessons,
                        pagination: {
                            items_per_page: pagination?.items_per_page ?? pagination?.ItemsPerPage ?? 10,
                            total_items: pagination?.total_items ?? pagination?.TotalItems ?? 0,
                            current_page: pagination?.current_page ?? pagination?.CurrentPage ?? 1,
                            total_pages: pagination?.total_pages ?? pagination?.TotalPages ?? 1,
                        },
                    },
                };
            },
            providesTags: (result) =>
                result?.data?.lessons
                    ? [
                        ...result.data.lessons.map(({ id }) => ({ type: 'Lesson' as const, id })),
                        { type: 'Lesson', id: 'LIST' },
                    ]
                    : [{ type: 'Lesson', id: 'LIST' }],
        }),
        getLessonById: builder.query<LessonResponse, string>({
            query: (id) => `/v1/lessons/${id}`,
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: response?.data || response,
            }),
            providesTags: (_result, _error, id) => [{ type: 'Lesson', id }],
        }),
    }),
});

export const {
    useGetLessonsQuery,
    useGetLessonByIdQuery,
} = lessonApi;
