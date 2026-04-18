import { baseApi } from './baseApi';

export interface TeacherLessonShift {
    id: string;
    code: string;
    name: string;
    start_time: string;
    end_time: string;
    duration_minutes: number;
    session_type: string;
}

export interface TeacherLesson {
    id: string;
    class_id: string;
    class_name: string;
    class_code: string;
    room_id?: string;
    room_name?: string;
    date_start: string;
    date_end: string;
    notes: string;
    shift?: TeacherLessonShift | null;
}

export interface TeacherLessonsResponse {
    success: boolean;
    message: string;
    data: {
        teacher_id: string;
        lessons: TeacherLesson[];
    };
}

export interface TeacherLessonAttendanceRecord {
    attendance_id?: string;
    student: {
        id: string;
        code: string;
        full_name: string;
    };
    status: number;
    note: string;
    marked_at?: string;
}

export interface TeacherLessonAttendanceResponse {
    success: boolean;
    message: string;
    data: {
        lesson: TeacherLesson;
        records: TeacherLessonAttendanceRecord[];
    };
}

export interface TeacherAttendanceSummaryStudent {
    student: {
        id: string;
        code: string;
        full_name: string;
    };
    total_lessons: number;
    marked_count: number;
    present_count: number;
    absent_count: number;
    late_count: number;
    excused_count: number;
    unmarked_count: number;
    attendance_rate: number;
}

export interface TeacherAttendanceSummaryResponse {
    success: boolean;
    message: string;
    data: {
        teacher_id: string;
        class_id: string;
        total_lessons: number;
        students: TeacherAttendanceSummaryStudent[];
    };
}

export interface TeacherLessonSummary {
    id: string;
    lesson_id: string;
    topic: string;
    lesson_content: string;
    class_feedback: string;
    homework: string;
    homework_deadline?: string;
    teacher_notes: string;
    created_by_id?: string;
    created_at: string;
    updated_at: string;
}

export interface TeacherLessonSummaryResponse {
    success: boolean;
    message: string;
    data: {
        lesson: TeacherLesson;
        summary: TeacherLessonSummary | null;
    };
}

export interface TeacherLeaveRequest {
    id: string;
    student: {
        id: string;
        code: string;
        full_name: string;
    };
    leave_type: string;
    apply_date: string;
    late_minutes: number;
    early_minutes: number;
    reason: string;
    documents: string[];
    class?: {
        id: string;
        code: string;
        name: string;
    } | null;
    lesson?: {
        id: string;
        date_start: string;
        date_end: string;
    } | null;
    subject: string;
    status: string;
    approved_by_id?: string;
    approved_at?: string;
    rejection_reason: string;
    created_at: string;
    updated_at: string;
}

export interface TeacherLeaveRequestsResponse {
    success: boolean;
    message: string;
    data: {
        requests: TeacherLeaveRequest[];
    };
}

export interface GetTeacherLessonsParams {
    class_id?: string;
    from?: string;
    to?: string;
}

export interface GetTeacherLeaveRequestsParams {
    class_id?: string;
    status?: string;
    student_id?: string;
}

export const teacherPortalApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getTeacherLessons: builder.query<TeacherLessonsResponse, GetTeacherLessonsParams | void>({
            query: (params) => ({
                url: '/v1/teacher/lessons',
                params: params ?? undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    teacher_id: response?.data?.teacher_id ?? response?.data?.TeacherID ?? '',
                    lessons: response?.data?.lessons ?? response?.data?.Lessons ?? [],
                },
            }),
            providesTags: (result) =>
                result?.data?.lessons
                    ? [
                        ...result.data.lessons.map(({ id }) => ({ type: 'Lesson' as const, id })),
                        { type: 'Lesson', id: 'TEACHER-LIST' },
                    ]
                    : [{ type: 'Lesson', id: 'TEACHER-LIST' }],
        }),
        getTeacherLessonAttendance: builder.query<TeacherLessonAttendanceResponse, string>({
            query: (lessonId) => `/v1/teacher/lessons/${lessonId}/attendance`,
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    lesson: response?.data?.lesson,
                    records: response?.data?.records ?? [],
                },
            }),
            providesTags: (_result, _error, lessonId) => [{ type: 'Lesson', id: `TEACHER-ATTENDANCE-${lessonId}` }],
        }),
        submitTeacherLessonAttendance: builder.mutation<{ success: boolean; message: string; data?: { saved_count?: number } }, {
            lessonId: string;
            records: Array<{
                student_id: string;
                status: number;
                note?: string;
            }>;
        }>({
            query: ({ lessonId, records }) => ({
                url: `/v1/teacher/lessons/${lessonId}/attendance`,
                method: 'POST',
                body: { records },
            }),
            invalidatesTags: (_result, _error, { lessonId }) => [{ type: 'Lesson', id: `TEACHER-ATTENDANCE-${lessonId}` }],
        }),
        getTeacherAttendanceSummary: builder.query<TeacherAttendanceSummaryResponse, string>({
            query: (classId) => `/v1/teacher/classes/${classId}/attendance-summary`,
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    teacher_id: response?.data?.teacher_id ?? '',
                    class_id: response?.data?.class_id ?? '',
                    total_lessons: response?.data?.total_lessons ?? 0,
                    students: response?.data?.students ?? [],
                },
            }),
            providesTags: (_result, _error, classId) => [{ type: 'Lesson', id: `TEACHER-ATTENDANCE-SUMMARY-${classId}` }],
        }),
        getTeacherLessonSummary: builder.query<TeacherLessonSummaryResponse, string>({
            query: (lessonId) => `/v1/teacher/lessons/${lessonId}/summary`,
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    lesson: response?.data?.lesson,
                    summary: response?.data?.summary ?? null,
                },
            }),
            providesTags: (_result, _error, lessonId) => [{ type: 'Lesson', id: `TEACHER-SUMMARY-${lessonId}` }],
        }),
        getTeacherLeaveRequests: builder.query<TeacherLeaveRequestsResponse, GetTeacherLeaveRequestsParams | void>({
            query: (params) => ({
                url: '/v1/teacher/leave-requests',
                params: params ?? undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    requests: response?.data?.requests ?? [],
                },
            }),
            providesTags: [{ type: 'LeaveRequest', id: 'TEACHER-LIST' }],
        }),
    }),
});

export const {
    useGetTeacherLessonsQuery,
    useGetTeacherLessonAttendanceQuery,
    useSubmitTeacherLessonAttendanceMutation,
    useGetTeacherAttendanceSummaryQuery,
    useGetTeacherLessonSummaryQuery,
    useGetTeacherLeaveRequestsQuery,
} = teacherPortalApi;
