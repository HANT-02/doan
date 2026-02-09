package controllers

import (
	"doan/cmd/http/controllers/teacher"
	"doan/cmd/http/controllers/user"

	"github.com/google/wire"
)

var ControllerProviders = wire.NewSet(
	user.NewUserControllerV1,
	user.NewUserControllerV2,
	teacher.ControllerProviders,
)
