import { baseApi } from './baseApi';

export interface SchedulingFilters {
    date_from?: string;
    date_to?: string;
    effective_date_from?: string;
    mode?: 'cold_start' | 'replan_draft' | 'replan_with_published_lock';
    class_ids?: string[];
    teacher_ids?: string[];
    room_ids?: string[];
}

export interface SchedulingAssignment {
    variable_id: string;
    class_id: string;
    class_code: string;
    class_name: string;
    session_index: number;
    session_total: number;
    teacher_id: string;
    teacher_label: string;
    room_id: string;
    room_name: string;
    room_capacity: number;
    replace_lesson_id?: string;
    shift_id?: string;
    shift_code?: string;
    shift_name?: string;
    shift_type?: string;
    start_time: string;
    end_time: string;
    constraint_fit: string;
}

export interface SchedulingConflict {
    variable_id: string;
    class_id: string;
    class_code: string;
    class_name: string;
    session_index: number;
    session_total: number;
    type: string;
    message: string;
}

export interface SchedulingCandidateOption {
    key: string;
    room_id: string;
    room_name: string;
    room_capacity: number;
    shift_id?: string;
    shift_code?: string;
    shift_name?: string;
    shift_type?: string;
    start_time: string;
    end_time: string;
}

export interface SchedulingExistingLesson {
    lesson_id: string;
    class_id: string;
    class_code: string;
    class_name: string;
    status?: string;
    teacher_id?: string;
    teacher_label?: string;
    room_id?: string;
    room_name?: string;
    start_time: string;
    end_time: string;
    notes?: string;
    student_ids?: string[];
}

export interface SchedulingPreview {
    run_id: string;
    mode?: string;
    status: 'FAILED' | 'PARTIAL' | 'COMPLETED';
    generated_at: string;
    effective_date_from?: string;
    filters: SchedulingFilters;
    summary: {
        requested_classes: number;
        requested_sessions: number;
        scheduled_lessons: number;
        unscheduled_lessons: number;
        conflict_count: number;
        soft_score: number;
        schedule_change_count: number;
        teacher_change_count: number;
        room_change_count: number;
        average_capacity_utilization: number;
    };
    assignments: SchedulingAssignment[];
    conflicts: SchedulingConflict[];
    existing_lessons: SchedulingExistingLesson[];
    candidate_options: Record<string, SchedulingCandidateOption[]>;
}

export interface SchedulingPreviewResponse {
    success: boolean;
    message?: string;
    data: SchedulingPreview;
}

export interface CommitSchedulingResponse {
    success: boolean;
    message?: string;
    data: {
        message: string;
        scheduled_lessons: number;
        status: string;
    };
}

export interface SubstituteSuggestion {
    teacher_id: string;
    teacher_name: string;
    teacher_code: string;
    score: number;
    match_reasons: string[];
    is_available: boolean;
}

export interface SuggestSubstituteResponse {
    success: boolean;
    message?: string;
    data: SubstituteSuggestion[];
}

export interface AssignSubstituteRequest {
    new_teacher_id: string;
    reason: string;
}

export interface AssignSubstituteResponse {
    success: boolean;
    message?: string;
}

export interface MakeupSpot {
    lesson_id: string;
    class_id: string;
    class_code: string;
    class_name: string;
    course_id?: string;
    course_name?: string;
    teacher_id?: string;
    teacher_name?: string;
    room_id?: string;
    room_name?: string;
    start_time: string;
    end_time: string;
    match_type: string;
    expected_student_count: number;
    capacity_limit: number;
    remaining_capacity: number;
    capacity_utilization: number;
    eligible: boolean;
    reasons: string[];
}

export interface FindMakeupSpotsResponse {
    success: boolean;
    message?: string;
    data: {
        student_id: string;
        lesson_id: string;
        spots: MakeupSpot[];
    };
}

type RawSchedulingPreviewResponse = {
    success?: boolean;
    message?: string;
    data?: SchedulingPreview;
};

type RawCommitSchedulingResponse = {
    success?: boolean;
    message?: string;
    data?: {
        message?: string;
        scheduled_lessons?: number;
        status?: string;
    };
};

