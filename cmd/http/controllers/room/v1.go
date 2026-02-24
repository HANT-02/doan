package room

import (
	"doan/cmd/http/rest"
	"doan/internal/usecases/room"
	"net/http"
	"strconv"

	"doan/pkg/logger"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	createRoomUseCase room.CreateRoomUseCase
	getRoomUseCase    room.GetRoomUseCase
	updateRoomUseCase room.UpdateRoomUseCase
	deleteRoomUseCase room.DeleteRoomUseCase
	listRoomsUseCase  room.ListRoomsUseCase
}

func NewRoomControllerV1(
	createRoomUseCase room.CreateRoomUseCase,
	getRoomUseCase room.GetRoomUseCase,
	updateRoomUseCase room.UpdateRoomUseCase,
	deleteRoomUseCase room.DeleteRoomUseCase,
	listRoomsUseCase room.ListRoomsUseCase,
) *ControllerV1 {
	return &ControllerV1{
		createRoomUseCase: createRoomUseCase,
		getRoomUseCase:    getRoomUseCase,
		updateRoomUseCase: updateRoomUseCase,
		deleteRoomUseCase: deleteRoomUseCase,
		listRoomsUseCase:  listRoomsUseCase,
	}
}

func (ctrl *ControllerV1) CreateRoom(c *gin.Context) {
	ctxLogger := logger.NewLogger(c.Request.Context())
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctxLogger.Errorf("Failed to bind request: %v", err)
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.createRoomUseCase.Execute(c.Request.Context(), room.CreateRoomInput{
		Name:     req.Name,
		Capacity: req.Capacity,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to create room", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusCreated, "Room created successfully", output.Room)
}

func (ctrl *ControllerV1) GetRoom(c *gin.Context) {
	id := c.Param("id")
	output, err := ctrl.getRoomUseCase.Execute(c.Request.Context(), room.GetRoomInput{ID: id})
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "Room not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Room retrieved successfully", output.Room)
}

func (ctrl *ControllerV1) UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	output, err := ctrl.updateRoomUseCase.Execute(c.Request.Context(), room.UpdateRoomInput{
		ID:       id,
		Name:     req.Name,
		Capacity: req.Capacity,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to update room", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Room updated successfully", output.Room)
}

func (ctrl *ControllerV1) DeleteRoom(c *gin.Context) {
	id := c.Param("id")

	output, err := ctrl.deleteRoomUseCase.Execute(c.Request.Context(), room.DeleteRoomInput{ID: id})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to delete room", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, output.Message, nil)
}

func (ctrl *ControllerV1) ListRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")

	output, err := ctrl.listRoomsUseCase.Execute(c.Request.Context(), room.ListRoomsInput{
		Search:    search,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})

	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list rooms", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Rooms retrieved successfully", output)
}
