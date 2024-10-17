package producer

import (
	"log/slog"
	"practice/internal/pkg/config"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Options struct {
	fx.In

	Logger *slog.Logger
	Cfg    *config.Config
}

type MsgBroker struct {
	channel *amqp.Channel
	logger  *slog.Logger
}

func NewChannel(cfg *config.Config) (*amqp.Channel, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ_ADDRESS)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func New(opts Options) *MsgBroker {
	ch, err := NewChannel(opts.Cfg)
	if err != nil {
		return nil
	}

	return &MsgBroker{
		channel: ch,
		logger:  opts.Logger,
	}
}

func (m *MsgBroker) Publish(queueName string, body []byte) error {
	err := m.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		m.logger.Error("failed to publish to rabbitmq", "queue", queueName, "error", err.Error())
		return err
	}

	m.logger.Info("published to rabbitmq", "queue", queueName)
	return nil
}
