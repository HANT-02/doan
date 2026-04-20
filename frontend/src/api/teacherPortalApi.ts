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

export interface UpsertTeacherLessonSummaryPayload {
    topic: string;
    lesson_content: string;
    class_feedback: string;
    homework: string;
    homework_deadline?: string | null;
    teacher_notes: string;
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

export interface TeacherLeaveRequestStatusResponse {
    success: boolean;
    message: string;
    data: {
        request_id: string;
        status: string;
    };
}

export interface TeacherLessonAcademicRecord {
    record_id?: string;
    lesson_summary_id?: string;
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
    created_at?: string;
    updated_at?: string;
}

export interface TeacherLessonAcademicRecordsResponse {
    success: boolean;
    message: string;
    data: {
        lesson: TeacherLesson;
        records: TeacherLessonAcademicRecord[];
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
        upsertTeacherLessonSummary: builder.mutation<{ success: boolean; message: string; data: TeacherLessonSummary | null }, {
            lessonId: string;
            body: UpsertTeacherLessonSummaryPayload;
        }>({
            query: ({ lessonId, body }) => ({
                url: `/v1/teacher/lessons/${lessonId}/summary`,
                method: 'PUT',
                body,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: response?.data ?? null,
            }),
            invalidatesTags: (_result, _error, { lessonId }) => [{ type: 'Lesson', id: `TEACHER-SUMMARY-${lessonId}` }],
        }),
        getTeacherLessonAcademicRecords: builder.query<TeacherLessonAcademicRecordsResponse, string>({
            query: (lessonId) => `/v1/teacher/lessons/${lessonId}/records`,
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    lesson: response?.data?.lesson,
                    records: response?.data?.records ?? [],
                },
            }),
            providesTags: (_result, _error, lessonId) => [{ type: 'AcademicRecord', id: `TEACHER-LESSON-${lessonId}` }],
        }),
        upsertTeacherLessonAcademicRecord: builder.mutation<{ success: boolean; message: string; data?: { saved_count?: number } }, {
            lessonId: string;
            studentId: string;
            body: {
                homework_completed: boolean;
                homework_score: number;
                attitude_rating: number;
                participation_score: number;
                personal_comment: string;
            };
        }>({
            query: ({ lessonId, studentId, body }) => ({
                url: `/v1/teacher/lessons/${lessonId}/records/${studentId}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { lessonId }) => [{ type: 'AcademicRecord', id: `TEACHER-LESSON-${lessonId}` }],
        }),
        finalizeTeacherLessonAcademicRecords: builder.mutation<{ success: boolean; message: string; data?: { finalized_count?: number } }, string>({
            query: (lessonId) => ({
                url: `/v1/teacher/lessons/${lessonId}/records/finalize`,
                method: 'POST',
            }),
            invalidatesTags: (_result, _error, lessonId) => [{ type: 'AcademicRecord', id: `TEACHER-LESSON-${lessonId}` }],
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
        approveTeacherLeaveRequest: builder.mutation<TeacherLeaveRequestStatusResponse, string>({
            query: (id) => ({
                url: `/v1/teacher/leave-requests/${id}/approve`,
                method: 'POST',
            }),
            invalidatesTags: [{ type: 'LeaveRequest', id: 'TEACHER-LIST' }],
        }),
        rejectTeacherLeaveRequest: builder.mutation<TeacherLeaveRequestStatusResponse, { id: string; rejection_reason: string }>({
            query: ({ id, rejection_reason }) => ({
                url: `/v1/teacher/leave-requests/${id}/reject`,
                method: 'POST',
                body: { rejection_reason },
            }),
            invalidatesTags: [{ type: 'LeaveRequest', id: 'TEACHER-LIST' }],
        }),
    }),
});

export const {
    useGetTeacherLessonsQuery,
    useGetTeacherLessonAttendanceQuery,
    useSubmitTeacherLessonAttendanceMutation,
    useGetTeacherAttendanceSummaryQuery,
    useGetTeacherLessonSummaryQuery,
    useUpsertTeacherLessonSummaryMutation,
    useGetTeacherLessonAcademicRecordsQuery,
    useUpsertTeacherLessonAcademicRecordMutation,
    useFinalizeTeacherLessonAcademicRecordsMutation,
    useGetTeacherLeaveRequestsQuery,
    useApproveTeacherLeaveRequestMutation,
    useRejectTeacherLeaveRequestMutation,
} = teacherPortalApi;
