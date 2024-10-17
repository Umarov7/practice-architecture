package handler

import (
	"log"
	"log/slog"
	kafkaCons "practice/internal/kafka/consumer"
	kafkaProd "practice/internal/kafka/producer"
	"practice/internal/pkg/config"
	rabbitmqCons "practice/internal/rabbitmq/consumer"
	rabbitmqProd "practice/internal/rabbitmq/producer"
	"practice/internal/repository/mongodb"
	"practice/internal/repository/postgres"
	"practice/internal/service/computer"
	"practice/internal/service/user"

	"go.uber.org/fx"
)

type Handler struct {
	cfg                  *config.Config
	logger               *slog.Logger
	repositoryPostgres   *postgres.Postgres
	repositoryMongo      *mongodb.MongoDB
	serviceUser          user.ServiceUser
	serviceComputer      computer.ServiceComputer
	kafkaProducer        kafkaProd.IKafkaProducer
	rabbitProducer       *rabbitmqProd.MsgBroker
	rabbitConsumer       *rabbitmqCons.MsgBroker
	topicUserCreated     string
	topicUserUpdated     string
	topicUserDeleted     string
	topicComputerCreated string
	topicComputerUpdated string
	topicComputerDeleted string
	queueUserCreated     string
	queueUserUpdated     string
	queueUserDeleted     string
	queueComputerCreated string
	queueComputerUpdated string
	queueComputerDeleted string
}

type Options struct {
	fx.In
	Cfg                *config.Config
	Logger             *slog.Logger
	RepositoryPostgres *postgres.Postgres
	RepositoryMongo    *mongodb.MongoDB
	ServiceUser        user.ServiceUser
	ServiceComputer    computer.ServiceComputer
	RabbitmqProducer   *rabbitmqProd.MsgBroker
	RabbitmqConsumer   *rabbitmqCons.MsgBroker
}

var Module = fx.Provide(New)

func New(opts Options) *Handler {
	consumers := map[string]func([]byte){
		opts.Cfg.KAFKA_TOPIC_USER_CREATED:     kafkaCons.ConsumeCreateUser(opts.Cfg, opts.ServiceUser),
		opts.Cfg.KAFKA_TOPIC_USER_UPDATED:     kafkaCons.ConsumeUpdateUser(opts.Cfg, opts.ServiceUser),
		opts.Cfg.KAFKA_TOPIC_USER_DELETED:     kafkaCons.ConsumeDeleteUser(opts.Cfg, opts.ServiceUser),
		opts.Cfg.KAFKA_TOPIC_COMPUTER_CREATED: kafkaCons.ConsumeCreateComputer(opts.Cfg, opts.ServiceComputer),
		opts.Cfg.KAFKA_TOPIC_COMPUTER_UPDATED: kafkaCons.ConsumeUpdateComputer(opts.Cfg, opts.ServiceComputer),
		opts.Cfg.KAFKA_TOPIC_COMPUTER_DELETED: kafkaCons.ConsumeDeleteComputer(opts.Cfg, opts.ServiceComputer),
	}

	for topic, handler := range consumers {
		consumer := kafkaCons.NewKafkaConsumer([]string{opts.Cfg.KAFKA_ADDRESS}, topic)
		go func(t string, h func([]byte)) {
			log.Printf("Starting consumer for topic: %s", t)
			consumer.Consume(h)
		}(topic, handler)
	}

	return &Handler{
		cfg:                  opts.Cfg,
		logger:               opts.Logger,
		repositoryPostgres:   opts.RepositoryPostgres,
		repositoryMongo:      opts.RepositoryMongo,
		serviceUser:          opts.ServiceUser,
		serviceComputer:      opts.ServiceComputer,
		kafkaProducer:        kafkaProd.NewKafkaProducer([]string{opts.Cfg.KAFKA_ADDRESS}),
		rabbitProducer:       opts.RabbitmqProducer,
		rabbitConsumer:       opts.RabbitmqConsumer,
		topicUserCreated:     opts.Cfg.KAFKA_TOPIC_USER_CREATED,
		topicUserUpdated:     opts.Cfg.KAFKA_TOPIC_USER_UPDATED,
		topicUserDeleted:     opts.Cfg.KAFKA_TOPIC_USER_DELETED,
		topicComputerCreated: opts.Cfg.KAFKA_TOPIC_COMPUTER_CREATED,
		topicComputerUpdated: opts.Cfg.KAFKA_TOPIC_COMPUTER_UPDATED,
		topicComputerDeleted: opts.Cfg.KAFKA_TOPIC_COMPUTER_DELETED,
		queueUserCreated:     opts.Cfg.RabbitMQ_QUEUE_USER_CREATED,
		queueUserUpdated:     opts.Cfg.RabbitMQ_QUEUE_USER_UPDATED,
		queueUserDeleted:     opts.Cfg.RabbitMQ_QUEUE_USER_DELETED,
		queueComputerCreated: opts.Cfg.RabbitMQ_QUEUE_COMPUTER_CREATED,
		queueComputerUpdated: opts.Cfg.RabbitMQ_QUEUE_COMPUTER_UPDATED,
		queueComputerDeleted: opts.Cfg.RabbitMQ_QUEUE_COMPUTER_DELETED,
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
