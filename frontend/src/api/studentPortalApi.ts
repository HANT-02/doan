import { baseApi } from './baseApi';

export interface StudentTimetableShift {
    id: string;
    code: string;
    name: string;
    start_time: string;
    end_time: string;
    duration_minutes: number;
    session_type: string;
}

export interface StudentTimetableTeacher {
    id?: string;
    code?: string;
    full_name?: string;
}

export interface StudentTimetableLesson {
    id: string;
    class_id: string;
    class_name: string;
    class_code: string;
    teacher: StudentTimetableTeacher;
    room_id?: string;
    room_name?: string;
    date_start: string;
    date_end: string;
    notes: string;
    shift?: StudentTimetableShift | null;
}

export interface StudentTimetableResponse {
    success: boolean;
    message: string;
    data: {
        student_id: string;
        lessons: StudentTimetableLesson[];
    };
}

export interface StudentPortalLeaveRequest {
    id: string;
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
    rejection_reason: string;
    created_at: string;
    updated_at: string;
}

export interface StudentPortalLeaveRequestsResponse {
    success: boolean;
    message: string;
    data: {
        requests: StudentPortalLeaveRequest[];
    };
}

export interface CreateStudentLeaveRequestPayload {
    leave_type: string;
    apply_date: string;
    late_minutes?: number;
    early_minutes?: number;
    reason: string;
    documents?: string[];
    class_id?: string;
    lesson_id?: string;
    subject?: string;
}

export interface GetStudentTimetableParams {
    class_id?: string;
    from?: string;
    to?: string;
}

export interface GetStudentLeaveRequestsParams {
    class_id?: string;
    status?: string;
}

export interface StudentAttendanceSummary {
    total_lessons: number;
    marked_count: number;
    present_count: number;
    absent_count: number;
    late_count: number;
    excused_count: number;
    unmarked_count: number;
    attendance_rate: number;
    absent_rate: number;
    warning: boolean;
    warning_message?: string;
}

export interface StudentAttendanceRecord {
    lesson: StudentTimetableLesson;
    status?: number;
    note: string;
    marked_at?: string;
}

export interface StudentAttendanceResponse {
    success: boolean;
    message: string;
    data: {
        student_id: string;
        class_id?: string;
        summary: StudentAttendanceSummary;
        records: StudentAttendanceRecord[];
    };
}

export interface GetStudentAttendanceParams {
    class_id?: string;
    from?: string;
    to?: string;
}

export interface StudentAcademicRecordLessonSummary {
    id: string;
    topic: string;
    homework: string;
}

export interface StudentAcademicRecordLesson extends StudentTimetableLesson {
    summary: StudentAcademicRecordLessonSummary;
}