const normalizePreviewResponse = (response: RawSchedulingPreviewResponse): SchedulingPreviewResponse => ({
    success: !!response.success,
    message: response.message,
    data: response.data || {
        run_id: '',
        mode: 'replan_with_published_lock',
        status: 'FAILED',
        generated_at: '',
        effective_date_from: '',
        filters: {
            date_from: '',
            date_to: '',
            effective_date_from: '',
            mode: 'replan_with_published_lock',
            class_ids: [],
            teacher_ids: [],
            room_ids: [],
        },
        summary: {
            requested_classes: 0,
            requested_sessions: 0,
            scheduled_lessons: 0,
            unscheduled_lessons: 0,
            conflict_count: 0,
            soft_score: 0,
            schedule_change_count: 0,
            teacher_change_count: 0,
            room_change_count: 0,
            average_capacity_utilization: 0,
        },
        assignments: [],
        conflicts: [],
        existing_lessons: [],
        candidate_options: {},
    },
});

export const schedulingApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        previewScheduling: builder.mutation<SchedulingPreviewResponse, SchedulingFilters>({
            query: (body) => ({
                url: '/v1/scheduling/preview',
                method: 'POST',
                body,
            }),
            transformResponse: normalizePreviewResponse,
            invalidatesTags: [{ type: 'Scheduling', id: 'LATEST' }],
        }),
        getSchedulingPreview: builder.query<SchedulingPreviewResponse, string>({
            query: (id) => `/v1/scheduling/preview/${id}`,
            transformResponse: normalizePreviewResponse,
            providesTags: (_result, _error, id) => [{ type: 'Scheduling', id }],
        }),
        getLatestSchedulingPreview: builder.query<SchedulingPreviewResponse, void>({
            query: () => '/v1/scheduling/preview/latest',
            transformResponse: normalizePreviewResponse,
            providesTags: [{ type: 'Scheduling', id: 'LATEST' }],
        }),
        commitSchedulingPreview: builder.mutation<CommitSchedulingResponse, {
            run_id: string;
            manual_assignments?: Array<{
                variable_id: string;
                option_key: string;
            }>;
        }>({
            query: (body) => ({
                url: '/v1/scheduling/commit',
                method: 'POST',
                body,
            }),
            transformResponse: (response: RawCommitSchedulingResponse) => ({
                success: !!response.success,
                message: response.message,
                data: {
                    message: response.data?.message || '',
                    scheduled_lessons: response.data?.scheduled_lessons || 0,
                    status: response.data?.status || 'FAILED',
                },
            }),
            invalidatesTags: [
                { type: 'Lesson', id: 'LIST' },
                { type: 'Lesson', id: 'TEACHER-LIST' },
                { type: 'Lesson', id: 'STUDENT-TIMETABLE' },
                { type: 'Scheduling', id: 'LATEST' },
            ],
        }),
        suggestSubstitute: builder.query<SuggestSubstituteResponse, string>({
            query: (lessonId) => `/v1/scheduling/lessons/${lessonId}/suggest-substitutes`,
            providesTags: ['Scheduling'],
        }),
        assignSubstitute: builder.mutation<AssignSubstituteResponse, { lessonId: string, data: AssignSubstituteRequest }>({
            query: ({ lessonId, data }) => ({
                url: `/v1/scheduling/lessons/${lessonId}/assign-substitute`,
                method: 'POST',
                body: data,
            }),
            invalidatesTags: [
                { type: 'Lesson', id: 'LIST' },
                { type: 'Lesson', id: 'TEACHER-LIST' },
            ],
        }),
        findMakeupSpots: builder.query<FindMakeupSpotsResponse, { lessonId: string; studentId: string; limit?: number }>({
            query: ({ lessonId, studentId, limit }) => ({
                url: `/v1/scheduling/lessons/${lessonId}/find-makeup-spots`,
                params: {
                    student_id: studentId,
                    limit,
                },
            }),
            providesTags: ['Scheduling'],
        }),
    }),
});

export const {
    usePreviewSchedulingMutation,
    useLazyGetSchedulingPreviewQuery,
    useLazyGetLatestSchedulingPreviewQuery,
    useCommitSchedulingPreviewMutation,
    useLazySuggestSubstituteQuery,
    useAssignSubstituteMutation,
    useLazyFindMakeupSpotsQuery,
} = schedulingApi;
