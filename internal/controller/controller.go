package controller

import (
	"practice/internal/controller/http"

	"go.uber.org/fx"
)

var Module = fx.Options(
	http.Module,
)
