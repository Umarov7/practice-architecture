package rabbitmq

import (
	"practice/internal/rabbitmq/consumer"
	"practice/internal/rabbitmq/producer"

	"go.uber.org/fx"
)

var Module = fx.Options(
	producer.Module,
	consumer.Module,
)
