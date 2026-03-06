import { baseApi } from './baseApi';
import type { Course } from './courseApi';
import { normalizeDetailEnvelope, normalizeListEnvelope } from './responseUtils';

export interface Program {
    id: string;
    code: string;
    name: string;
    track: string;
    effective_from?: string;
    effective_to?: string;
    approval_note?: string;
    created_at: string;
    updated_at: string;
    courses?: Course[];
}

export interface ListProgramsParams {
    page?: number;
    limit?: number;
    search?: string;
    track?: string;
}

export interface ListProgramsResponse {
    success: boolean;
    message: string;
    data: {
        programs: Program[];
        pagination: {
            total_pages: number;
            total_items: number;
            current_page: number;
            items_per_page: number;
        };
    };
}

export interface ProgramResponse {
    success: boolean;
    message: string;
    data: Program | null;
}

interface RawProgramListResponse {
    success?: boolean;
    message?: string;
    data?: {
        programs?: Program[];
        Programs?: Program[];
        pagination?: {
            total_pages?: number;
            total_items?: number;
            current_page?: number;
            items_per_page?: number;
            limit?: number;
        };
        Pagination?: {
            TotalPages?: number;
            TotalItems?: number;
            CurrentPage?: number;
            ItemsPerPage?: number;
        };
    };
}

interface RawProgramResponse {
    success?: boolean;
    message?: string;
    data?: {
        Program?: Program;
    } | Program;
}

export const programApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getPrograms: builder.query<ListProgramsResponse, ListProgramsParams | void>({
            query: (params) => ({
                url: '/v1/programs',
                params: params || undefined,
            }),
            transformResponse: (response: RawProgramListResponse) => {
                const normalized = normalizeListEnvelope<Program>(response, ['programs', 'Programs']);
                return {
                    success: normalized.success,
                    message: normalized.message,
                    data: {
                        programs: normalized.items,
                        pagination: normalized.pagination,
                    },
                };
            },
            providesTags: (result) =>
                result?.data?.programs
                    ? [
                        ...result.data.programs.map(({ id }) => ({ type: 'Program' as const, id })),
                        { type: 'Program', id: 'LIST' },
                    ]
                    : [{ type: 'Program', id: 'LIST' }],
        }),
        getProgramById: builder.query<ProgramResponse, string>({
            query: (id) => `/v1/programs/${id}`,
            transformResponse: (response: RawProgramResponse) => normalizeDetailEnvelope<Program>(response, ['Program']),
            providesTags: (_result, _error, id) => [{ type: 'Program', id }],
        }),
        createProgram: builder.mutation<ProgramResponse, Partial<Program>>({
            query: (body) => ({
                url: '/v1/programs',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'Program', id: 'LIST' }],
        }),
        updateProgram: builder.mutation<ProgramResponse, Partial<Program> & { id: string }>({
            query: ({ id, ...body }) => ({
                url: `/v1/programs/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Program', id }, { type: 'Program', id: 'LIST' }],
        }),
        deleteProgram: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/programs/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: [{ type: 'Program', id: 'LIST' }],
        }),
        addCoursesToProgram: builder.mutation<{ success: boolean; message: string }, { programId: string, courseIds: string[] }>({
            query: ({ programId, courseIds }) => ({
                url: `/v1/programs/${programId}/courses`,
                method: 'POST',
                body: { course_ids: courseIds },
            }),
            invalidatesTags: (_result, _error, { programId }) => [{ type: 'Program', id: programId }, { type: 'Program', id: 'LIST' }, 'Course'],
        }),
        removeCoursesFromProgram: builder.mutation<{ success: boolean; message: string }, { programId: string, courseIds: string[] }>({
            query: ({ programId, courseIds }) => ({
                url: `/v1/programs/${programId}/courses`,
                method: 'DELETE',
                body: { course_ids: courseIds },
            }),
            invalidatesTags: (_result, _error, { programId }) => [{ type: 'Program', id: programId }, { type: 'Program', id: 'LIST' }, 'Course'],
        }),
    }),
});

export const {
    useGetProgramsQuery,
    useGetProgramByIdQuery,
    useCreateProgramMutation,
    useUpdateProgramMutation,
    useDeleteProgramMutation,
    useAddCoursesToProgramMutation,
    useRemoveCoursesFromProgramMutation,
} = programApi;
