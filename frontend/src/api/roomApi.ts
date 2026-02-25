import { baseApi } from './baseApi';

export interface Room {
    id: string;
    name: string;
    capacity: number;
    location?: string;
    status: 'ACTIVE' | 'MAINTENANCE' | 'INACTIVE';
    created_at: string;
    updated_at: string;
}

export interface ListRoomsResponse {
    success: boolean;
    message: string;
    data: {
        rooms: Room[];
        pagination: {
            items_per_page: number;
            total_items: number;
            current_page: number;
            total_pages: number;
        };
    };
}

export interface RoomResponse {
    success: boolean;
    message: string;
    data: Room;
}

export const roomApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getRooms: builder.query<ListRoomsResponse, { page?: number; limit?: number; search?: string; sortBy?: string; sortOrder?: string }>({
            query: (params) => ({
                url: '/v1/rooms',
                params,
            }),
            transformResponse: (response: any) => {
                return {
                    success: response.success,
                    message: response.message,
                    data: {
                        rooms: response.data?.Rooms || [],
                        pagination: {
                            items_per_page: response.data?.Pagination?.ItemsPerPage || 10,
                            total_items: response.data?.Pagination?.TotalItems || 0,
                            current_page: response.data?.Pagination?.CurrentPage || 1,
                            total_pages: response.data?.Pagination?.TotalPages || 1,
                        }
                    }
                };
            },
            providesTags: (result) =>
                result?.data?.rooms
                    ? [
                        ...result.data.rooms.map(({ id }) => ({ type: 'Room' as const, id })),
                        { type: 'Room', id: 'LIST' },
                    ]
                    : [{ type: 'Room', id: 'LIST' }],
        }),
        getRoomById: builder.query<RoomResponse, string>({
            query: (id) => `/v1/rooms/${id}`,
            transformResponse: (response: any) => ({
                ...response,
                data: response.data?.Room || response.data
            }),
            providesTags: (_result, _error, id) => [{ type: 'Room', id }],
        }),
        createRoom: builder.mutation<RoomResponse, Partial<Room>>({
            query: (body) => ({
                url: '/v1/rooms',
                method: 'POST',
                body,
            }),
            invalidatesTags: [{ type: 'Room', id: 'LIST' }],
        }),
        updateRoom: builder.mutation<RoomResponse, { id: string; body: Partial<Room> }>({
            query: ({ id, body }) => ({
                url: `/v1/rooms/${id}`,
                method: 'PUT',
                body,
            }),
            invalidatesTags: (_result, _error, { id }) => [{ type: 'Room', id }, { type: 'Room', id: 'LIST' }],
        }),
        deleteRoom: builder.mutation<{ success: boolean; message: string }, string>({
            query: (id) => ({
                url: `/v1/rooms/${id}`,
                method: 'DELETE',
            }),
            invalidatesTags: [{ type: 'Room', id: 'LIST' }],
        }),
    }),
});

export const {
    useGetRoomsQuery,
    useGetRoomByIdQuery,
    useCreateRoomMutation,
    useUpdateRoomMutation,
    useDeleteRoomMutation,
} = roomApi;
