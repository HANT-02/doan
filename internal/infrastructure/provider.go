package infrastructure

import (
	"doan/internal/infrastructure/database"
	_interface "doan/internal/infrastructure/queue/interface"
	"doan/internal/infrastructure/queue/noop"

	"github.com/google/wire"
)

// InfrastructureProviders provides all infrastructure dependencies
// Including: Database, Queue, External services
var InfrastructureProviders = wire.NewSet(
	// Database layer
	database.DBProvider,

	// Queue infrastructure
	ProvideQueue,
)

// ProvideQueue provides queue implementation (noop for now)
func ProvideQueue() _interface.Queue {
	return noop.New()
}
