import { baseApi } from './baseApi';
import { normalizeDetailEnvelope } from './responseUtils';

export interface AtRiskPrediction {
    student_id: string;
    student_code: string;
    student_name: string;
    grade_level: string;
    class_id: string;
    class_code: string;
    class_name: string;
    snapshot_at: string;
    label: 'AT_RISK' | 'NOT_AT_RISK';
    risk_score: number;
    risk_band: 'CAO' | 'TRUNG_BINH' | 'THAP';
    model_name: string;
    model_version: string;
    primary_reason: string;
    reasons: string[];
    feature_summary: {
        attendance_rate_28d: number;
        absence_count_28d: number;
        average_total_score_28d: number;
        homework_completion_rate_28d: number;
        active_enrollment_count_28d: number;
        weekly_lesson_load_28d: number;
        approved_leave_count_28d: number;
        days_since_last_lesson: number;
    };
}

export interface PredictiveModelMetadata {
    version: string;
    model_name: string;
    dataset_name: string;
    dataset_source: string;
    definition_name: string;
    definition_version: string;
    trained_at: string;
    train_size: number;
    test_size: number;
    feature_names: string[];
    model_reports: Array<{
        name: string;
        positive_class: string;
        metrics: {
            accuracy: number;
            precision: number;
            recall: number;
            f1_score: number;
            support: number;
        };
        notes?: string[];
    }>;
    recommendation: string[];
}

export interface ListAtRiskPredictionsResponse {
    success: boolean;
    message: string;
    data: {
        items: AtRiskPrediction[];
        pagination: {
            current_page: number;
            items_per_page: number;
            total_items: number;
            total_pages: number;
        };
        summary: {
            total_students_evaluated: number;
            at_risk_count: number;
            not_at_risk_count: number;
        };
        model_metadata: PredictiveModelMetadata;
    };
}

export interface ModelMetadataResponse {
    success: boolean;
    message: string;
    data: {
        model_metadata: PredictiveModelMetadata | null;
    };
}

export interface TrainAtRiskFromDBStep {
    name: string;
    command: string;
    status: 'running' | 'success' | 'failed' | 'warning';
    duration_ms: number;
    output_tail?: string;
}

export interface TrainAtRiskFromDBResponse {
    success: boolean;
    message: string;
    data: {
        message: string;
        dataset_name: string;
        started_at: string;
        finished_at: string;
        duration_ms: number;
        ml_dir: string;
        model_metadata: PredictiveModelMetadata;
        steps: TrainAtRiskFromDBStep[];
    };
}

export const predictiveApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getAtRiskPredictions: builder.query<
            ListAtRiskPredictionsResponse,
            { page?: number; limit?: number; search?: string; only_at_risk?: boolean; refresh?: boolean }
        >({
            query: (params) => ({
                url: '/v1/predictive/at-risk/students',
                params,
            }),
            providesTags: [{ type: 'Predictive', id: 'AT_RISK_LIST' }],
        }),
        getPredictiveModelMetadata: builder.query<ModelMetadataResponse, { refresh?: boolean } | void>({
            query: (params) => params
                ? {
                    url: '/v1/predictive/at-risk/model-metadata',
                    params,
                }
                : '/v1/predictive/at-risk/model-metadata',
            transformResponse: (response: any) => {
                const normalized = normalizeDetailEnvelope<PredictiveModelMetadata>(
                    response,
                    ['model_metadata'],
                );

                return {
                    success: normalized.success,
                    message: normalized.message,
                    data: {
                        model_metadata: normalized.data ?? response.data?.model_metadata ?? null,
                    },
                };
            },
            providesTags: [{ type: 'Predictive', id: 'MODEL_METADATA' }],
        }),
        trainAtRiskFromDB: builder.mutation<TrainAtRiskFromDBResponse, void>({
            query: () => ({
                url: '/v1/predictive/at-risk/train-from-db',
                method: 'POST',
            }),
            invalidatesTags: [
                { type: 'Predictive', id: 'AT_RISK_LIST' },
                { type: 'Predictive', id: 'MODEL_METADATA' },
            ],
        }),
    }),
});

export const {
    useGetAtRiskPredictionsQuery,
    useGetPredictiveModelMetadataQuery,
    useTrainAtRiskFromDBMutation,
} = predictiveApi;
