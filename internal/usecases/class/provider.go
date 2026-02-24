package class

import "github.com/google/wire"

var ClassUseCaseProviders = wire.NewSet(
	NewCreateClassUseCase,
	NewGetClassUseCase,
	NewUpdateClassUseCase,
	NewDeleteClassUseCase,
	NewListClassesUseCase,
)
