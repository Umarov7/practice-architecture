package service

import (
	"practice/internal/service/computer"
	"practice/internal/service/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	user.Module,
	computer.Module,
)
