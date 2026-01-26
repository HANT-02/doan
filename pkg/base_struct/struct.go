package base_struct

import (
	"doan/pkg/config"
	"doan/pkg/logger"
)

type BaseDependency struct {
	Log           logger.Logger
	ConfigManager config.Manager
}
