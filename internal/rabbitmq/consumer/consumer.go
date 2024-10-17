package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"practice/internal/pkg/config"
	compRepo "practice/internal/repository/mongodb/computer"
	userRepo "practice/internal/repository/postgres/user"
	"practice/internal/service/computer"
	"practice/internal/service/user"
	"sync"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Options struct {
	fx.In
	fx.Lifecycle
	UserService     user.ServiceUser
	ComputerService computer.ServiceComputer
	Logger          *slog.Logger
	Cfg             *config.Config
}

type MsgBroker struct {
	user           user.ServiceUser
	computer       computer.ServiceComputer
	channel        *amqp.Channel
	logger         *slog.Logger
	cfg            *config.Config
	wg             *sync.WaitGroup
	numOfServices  int
	createUser     <-chan amqp.Delivery
	updateUser     <-chan amqp.Delivery
	deleteUser     <-chan amqp.Delivery
	createComputer <-chan amqp.Delivery
	updateComputer <-chan amqp.Delivery
	deleteComputer <-chan amqp.Delivery
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
		opts.Logger.Error("error initializing amqp channel", err.Error())
		return nil
	}

	msgBroker := &MsgBroker{
		user:          opts.UserService,
		computer:      opts.ComputerService,
		channel:       ch,
		logger:        opts.Logger,
		cfg:           opts.Cfg,
		wg:            &sync.WaitGroup{},
		numOfServices: 6,
	}

	if err := declareQueues(msgBroker); err != nil {
		msgBroker.logger.Error("error declaring queues: %s", err.Error())
		return nil
	}

	msgBroker.wg.Add(msgBroker.numOfServices)
	opts.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go msgBroker.consumeMessages(ctx, msgBroker.createUser, "create_user")
			go msgBroker.consumeMessages(ctx, msgBroker.updateUser, "update_user")
			go msgBroker.consumeMessages(ctx, msgBroker.deleteUser, "delete_user")
			go msgBroker.consumeMessages(ctx, msgBroker.createComputer, "create_computer")
			go msgBroker.consumeMessages(ctx, msgBroker.updateComputer, "update_computer")
			go msgBroker.consumeMessages(ctx, msgBroker.deleteComputer, "delete_computer")

			log.Println("RabbitMQ consuming")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			msgBroker.wg.Wait()
			return msgBroker.channel.Close()
		},
	})

	return msgBroker
}

func (m *MsgBroker) consumeMessages(ctx context.Context, messages <-chan amqp.Delivery, logPrefix string) {
	defer m.wg.Done()
	log.Printf("Consuming %s queue", logPrefix)
	
	for {
		select {
		case msg := <-messages:
			log.Printf("Received data through RabbitMQ: %s", msg.Body)
			if err := m.processMessage(ctx, logPrefix, msg); err != nil {
				handleError(m, msg, logPrefix, err, "processing message")
				msg.Nack(false, false)
				continue
			}
			msg.Ack(false)

		case <-ctx.Done():
			log.Printf("context done, stopping %s\n", logPrefix)
			return
		}
	}
}

func (m *MsgBroker) processMessage(ctx context.Context, logPrefix string, val amqp.Delivery) error {
	var err error

	switch logPrefix {
	case "create_user":
		var req userRepo.User
		err = json.Unmarshal(val.Body, &req)
		if err == nil {
			resp, err2 := m.user.Create(ctx, &req)
			m.logger.Info("created user: %v", resp)
			err = err2
		}
	case "update_user":
		var req userRepo.User
		err = json.Unmarshal(val.Body, &req)
		if err == nil {
			resp, err2 := m.user.Update(ctx, &req)
			m.logger.Info("updated user: %v", resp)
			err = err2
		}
	case "delete_user":
		if err == nil {
			resp, err2 := m.user.Delete(ctx, string(val.Body))
			m.logger.Info("deleted user: %v", resp)
			err = err2
		}
	case "create_computer":
		var req compRepo.Computer
		err = json.Unmarshal(val.Body, &req)
		if err == nil {
			resp, err2 := m.computer.Create(ctx, &req)
			m.logger.Info("created computer: %v", resp)
			err = err2
		}
	case "update_computer":
		var req compRepo.Computer
		err = json.Unmarshal(val.Body, &req)
		if err == nil {
			resp, err2 := m.computer.Update(ctx, &req)
			m.logger.Info("updated computer: %v", resp)
			err = err2
		}
	case "delete_computer":
		if err == nil {
			resp, err2 := m.computer.Delete(ctx, string(val.Body))
			m.logger.Info("deleted computer: %v", resp)
			err = err2
		}
	default:
		return errors.Errorf("unknown log prefix: %s", logPrefix)
	}

	if err != nil {
		return errors.Errorf("error unmarshalling or processing: ", err)
	}
	return nil
}

func declareQueues(m *MsgBroker) error {
	queueNames := []string{
		m.cfg.RabbitMQ_QUEUE_USER_CREATED,
		m.cfg.RabbitMQ_QUEUE_USER_UPDATED,
		m.cfg.RabbitMQ_QUEUE_USER_DELETED,
		m.cfg.RabbitMQ_QUEUE_COMPUTER_CREATED,
		m.cfg.RabbitMQ_QUEUE_COMPUTER_UPDATED,
		m.cfg.RabbitMQ_QUEUE_COMPUTER_DELETED,
	}

	var consumers []<-chan amqp.Delivery

	for _, name := range queueNames {
		q, err := m.channel.QueueDeclare(name, true, false, false, false, nil)
		if err != nil {
			return err
		}

		con, err := m.channel.Consume(q.Name, "", false, false, false, false, nil)
		if err != nil {
			return err
		}

		consumers = append(consumers, con)
	}

	m.createUser = consumers[0]
	m.updateUser = consumers[1]
	m.deleteUser = consumers[2]
	m.createComputer = consumers[3]
	m.updateComputer = consumers[4]
	m.deleteComputer = consumers[5]

	return nil
}

func handleError(m *MsgBroker, val amqp.Delivery, logPrefix string, err error, msg string) {
	m.logger.Error(fmt.Sprintf("error %s in %s: %s\n", msg, logPrefix, err.Error()))
}
