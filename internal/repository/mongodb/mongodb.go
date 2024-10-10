package mongodb

import (
	"context"
	"log/slog"
	"practice/internal/pkg/config"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

type MongoDB struct {
	Cfg    *config.Config
	Client *mongo.Client
	DB     *mongo.Database
	Logger *slog.Logger
}

type Options struct {
	fx.In
	fx.Lifecycle
	Config *config.Config
	Logger *slog.Logger
}

var Module = fx.Options(fx.Provide(New))

func New(opts Options) *MongoDB {
	var (
		client *mongo.Client
		db     *mongo.Database
	)

	mongoDB := &MongoDB{
		Cfg:    opts.Config,
		Client: client,
		DB:     db,
		Logger: opts.Logger,
	}

	opts.Lifecycle.Append(fx.Hook{
		OnStart: mongoDB.onStart,
		OnStop:  mongoDB.onStop,
	})

	return mongoDB
}

func (r *MongoDB) onStart(ctx context.Context) error {
	connectOptions := options.Client().ApplyURI(r.Cfg.MongoDB_URI)

	connect, err := mongo.Connect(ctx, connectOptions)
	if err != nil {
		return errors.Wrap(err, "error while connecting to mongodb")
	}

	if err = connect.Ping(ctx, nil); err != nil {
		return errors.Wrap(err, "error while pinging mongodb")
	}

	r.Client = connect
	r.DB = r.Client.Database(r.Cfg.MongoDB_NAME)

	r.Logger.Info("connected to mongodb")
	return nil
}

func (r *MongoDB) onStop(ctx context.Context) error {
	return r.Client.Disconnect(ctx)
}
