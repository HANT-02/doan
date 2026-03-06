import { baseApi } from './baseApi';

export interface MaterialLabel {
    id: string;
    code: string;
    name: string;
    severity: 'SAFE' | 'WARNING' | 'DANGER';
    description: string;
}

export interface MaterialAuditLog {
    id: string;
    status: string;
    provider: string;
    raw_ocr_text: string;
    confidence_score: number;
    reasoning: string;
    detected_issues: string[];
    triggered_at: string;
    completed_at?: string;
    label?: MaterialLabel;
}

export interface MaterialDecision {
    id: string;
    compliance_officer_id: string;
    approved: boolean;
    reject_reason: string;
    notes: string;
    decided_at: string;
}

export interface MaterialItem {
    id: string;
    teacher_id: string;
    title: string;
    description: string;
    file_name: string;
    file_path: string;
    file_type: string;
    status: 'UPLOADED' | 'SCANNING' | 'AI_REVIEWED' | 'APPROVED' | 'REJECTED';
    uploaded_at: string;
    latest_label?: MaterialLabel;
    latest_audit?: MaterialAuditLog;
    last_decision?: MaterialDecision;
    audit_logs: MaterialAuditLog[];
}

interface RawMaterialListResponse {
    success?: boolean;
    message?: string;
    data?: {
        materials?: MaterialItem[];
        Materials?: MaterialItem[];
    };
}

interface RawMaterialResponse {
    success?: boolean;
    message?: string;
    data?: MaterialItem;
}

interface MaterialListResponse {
    success: boolean;
    message?: string;
    data: {
        materials: MaterialItem[];
    };
}

interface MaterialResponse {
    success: boolean;
    message?: string;
    data: MaterialItem;
}

const normalizeMaterialResponse = (response: RawMaterialResponse): MaterialResponse => ({
    success: !!response.success,
    message: response.message,
    data: response.data as MaterialItem,
});

const normalizeMaterialListResponse = (response: RawMaterialListResponse): MaterialListResponse => ({
    success: !!response.success,
    message: response.message,
    data: {
        materials: response.data?.materials || response.data?.Materials || [],
    },
});

export const materialApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getMaterials: builder.query<MaterialListResponse, { teacher_id?: string; status?: string; queue?: string }>({
            query: (params) => ({
                url: '/v1/materials',
                params,
            }),
            transformResponse: normalizeMaterialListResponse,
            providesTags: (result) =>
                result?.data?.materials?.length
                    ? [
                        ...result.data.materials.map(({ id }) => ({ type: 'Material' as const, id })),
                        { type: 'Material', id: 'LIST' },
                    ]
                    : [{ type: 'Material', id: 'LIST' }],
        }),
        getMaterialById: builder.query<MaterialResponse, string>({
            query: (id) => `/v1/materials/${id}`,
            transformResponse: normalizeMaterialResponse,
            providesTags: (_result, _error, id) => [{ type: 'Material', id }],
        }),
        uploadMaterial: builder.mutation<MaterialResponse, FormData>({
            query: (body) => ({
                url: '/v1/materials/upload',
                method: 'POST',
                body,
            }),
            transformResponse: normalizeMaterialResponse,
            invalidatesTags: [{ type: 'Material', id: 'LIST' }],
        }),
        reviewMaterial: builder.mutation<MaterialResponse, { id: string; compliance_officer_id: string; approved: boolean; reject_reason?: string; notes?: string }>({
            query: ({ id, ...body }) => ({
                url: `/v1/materials/${id}/review`,
                method: 'POST',
                body,
            }),
            transformResponse: normalizeMaterialResponse,
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Material', id }, { type: 'Material', id: 'LIST' }],
        }),
    }),
});

export const {
    useGetMaterialsQuery,
    useGetMaterialByIdQuery,
    useUploadMaterialMutation,
    useReviewMaterialMutation,
} = materialApi;
