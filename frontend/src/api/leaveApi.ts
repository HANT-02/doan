import { baseApi } from './baseApi';

export interface LeaveRequest {
    id: string;
    leave_type: string;
    apply_date: string;
    late_minutes: number;
    early_minutes: number;
    reason: string;
    documents: string[];
    class_id?: string;
    class?: {
        id: string;
        code: string;
        name: string;
    };
    lesson_id?: string;
    lesson?: {
        id: string;
        class_id: string;
        date_start: string;
        date_end: string;
    };
    subject: string;
    status: string;
    approved_by_id?: string;
    approved_at?: string;
    rejection_reason: string;
    student: {
        id: string;
        code: string;
        full_name: string;
        email?: string;
    };
    created_at: string;
    updated_at: string;
}

export interface LeaveRequestsResponse {
    success: boolean;
    message: string;
    data: {
        requests: LeaveRequest[];
    };
}

export const leaveApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getLeaveRequests: builder.query<LeaveRequestsResponse, { status?: string; class_id?: string } | void>({
            query: (params) => ({
                url: '/v1/leave-requests',
                params: params || undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    requests: response?.data?.requests || [],
                },
            }),
            providesTags: [{ type: 'LeaveRequest', id: 'LIST' }],
        }),
        createLeaveRequest: builder.mutation<{ success: boolean; message: string }, Record<string, unknown>>({
            query: (body) => ({
                url: '/v1/leave-requests',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'LeaveRequest', id: 'LIST' }],
        }),
        approveLeaveRequest: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/leave-requests/${id}/approve`,
                method: 'POST',
            }),
            invalidatesTags: [{ type: 'LeaveRequest', id: 'LIST' }],
        }),
        rejectLeaveRequest: builder.mutation<{ success: boolean; message: string }, { id: string; rejection_reason: string }>({
            query: ({ id, rejection_reason }) => ({
                url: `/v1/leave-requests/${id}/reject`,
                method: 'POST',
                body: { rejection_reason },
            }),
            invalidatesTags: [{ type: 'LeaveRequest', id: 'LIST' }],
        }),
        cancelLeaveRequest: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/leave-requests/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: [{ type: 'LeaveRequest', id: 'LIST' }],
        }),
    }),
});

export const {
    useGetLeaveRequestsQuery,
    useCreateLeaveRequestMutation,
    useApproveLeaveRequestMutation,
    useRejectLeaveRequestMutation,
    useCancelLeaveRequestMutation,
} = leaveApi;
