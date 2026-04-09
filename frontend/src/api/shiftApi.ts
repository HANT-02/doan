import { baseApi } from './baseApi';
import { normalizeDetailEnvelope, normalizeListEnvelope } from './responseUtils';

export type ShiftSessionType = 'MORNING' | 'AFTERNOON' | 'EVENING' | 'CUSTOM';

export interface Shift {
    id: string;
    code: string;
    name: string;
    start_time: string;
    end_time: string;
    duration_minutes: number;
    session_type: ShiftSessionType;
    is_active: boolean;
    notes?: string;
    created_at: string;
    updated_at: string;
}

export interface ListShiftsResponse {
    success: boolean;
    message: string;
    data: {
        shifts: Shift[];
        pagination: {
            items_per_page: number;
            total_items: number;
            current_page: number;
            total_pages: number;
        };
    };
}

export interface ShiftResponse {
    success: boolean;
    message: string;
    data: Shift | null;
}

interface RawShiftListResponse {
    success?: boolean;
    message?: string;
    data?: {
        shifts?: Shift[];
        Shifts?: Shift[];
        pagination?: {
            items_per_page?: number;
            total_items?: number;
            current_page?: number;
            total_pages?: number;
        };
        Pagination?: {
            items_per_page?: number;
            total_items?: number;
            current_page?: number;
            total_pages?: number;
            ItemsPerPage?: number;
            TotalItems?: number;
            CurrentPage?: number;
            TotalPages?: number;
        };
    };
}

interface RawShiftResponse {
    success?: boolean;
    message?: string;
    data?: {
        Shift?: Shift;
    } | Shift;
}

export const shiftApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getShifts: builder.query<ListShiftsResponse, { page?: number; limit?: number; search?: string; session_type?: string; is_active?: boolean; sortBy?: string; sortOrder?: string }>({
            query: (params) => ({
                url: '/v1/shifts',
                params,
            }),
            transformResponse: (response: RawShiftListResponse) => {
                const normalized = normalizeListEnvelope<Shift>(response, ['shifts', 'Shifts']);
                return {
                    success: normalized.success,
                    message: normalized.message,
                    data: {
                        shifts: normalized.items,
                        pagination: normalized.pagination,
                    },
                };
            },
            providesTags: (result) =>
                result?.data?.shifts
                    ? [
                        ...result.data.shifts.map(({ id }) => ({ type: 'Shift' as const, id })),
                        { type: 'Shift', id: 'LIST' },
                    ]
                    : [{ type: 'Shift', id: 'LIST' }],
        }),
        getShiftById: builder.query<ShiftResponse, string>({
            query: (id) => `/v1/shifts/${id}`,
            transformResponse: (response: RawShiftResponse) => normalizeDetailEnvelope<Shift>(response, ['Shift']),
            providesTags: (_result, _error, id) => [{ type: 'Shift', id }],
        }),
        createShift: builder.mutation<ShiftResponse, Partial<Shift>>({
            query: (body) => ({
                url: '/v1/shifts',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'Shift', id: 'LIST' }],
        }),
        updateShift: builder.mutation<ShiftResponse, { id: string; body: Partial<Shift> }>({
            query: ({ id, body }) => ({
                url: `/v1/shifts/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Shift', id }, { type: 'Shift', id: 'LIST' }],
        }),
        deleteShift: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/shifts/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: [{ type: 'Shift', id: 'LIST' }],
        }),
    }),
});

export const {
    useGetShiftsQuery,
    useGetShiftByIdQuery,
    useCreateShiftMutation,
    useUpdateShiftMutation,
    useDeleteShiftMutation,
} = shiftApi;
