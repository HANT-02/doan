package account

import (
	"net/http"
	"strconv"

	"doan/cmd/http/rest"
	accountuc "doan/internal/usecases/account"

	"github.com/gin-gonic/gin"
)

var _ Controller = (*ControllerV1)(nil)

type ControllerV1 struct {
	listUsersUseCase         accountuc.ListUsersUseCase
	getUserUseCase           accountuc.GetUserUseCase
	createUserUseCase        accountuc.CreateUserUseCase
	updateUserUseCase        accountuc.UpdateUserUseCase
	resetUserPasswordUseCase accountuc.ResetUserPasswordUseCase
}

func NewAccountControllerV1(
	listUsersUseCase accountuc.ListUsersUseCase,
	getUserUseCase accountuc.GetUserUseCase,
	createUserUseCase accountuc.CreateUserUseCase,
	updateUserUseCase accountuc.UpdateUserUseCase,
	resetUserPasswordUseCase accountuc.ResetUserPasswordUseCase,
) *ControllerV1 {
	return &ControllerV1{
		listUsersUseCase:         listUsersUseCase,
		getUserUseCase:           getUserUseCase,
		createUserUseCase:        createUserUseCase,
		updateUserUseCase:        updateUserUseCase,
		resetUserPasswordUseCase: resetUserPasswordUseCase,
	}
}

func (ctrl *ControllerV1) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sortBy := c.Query("sortBy")
	sortOrder := c.Query("sortOrder")
	search := c.Query("search")
	role := c.Query("role")

	var isActive *bool
	if activeRaw := c.Query("is_active"); activeRaw != "" {
		parsed, parseErr := strconv.ParseBool(activeRaw)
		if parseErr == nil {
			isActive = &parsed
		}
	}

	output, err := ctrl.listUsersUseCase.Execute(c.Request.Context(), accountuc.ListUsersInput{
		Search:    search,
		Role:      role,
		IsActive:  isActive,
		Page:      page,
		Limit:     limit,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusInternalServerError, "Failed to list users", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Users retrieved successfully", output)
}

func (ctrl *ControllerV1) GetUser(c *gin.Context) {
	output, err := ctrl.getUserUseCase.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		rest.ResponseError(c, http.StatusNotFound, "User not found", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "User retrieved successfully", output)
}

func (ctrl *ControllerV1) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	output, err := ctrl.createUserUseCase.Execute(c.Request.Context(), accountuc.CreateUserInput{
		Code:     req.Code,
		FullName: req.FullName,
		Email:    req.Email,
		Role:     req.Role,
		IsActive: isActive,
		Password: req.Password,
	})
	if err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to create user", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusCreated, "User created successfully", output)
}

func (ctrl *ControllerV1) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := ctrl.updateUserUseCase.Execute(c.Request.Context(), accountuc.UpdateUserInput{
		ID:       c.Param("id"),
		FullName: req.FullName,
		Role:     req.Role,
		IsActive: req.IsActive,
	}); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to update user", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "User updated successfully", nil)
}

func (ctrl *ControllerV1) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := ctrl.resetUserPasswordUseCase.Execute(c.Request.Context(), accountuc.ResetUserPasswordInput{
		ID:          c.Param("id"),
		NewPassword: req.NewPassword,
	}); err != nil {
		rest.ResponseError(c, http.StatusBadRequest, "Failed to reset password", err)
		return
	}

	rest.ResponseSuccess(c, http.StatusOK, "Password reset successfully", nil)
}
