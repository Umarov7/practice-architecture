package handler

import (
	"log/slog"
	"practice/internal/pkg/config"
	"practice/internal/repository/mongodb"
	"practice/internal/repository/postgres"
	"practice/internal/service/computer"
	"practice/internal/service/user"

	"go.uber.org/fx"
)

type Handler struct {
	cfg                *config.Config
	logger             *slog.Logger
	repositoryPostgres *postgres.Postgres
	repositoryMongo    *mongodb.MongoDB
	serviceUser        user.ServiceUser
	serviceComputer    computer.ServiceComputer
}

type Options struct {
	fx.In
	Cfg                *config.Config
	Logger             *slog.Logger
	RepositoryPostgres *postgres.Postgres
	RepositoryMongo    *mongodb.MongoDB
	ServiceUser        user.ServiceUser
	ServiceComputer    computer.ServiceComputer
}

var Module = fx.Provide(New)

func New(opts Options) *Handler {
	return &Handler{
		cfg:                opts.Cfg,
		logger:             opts.Logger,
		repositoryPostgres: opts.RepositoryPostgres,
		repositoryMongo:    opts.RepositoryMongo,
		serviceUser:        opts.ServiceUser,
		serviceComputer:    opts.ServiceComputer,
	}
}

type UserReq struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type ComputerReq struct {
	IP           string `json:"ip"`
	Manufacturer string `json:"manufacturer"`
	CPU          string `json:"cpu" bson:"cpu"`
	RAM          string `json:"ram" bson:"ram"`
	HDD          string `json:"hdd" bson:"hdd"`
	GPU          string `json:"gpu" bson:"gpu"`
	OS           string `json:"os" bson:"os"`
}
