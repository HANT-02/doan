import { baseApi } from './baseApi';
import { normalizeDetailEnvelope, normalizeListEnvelope } from './responseUtils';
import type { Student } from './studentApi';

export interface Class {
    id: string;
    code: string;
    name: string;
    notes?: string;
    start_date: string;
    end_date?: string;
    max_students: number;
    status: 'OPEN' | 'CLOSED' | 'CANCELLED';
    price: number;
    program_id?: string;
    course_id?: string;
    teacher_id?: string;
    created_at: string;
    updated_at: string;
}

export interface ClassSchedule {
    id: string;
    class_id: string;
    shift_id: string;
    room_id?: string;
    day_of_week: string;
    shift?: {
        id: string;
        name: string;
        start_time: string;
        end_time: string;
    };
    room?: {
        id: string;
        name: string;
        code: string;
    };
}


export interface ListClassesResponse {
    success: boolean;
    message: string;
    data: {
        classes: Class[];
        pagination: {
            items_per_page: number;
            total_items: number;
            current_page: number;
            total_pages: number;
        };
    };
}

export interface ClassResponse {
    success: boolean;
    message: string;
    data: Class;
}

export interface ClassRoster {
    class_id: string;
    max_students: number;
    capacity_limit: number;
    current_count: number;
    students: Student[];
}

export interface ClassRosterResponse {
    success: boolean;
    message: string;
    data: ClassRoster;
}

export interface ClassScheduleListResponse {
    success: boolean;
    message: string;
    data: {
        schedules: ClassSchedule[];
    };
}


interface RawClassListResponse {
    success?: boolean;
    message?: string;
    data?: {
        classes?: Class[];
        Classes?: Class[];
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

interface RawClassResponse {
    success?: boolean;
    message?: string;
    data?: {
        Class?: Class;
    } | Class;
}

interface RawClassRosterResponse {
    success?: boolean;
    message?: string;
    data?: {
        class_id?: string;
        max_students?: number;
        capacity_limit?: number;
        current_count?: number;
        students?: Student[];
        ClassID?: string;
        MaxStudents?: number;
        CapacityLimit?: number;
        CurrentCount?: number;
        Students?: Student[];
    };
}

export const classApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getClasses: builder.query<ListClassesResponse, { page?: number; limit?: number; search?: string; status?: string; sortBy?: string; sortOrder?: string }>({
            query: (params) => ({
                url: '/v1/classes',
                params,
            }),
            transformResponse: (response: RawClassListResponse) => {
                const normalized = normalizeListEnvelope<Class>(response, ['classes', 'Classes']);
                return {
                    success: normalized.success,
                    message: normalized.message,
                    data: {
                        classes: normalized.items,
                        pagination: normalized.pagination,
                    },
                };
            },
            providesTags: (result) =>
                result?.data?.classes
                    ? [
                        ...result.data.classes.map(({ id }) => ({ type: 'Class' as const, id })),
                        { type: 'Class', id: 'LIST' },
                    ]
                    : [{ type: 'Class', id: 'LIST' }],
        }),
        getClassById: builder.query<ClassResponse, string>({
            query: (id) => `/v1/classes/${id}`,
            transformResponse: (response: RawClassResponse) => normalizeDetailEnvelope<Class>(response, ['Class']) as ClassResponse,
            providesTags: (_result, _error, id) => [{ type: 'Class', id }],
        }),
        getClassRoster: builder.query<ClassRosterResponse, string>({
            query: (id) => `/v1/classes/${id}/students`,
            transformResponse: (response: RawClassRosterResponse) => ({
                success: response.success ?? false,
                message: response.message ?? '',
                data: {
                    class_id: response.data?.class_id || response.data?.ClassID || '',
                    max_students: response.data?.max_students || response.data?.MaxStudents || 0,
                    capacity_limit: response.data?.capacity_limit || response.data?.CapacityLimit || 0,
                    current_count: response.data?.current_count || response.data?.CurrentCount || 0,
                    students: response.data?.students || response.data?.Students || [],
                },
            }),
            providesTags: (_result, _error, id) => [{ type: 'Class', id }, { type: 'Class', id: `ROSTER-${id}` }],
        }),
        createClass: builder.mutation<ClassResponse, Partial<Class>>({
            query: (body) => ({
                url: '/v1/classes',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'Class', id: 'LIST' }],
        }),
        updateClass: builder.mutation<ClassResponse, { id: string; body: Partial<Class> }>({
            query: ({ id, body }) => ({
                url: `/v1/classes/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Class', id }, { type: 'Class', id: 'LIST' }],
        }),
        deleteClass: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/classes/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: [{ type: 'Class', id: 'LIST' }],
        }),
        enrollStudents: builder.mutation<{ success: boolean; message: string; data?: { enrolled?: number } }, { id: string; student_ids: string[] }>({
            query: ({ id, student_ids }) => ({
                url: `/v1/classes/${id}/students`,
                method: 'POST',
                body: { student_ids },
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Class', id }, { type: 'Class', id: `ROSTER-${id}` }, { type: 'Student', id: 'LIST' }],
        }),
        removeStudents: builder.mutation<{ success: boolean; message: string }, { id: string; student_ids: string[] }>({
            query: ({ id, student_ids }) => ({
                url: `/v1/classes/${id}/students`,
                method: 'DELETE',
                body: { student_ids },
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Class', id }, { type: 'Class', id: `ROSTER-${id}` }, { type: 'Student', id: 'LIST' }],
        }),
        assignTeacher: builder.mutation<{ success: boolean; message: string }, { id: string; teacher_id: string }>({
            query: ({ id, teacher_id }) => ({
                url: `/v1/classes/${id}/teacher`,
                method: 'PUT',
                body: { teacher_id },
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Class', id }, { type: 'Class', id: 'LIST' }, { type: 'Teacher', id: 'LIST' }],
        }),
        getClassSchedules: builder.query<ClassScheduleListResponse, string>({
            query: (id) => `/v1/classes/${id}/schedules`,
            providesTags: (_result, _error, id) => [{ type: 'Class', id: `SCHEDULES-${id}` }],
        }),
        createClassSchedule: builder.mutation<{ success: boolean; message: string; data: { schedule: ClassSchedule } }, { classId: string; shift_id: string; day_of_week: string; room_id?: string }>({
            query: ({ classId, ...body }) => ({
                url: `/v1/classes/${classId}/schedules`,
                method: 'POST',
                body,
            }),
            invalidatesTags: (_result, _error, { classId }) => [{ type: 'Class', id: `SCHEDULES-${classId}` }],
        }),
        deleteClassSchedule: builder.mutation<{ success: boolean; message: string }, { classId: string; scheduleId: string }>({
            query: ({ classId, scheduleId }) => ({
                url: `/v1/classes/${classId}/schedules/${scheduleId}`,
                method: 'DELETE',
            }),
            invalidatesTags: (_result, _error, { classId }) => [{ type: 'Class', id: `SCHEDULES-${classId}` }],
        }),
    }),
});

export const {
    useGetClassesQuery,
    useGetClassByIdQuery,
    useGetClassRosterQuery,
    useCreateClassMutation,
    useUpdateClassMutation,
    useDeleteClassMutation,
    useEnrollStudentsMutation,
    useRemoveStudentsMutation,
    useAssignTeacherMutation,
    useGetClassSchedulesQuery,
    useCreateClassScheduleMutation,
    useDeleteClassScheduleMutation,
} = classApi;
