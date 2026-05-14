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
    status: 'DRAFT' | 'PUBLISHED' | 'HISTORY' | 'UNPLANNED';
    published_at?: string;
    source_preview_run_id?: string;
    change_reason?: string;
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

export interface LessonAttendanceRecord {
    student: {
        id: string;
        code: string;
        full_name: string;
        email?: string;
    };
    attendance?: {
        id: string;
        lesson_id: string;
        student_id: string;
        status: number;
        note: string;
        marked_at: string;
        created_at: string;
        updated_at: string;
    } | null;
    status: number;
    note: string;
}

export interface LessonSummary {
    id: string;
    lesson_id: string;
    topic: string;
    lesson_content: string;
    class_feedback: string;
    homework: string;
    homework_deadline?: string;
    teacher_notes: string;
    created_by_id?: string;
    created_by?: {
        id: string;
        code: string;
        fullName?: string;
        full_name?: string;
        email: string;
    };
    created_at: string;
    updated_at: string;
}

export interface LessonAcademicRecord {
    id: string;
    lesson_summary_id: string;
    lesson_summary?: {
        id: string;
        lesson?: {
            id: string;
            class_id: string;
            class?: {
                id: string;
                code: string;
                name: string;
            };
            date_start: string;
            date_end: string;
        };
    };
    student_id: string;
    student: {
        id: string;
        code: string;
        full_name: string;
    };
    homework_completed: boolean;
    homework_score: number;
    attitude_rating: number;
    participation_score: number;
    personal_comment: string;
    total_score: number;
    is_completed: boolean;
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

export interface LessonAttendanceResponse {
    success: boolean;
    message: string;
    data: {
        lesson: Lesson;
        records: LessonAttendanceRecord[];
    };
}

export interface LessonSummaryResponse {
    success: boolean;
    message: string;
    data: {
        lesson: Lesson;
        summary: LessonSummary | null;
    };
}

export interface LessonAcademicRecordsResponse {
    success: boolean;
    message: string;
    data: {
        lesson: Lesson;
        records: Array<{
            student: {
                id: string;
                code: string;
                full_name: string;
            };
            record?: LessonAcademicRecord | null;
        }>;
    };
}

export interface ListLessonsParams {
    page?: number;
    limit?: number;
    class_id?: string;
    teacher_id?: string;
    status?: string;
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
        getLessonAttendance: builder.query<LessonAttendanceResponse, string>({
            query: (id) => `/v1/lessons/${id}/attendance`,
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    lesson: response?.data?.lesson,
                    records: response?.data?.records || [],
                },
            }),
            providesTags: (_result, _error, id) => [{ type: 'Lesson', id }, { type: 'Lesson', id: `ATTENDANCE-${id}` }],
        }),
        upsertLessonAttendance: builder.mutation<{ success: boolean; message: string; data?: { saved_count?: number } }, {
            id: string;
            body: {
                records: Array<{
                    student_id: string;
                    status: number;
                    note?: string;
                }>;
            };
        }>({
            query: ({ id, body }) => ({
                url: `/v1/lessons/${id}/attendance`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Lesson', id }, { type: 'Lesson', id: `ATTENDANCE-${id}` }],
        }),
        getLessonSummary: builder.query<LessonSummaryResponse, string>({
            query: (id) => `/v1/lessons/${id}/summary`,
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    lesson: response?.data?.lesson,
                    summary: response?.data?.summary || null,
                },
            }),
            providesTags: (_result, _error, id) => [{ type: 'Lesson', id }, { type: 'Lesson', id: `SUMMARY-${id}` }],
        }),
        upsertLessonSummary: builder.mutation<{ success: boolean; message: string }, {
            id: string;
            body: {
                topic: string;
                lesson_content: string;
                class_feedback: string;
                homework: string;
                homework_deadline?: string | null;
                teacher_notes: string;
            };
        }>({
            query: ({ id, body }) => ({
                url: `/v1/lessons/${id}/summary`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Lesson', id }, { type: 'Lesson', id: `SUMMARY-${id}` }],
        }),
        getLessonAcademicRecords: builder.query<LessonAcademicRecordsResponse, string>({
            query: (id) => `/v1/lessons/${id}/records`,
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    lesson: response?.data?.lesson,
                    records: response?.data?.records || [],
                },
            }),
            providesTags: (_result, _error, id) => [{ type: 'AcademicRecord', id: `LESSON-${id}` }],
        }),
        upsertLessonAcademicRecords: builder.mutation<{ success: boolean; message: string; data?: { saved_count?: number } }, {
            id: string;
            body: {
                records: Array<{
                    student_id: string;
                    homework_completed: boolean;
                    homework_score: number;
                    attitude_rating: number;
                    participation_score: number;
                    personal_comment: string;
                }>;
            };
        }>({
            query: ({ id, body }) => ({
                url: `/v1/lessons/${id}/records`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'AcademicRecord', id: `LESSON-${id}` }],
        }),
        finalizeLessonAcademicRecords: builder.mutation<{ success: boolean; message: string; data?: { finalized_count?: number } }, string>({
            query: (id) => ({
                url: `/v1/lessons/${id}/records/finalize`,
                method: 'POST',
            }),
            invalidatesTags: (_result, _error, id) => [{ type: 'AcademicRecord', id: `LESSON-${id}` }],
        }),
    }),
});

export const {
    useGetLessonsQuery,
    useGetLessonByIdQuery,
    useGetLessonAttendanceQuery,
    useUpsertLessonAttendanceMutation,
    useGetLessonSummaryQuery,
    useUpsertLessonSummaryMutation,
    useGetLessonAcademicRecordsQuery,
    useUpsertLessonAcademicRecordsMutation,
    useFinalizeLessonAcademicRecordsMutation,
} = lessonApi;
