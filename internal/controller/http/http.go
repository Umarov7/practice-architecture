package http

import (
	"practice/internal/controller/http/handler"
	"practice/internal/controller/http/router"

	"go.uber.org/fx"
)

var Module = fx.Options(
	router.Module,
	handler.Module,
)
