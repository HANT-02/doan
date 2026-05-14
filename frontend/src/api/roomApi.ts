import { baseApi } from './baseApi';
import { normalizeDetailEnvelope, normalizeListEnvelope } from './responseUtils';

export interface Room {
    id: string;
    code?: string;
    name: string;
    capacity: number;
    address?: string;
    campus_id?: string;
    campus?: {
        id: string;
        code: string;
        name: string;
    };
    location?: string;
    status?: 'ACTIVE' | 'MAINTENANCE' | 'INACTIVE';
    created_at: string;
    updated_at: string;
}

export type RoomStatus = NonNullable<Room['status']>;

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
    data: Room | null;
}

interface RawRoomListResponse {
    success?: boolean;
    message?: string;
    data?: {
        rooms?: Room[];
        Rooms?: Room[];
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

interface RawRoomResponse {
    success?: boolean;
    message?: string;
    data?: {
        Room?: Room;
    } | Room;
}

export const roomApi = baseApi.injectEndpoints({
    endpoints: (builder) => ({
        getRooms: builder.query<ListRoomsResponse, { page?: number; limit?: number; search?: string; sortBy?: string; sortOrder?: string }>({
            query: (params) => ({
                url: '/v1/rooms',
                params,
            }),
            transformResponse: (response: RawRoomListResponse) => {
                const normalized = normalizeListEnvelope<Room>(response, ['rooms', 'Rooms']);
                return {
                    success: normalized.success,
                    message: normalized.message,
                    data: {
                        rooms: normalized.items,
                        pagination: normalized.pagination,
                    },
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
            transformResponse: (response: RawRoomResponse) => normalizeDetailEnvelope<Room>(response, ['Room']),
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
