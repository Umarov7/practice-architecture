package repository

import (
	"practice/internal/repository/mongodb"
	"practice/internal/repository/mongodb/computer"
	"practice/internal/repository/postgres"
	"practice/internal/repository/postgres/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	postgres.Module,
	mongodb.Module,
	user.Module,
	computer.Module,
)
