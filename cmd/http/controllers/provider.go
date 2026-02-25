package controllers

import (
	"doan/cmd/http/controllers/class"
	"doan/cmd/http/controllers/room"
	"doan/cmd/http/controllers/student"
	"doan/cmd/http/controllers/teacher"
	"doan/cmd/http/controllers/user"

	"github.com/google/wire"
)

// ControllerProviders provides all HTTP controllers with interface bindings
var ControllerProviders = wire.NewSet(
	// User controllers
	user.NewUserControllerV1,
	user.NewUserControllerV2,

	// Class controller
	class.NewClassControllerV1,
	wire.Bind(new(class.Controller), new(*class.ControllerV1)),

	// Room controller
	room.NewRoomControllerV1,
	wire.Bind(new(room.Controller), new(*room.ControllerV1)),

	// Teacher controller
	teacher.NewTeacherControllerV1,
	wire.Bind(new(teacher.Controller), new(*teacher.ControllerV1)),

	// Student controller
	student.NewStudentControllerV1,
	wire.Bind(new(student.Controller), new(*student.ControllerV1)),
)
