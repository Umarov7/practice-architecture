package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"practice/internal/pkg/config"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type Postgres struct {
	Cfg    *config.Config
	DB     *sql.DB
	Logger *slog.Logger
}

type Options struct {
	fx.In
	fx.Lifecycle
	Config *config.Config
	Logger *slog.Logger
}

var Module = fx.Options(fx.Provide(New))

func New(opts Options) *Postgres {
	var db *sql.DB

	pg := &Postgres{
		Cfg:    opts.Config,
		DB:     db,
		Logger: opts.Logger,
	}

	opts.Lifecycle.Append(fx.Hook{
		OnStart: pg.onStart,
		OnStop:  pg.onStop,
	})

	return pg
}

func (r *Postgres) onStart(ctx context.Context) error {
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		r.Cfg.Postgres_HOST, r.Cfg.Postgres_PORT,
		r.Cfg.Postgres_USER, r.Cfg.Postgres_NAME,
		r.Cfg.Postgres_PASSWORD,
	)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return errors.Wrap(err, "error while connecting to postgres")
	}

	if err := db.Ping(); err != nil {
		return errors.Wrap(err, "error while pinging postgres")
	}

	db.

	r.DB = db

	r.Logger.Info("connected to postgres")
	return nil
}

func (r *Postgres) onStop(ctx context.Context) error {
	return r.DB.Close()
}
