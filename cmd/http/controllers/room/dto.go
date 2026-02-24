package room

import (
	"time"
)

type CreateRoomRequest struct {
	Name     string  `json:"name" binding:"required"`
	Capacity int     `json:"capacity" binding:"required,min=1"`
	Location *string `json:"location"`
	Status   string  `json:"status" binding:"required,oneof=ACTIVE MAINTENANCE INACTIVE"`
}

type UpdateRoomRequest struct {
	Name     string  `json:"name" binding:"required"`
	Capacity int     `json:"capacity" binding:"required,min=1"`
	Location *string `json:"location"`
	Status   string  `json:"status" binding:"required,oneof=ACTIVE MAINTENANCE INACTIVE"`
}

type RoomResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	Location  *string   `json:"location"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListRoomsResponse struct {
	Rooms      []RoomResponse `json:"rooms"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	ItemsPerPage int   `json:"items_per_page"`
	TotalItems   int64 `json:"total_items"`
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
}
