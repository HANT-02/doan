import { apiClient } from './client';

export interface Teacher {
    id: string;
    code: string;
    full_name: string;
    email: string;
    phone: string;
    is_school_teacher: boolean;
    school_name: string;
    employment_type: 'PART_TIME' | 'FULL_TIME';
    status: 'ACTIVE' | 'INACTIVE';
    notes: string;
    created_at: string;
    updated_at: string;
}

export interface CreateTeacherRequest {
    code?: string;
    full_name: string;
    email?: string;
    phone?: string;
    is_school_teacher?: boolean;
    school_name?: string;
    employment_type?: string;
    status?: string;
    notes?: string;
}

export interface UpdateTeacherRequest {
    code?: string;
    full_name?: string;
    email?: string;
    phone?: string;
    is_school_teacher?: boolean;
    school_name?: string;
    employment_type?: string;
    status?: string;
    notes?: string;
}

export interface ListTeachersParams {
    search?: string;
    status?: string;
    employment_type?: string;
    page?: number;
    limit?: number;
    sort_by?: string;
    sort_order?: 'asc' | 'desc';
}

export interface PaginationMeta {
    items_per_page: number;
    total_items: number;
    current_page: number;
    total_pages: number;
}

export interface ListTeachersResponse {
    teachers: Teacher[];
    pagination: PaginationMeta;
}

export interface ApiResponse<T> {
    success: boolean;
    message: string;
    data: T;
}

export const teacherApi = {
    list: async (params: ListTeachersParams): Promise<ListTeachersResponse> => {
        const response = await apiClient.get('/teachers', { params });
        return response.data as ListTeachersResponse;
    },

    getById: async (id: string): Promise<Teacher> => {
        const response = await apiClient.get(`/teachers/${id}`);
        return response.data as Teacher;
    },

    create: async (data: CreateTeacherRequest): Promise<Teacher> => {
        const response = await apiClient.post('/teachers', data);
        return response.data as Teacher;
    },

    update: async (id: string, data: UpdateTeacherRequest): Promise<Teacher> => {
        const response = await apiClient.put(`/teachers/${id}`, data);
        return response.data as Teacher;
    },

    delete: async (id: string): Promise<void> => {
        await apiClient.delete(`/teachers/${id}`);
    },
};
