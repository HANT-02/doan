import { baseApi } from './baseApi';
import type { LessonAcademicRecord } from './lessonApi';

export interface MyAcademicRecordsResponse {
    success: boolean;
    message: string;
    data: {
        records: LessonAcademicRecord[];
    };
}

export const academicApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getMyAcademicRecords: builder.query<MyAcademicRecordsResponse, { class_id?: string } | void>({
            query: (params) => ({
                url: '/v1/academic-records/my',
                params: params || undefined,
            }),
            transformResponse: (response: any) => ({
                success: response?.success ?? false,
                message: response?.message ?? '',
                data: {
                    records: response?.data?.records || [],
                },
            }),
            providesTags: [{ type: 'AcademicRecord', id: 'MY-LIST' }],
        }),
    }),
});

export const {
    useGetMyAcademicRecordsQuery,
} = academicApi;