export interface StudentAcademicRecordItem {
    record_id: string;
    lesson_summary_id: string;
    lesson: StudentAcademicRecordLesson;
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

export interface StudentAcademicClassSummary {
    class_id: string;
    class_name: string;
    class_code: string;
    records_count: number;
    completed_count: number;
    average_total_score: number;
}

export interface StudentAcademicRecordsResponse {
    success: boolean;
    message: string;
    data: {
        student_id: string;
        class_id?: string;
        class_summaries: StudentAcademicClassSummary[];
        records: StudentAcademicRecordItem[];
    };
}

export interface GetStudentAcademicRecordsParams {
    class_id?: string;
    from?: string;
    to?: string;
}

export interface StudentAtRiskTopFeature {
    key: string;
    label: string;
    value: number;
    display_value: string;
}

export interface StudentAtRiskFeatureSummary {
    attendance_rate_28d: number;
    absence_count_28d: number;
    average_total_score_28d: number;
    homework_completion_rate_28d: number;
    active_enrollment_count_28d: number;
    weekly_lesson_load_28d: number;
    approved_leave_count_28d: number;
    days_since_last_lesson: number;
}

export interface StudentAtRiskPrediction {
    student_id: string;
    student_code: string;
    student_name: string;
    grade_level: string;
    class_id: string;
    class_code: string;
    class_name: string;
    snapshot_at: string;
    risk_label: 'AT_RISK' | 'NOT_AT_RISK';
    risk_score: number;
    risk_band: 'CAO' | 'TRUNG_BINH' | 'THAP';
    model_name: string;
    model_version: string;
    primary_reason: string;
    reasons: string[];
    top_features: StudentAtRiskTopFeature[];
    feature_summary: StudentAtRiskFeatureSummary;
}

export interface StudentAtRiskResponse {
    success: boolean;
    message: string;
    data: {
        student_id: string;
        prediction?: StudentAtRiskPrediction | null;
    };
}

export const studentPortalApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getStudentTimetable: builder.query<StudentTimetableResponse, GetStudentTimetableParams | void>({
            query: (params) => ({
                url: '/v1/student/timetable',
                params: params ?? undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    student_id: response?.data?.student_id ?? '',
                    lessons: response?.data?.lessons ?? [],
                },
            }),
            providesTags: [{ type: 'Lesson', id: 'STUDENT-TIMETABLE' }],
        }),
        getMyStudentLeaveRequests: builder.query<StudentPortalLeaveRequestsResponse, GetStudentLeaveRequestsParams | void>({
            query: (params) => ({
                url: '/v1/student/leave-requests',
                params: params ?? undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    requests: response?.data?.requests ?? [],
                },
            }),
            providesTags: [{ type: 'LeaveRequest', id: 'STUDENT-LIST' }],
        }),
        createStudentLeaveRequest: builder.mutation<{ success: boolean; message: string }, CreateStudentLeaveRequestPayload>({
            query: (body) => ({
                url: '/v1/student/leave-requests',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'LeaveRequest', id: 'STUDENT-LIST' }],
        }),
        cancelStudentLeaveRequest: builder.mutation<{ success: boolean; message: string; data?: { request_id?: string; status?: string } }, string>({
            query: (id) => ({
                url: `/v1/student/leave-requests/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: [{ type: 'LeaveRequest', id: 'STUDENT-LIST' }],
        }),
        getStudentAttendance: builder.query<StudentAttendanceResponse, GetStudentAttendanceParams | void>({
            query: (params) => ({
                url: '/v1/student/attendance',
                params: params ?? undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    student_id: response?.data?.student_id ?? '',
                    class_id: response?.data?.class_id,
                    summary: response?.data?.summary ?? {
                        total_lessons: 0,
                        marked_count: 0,
                        present_count: 0,
                        absent_count: 0,
                        late_count: 0,
                        excused_count: 0,
                        unmarked_count: 0,
                        attendance_rate: 0,
                        absent_rate: 0,
                        warning: false,
                    },
                    records: response?.data?.records ?? [],
                },
            }),
            providesTags: [{ type: 'Lesson', id: 'STUDENT-ATTENDANCE' }],
        }),
        getStudentAcademicRecords: builder.query<StudentAcademicRecordsResponse, GetStudentAcademicRecordsParams | void>({
            query: (params) => ({
                url: '/v1/student/academic-records',
                params: params ?? undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    student_id: response?.data?.student_id ?? '',
                    class_id: response?.data?.class_id,
                    class_summaries: response?.data?.class_summaries ?? [],
                    records: response?.data?.records ?? [],
                },
            }),
            providesTags: [{ type: 'AcademicRecord', id: 'STUDENT-RESULTS' }],
        }),
        getStudentAtRiskPrediction: builder.query<StudentAtRiskResponse, { refresh?: boolean } | void>({
            query: (params) => ({
                url: '/v1/student/at-risk',
                params: params ?? undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    student_id: response?.data?.student_id ?? '',
                    prediction: response?.data?.prediction ?? null,
                },
            }),
            providesTags: [{ type: 'Predictive', id: 'STUDENT-AT-RISK' }],
        }),
    }),
});

export const {
    useGetStudentTimetableQuery,
    useGetMyStudentLeaveRequestsQuery,
    useCreateStudentLeaveRequestMutation,
    useCancelStudentLeaveRequestMutation,
    useGetStudentAttendanceQuery,
    useGetStudentAcademicRecordsQuery,
    useGetStudentAtRiskPredictionQuery,
} = studentPortalApi;
