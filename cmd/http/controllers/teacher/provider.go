package teacher

import "github.com/google/wire"

// ControllerProviders provides teacher controller for Wire DI
var ControllerProviders = wire.NewSet(
	NewTeacherControllerV1,
	wire.Bind(new(Controller), new(*ControllerV1)),
)
