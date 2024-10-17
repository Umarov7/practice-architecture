package app

import (
	"practice/internal/controller"
	"practice/internal/rabbitmq"
	"practice/internal/repository"
	"practice/internal/service"

	"go.uber.org/fx"
)

func New(opt fx.Option) *fx.App {
	return fx.New(
		opt,
		repository.Module,
		service.Module,
		controller.Module,
		rabbitmq.Module,
	)
}
