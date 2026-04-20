package controllers

import (
	"doan/cmd/http/controllers/academic"
	"doan/cmd/http/controllers/account"
	"doan/cmd/http/controllers/class"
	"doan/cmd/http/controllers/course"
	"doan/cmd/http/controllers/leave"
	"doan/cmd/http/controllers/lesson"
	"doan/cmd/http/controllers/material"
	"doan/cmd/http/controllers/predictive"
	"doan/cmd/http/controllers/program"
	"doan/cmd/http/controllers/room"
	"doan/cmd/http/controllers/scheduling"
	"doan/cmd/http/controllers/shift"
	"doan/cmd/http/controllers/student"
	"doan/cmd/http/controllers/studentportal"
	"doan/cmd/http/controllers/teacher"
	"doan/cmd/http/controllers/teacherportal"
	"doan/cmd/http/controllers/user"

	"github.com/google/wire"
)

// ControllerProviders provides all HTTP controllers with interface bindings
var ControllerProviders = wire.NewSet(
	account.NewAccountControllerV1,
	wire.Bind(new(account.Controller), new(*account.ControllerV1)),

	academic.NewAcademicControllerV1,
	wire.Bind(new(academic.Controller), new(*academic.ControllerV1)),

	// User controllers
	user.NewUserControllerV1,
	user.NewUserControllerV2,

	// Class controller
	class.NewClassControllerV1,
	wire.Bind(new(class.Controller), new(*class.ControllerV1)),

	// Room controller
	room.NewRoomControllerV1,
	wire.Bind(new(room.Controller), new(*room.ControllerV1)),

	// Shift controller
	shift.NewShiftControllerV1,
	wire.Bind(new(shift.Controller), new(*shift.ControllerV1)),

	// Teacher controller
	teacher.NewTeacherControllerV1,
	wire.Bind(new(teacher.Controller), new(*teacher.ControllerV1)),

	teacherportal.NewTeacherPortalControllerV1,
	wire.Bind(new(teacherportal.Controller), new(*teacherportal.ControllerV1)),

	// Student controller
	student.NewStudentControllerV1,
	wire.Bind(new(student.Controller), new(*student.ControllerV1)),

	studentportal.NewStudentPortalControllerV1,
	wire.Bind(new(studentportal.Controller), new(*studentportal.ControllerV1)),

	// Course controller
	course.NewCourseControllerV1,
	wire.Bind(new(course.Controller), new(*course.ControllerV1)),

	// Program controller
	program.NewProgramControllerV1,
	wire.Bind(new(program.Controller), new(*program.ControllerV1)),

	// Scheduling controller
	scheduling.NewSchedulingControllerV1,
	wire.Bind(new(scheduling.Controller), new(*scheduling.ControllerV1)),

	// Material controller
	material.NewMaterialControllerV1,
	wire.Bind(new(material.Controller), new(*material.ControllerV1)),

	// Predictive controller
	predictive.NewPredictiveControllerV1,
	wire.Bind(new(predictive.Controller), new(*predictive.ControllerV1)),

	// Lesson controller
	lesson.NewLessonControllerV1,
	wire.Bind(new(lesson.Controller), new(*lesson.ControllerV1)),

	leave.NewLeaveControllerV1,
	wire.Bind(new(leave.Controller), new(*leave.ControllerV1)),
)
