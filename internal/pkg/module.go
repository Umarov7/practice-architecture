package pkg

import (
	"practice/internal/pkg/config"
	"practice/internal/pkg/logger"

	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	logger.Module,
)
