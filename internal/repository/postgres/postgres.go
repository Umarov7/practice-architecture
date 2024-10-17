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

	// if err := migrateDB(ctx, db); err != nil {
	// 	return errors.Wrap(err, "error while migrating postgres")
	// }

	r.DB = db

	r.Logger.Info("connected to postgres")
	return nil
}

func (r *Postgres) onStop(ctx context.Context) error {
	return r.DB.Close()
}

func migrateDB(ctx context.Context, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "error while beginning transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query1 := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(50) NOT NULL,
		age INT NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		is_deleted BOOLEAN DEFAULT FALSE
	);
	`

	query2 := `
	INSERT INTO users (id, name, age, email) VALUES
	('550e8400-e29b-41d4-a716-446655440000', 'Alice Johnson', 28, 'alice.johnson@example.com'),
	('550e8400-e29b-41d4-a716-446655440001', 'Bob Smith', 34, 'bob.smith@example.com'),
	('550e8400-e29b-41d4-a716-446655440002', 'Carol Williams', 22, 'carol.williams@example.com'),
	('550e8400-e29b-41d4-a716-446655440003', 'David Brown', 45, 'david.brown@example.com'),
	('550e8400-e29b-41d4-a716-446655440004', 'Emma Davis', 30, 'emma.davis@example.com'),
	('550e8400-e29b-41d4-a716-446655440005', 'Frank Miller', 50, 'frank.miller@example.com'),
	('550e8400-e29b-41d4-a716-446655440006', 'Grace Wilson', 29, 'grace.wilson@example.com'),
	('550e8400-e29b-41d4-a716-446655440007', 'Hannah Moore', 41, 'hannah.moore@example.com'),
	('550e8400-e29b-41d4-a716-446655440008', 'Isaac Taylor', 35, 'isaac.taylor@example.com'),
	('550e8400-e29b-41d4-a716-446655440009', 'Jasmine Anderson', 27, 'jasmine.anderson@example.com');
	`

	_, err = tx.ExecContext(ctx, query1)
	if err != nil {
		return errors.Wrap(err, "error while creating table")
	}

	_, err = tx.ExecContext(ctx, query2)
	if err != nil {
		return errors.Wrap(err, "error while inserting data")
	}

	return tx.Commit()
}
