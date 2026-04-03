import { baseApi } from './baseApi';

export interface SchedulingFilters {
    date_from: string;
    date_to: string;
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

export interface SchedulingPreview {
    run_id: string;
    status: 'FAILED' | 'PARTIAL' | 'COMPLETED';
    generated_at: string;
    filters: SchedulingFilters;
    summary: {
        requested_classes: number;
        requested_sessions: number;
        scheduled_lessons: number;
        unscheduled_lessons: number;
        conflict_count: number;
        soft_score: number;
    };
    assignments: SchedulingAssignment[];
    conflicts: SchedulingConflict[];
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
        status: 'FAILED',
        generated_at: '',
        filters: {
            date_from: '',
            date_to: '',
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
        },
        assignments: [],
        conflicts: [],
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
        commitSchedulingPreview: builder.mutation<CommitSchedulingResponse, { run_id: string }>({
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
        }),
    }),
});

export const {
    usePreviewSchedulingMutation,
    useLazyGetSchedulingPreviewQuery,
    useLazyGetLatestSchedulingPreviewQuery,
    useCommitSchedulingPreviewMutation,
} = schedulingApi;
