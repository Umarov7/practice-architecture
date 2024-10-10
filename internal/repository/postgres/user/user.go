package user

import (
	"context"
	"database/sql"
	"log/slog"
	"practice/internal/pkg/config"
	"practice/internal/repository/postgres"

	"github.com/pkg/errors"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type RepositoryUser interface {
	Create(ctx context.Context, user *User) (*User, error)
	Read(ctx context.Context, userID string) (*User, error)
	Update(ctx context.Context, user *User) (string, error)
	Delete(ctx context.Context, userID string) (string, error)
}

type Repository struct {
	repo   *postgres.Postgres
	logger *slog.Logger
}

type Options struct {
	fx.In
	fx.Lifecycle
	Cfg      *config.Config
	Postgres *postgres.Postgres
	Logger   *slog.Logger
}

var _ RepositoryUser = (*Repository)(nil)

func New(opts Options) RepositoryUser {
	repo := &Repository{
		logger: opts.Logger,
	}

	opts.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			repo.repo = opts.Postgres
			return nil
		},
		OnStop: func(context.Context) error { return nil },
	})

	return repo
}

func (r *Repository) Create(ctx context.Context, user *User) (*User, error) {
	query := `
	insert into users
		(id, name, age, email) 
	values
		($1, $2, $3, $4)
	returning *
	`

	_, err := r.repo.DB.ExecContext(ctx, query, user.ID, user.Name, user.Age, user.Email)
	if err != nil {
		return nil, errors.Wrap(err, "error while inserting user")
	}

	return user, nil
}

func (r *Repository) Read(ctx context.Context, userID string) (*User, error) {
	query := `
	select
		name, age, email
	from
		users
	where
		id = $1 and is_deleted = false
	`

	u := User{ID: userID, IsDeleted: false}
	err := r.repo.DB.QueryRowContext(ctx, query, userID).Scan(&u.Name, &u.Age, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrap(err, "not found")
		}

		return nil, errors.Wrap(err, "error while finding user")
	}

	return &u, nil
}

func (r *Repository) Update(ctx context.Context, user *User) (string, error) {
	query := `
	update
		users
	set
		name = $2, age = $3, email = $4
	where
		id = $1 and is_deleted = false
	`

	_, err := r.repo.DB.ExecContext(ctx, query, user.ID, user.Name, user.Age, user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.Wrap(err, "not found")
		}
		return "", errors.Wrap(err, "error while updating user")
	}

	return user.ID, nil
}

func (r *Repository) Delete(ctx context.Context, userID string) (string, error) {
	query := `
	update
		users
	set
		is_deleted = true
	where
		id = $1
	`

	_, err := r.repo.DB.ExecContext(ctx, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.Wrap(err, "not found")
		}
		return "", errors.Wrap(err, "error while deleting user")
	}

	return userID, nil
}
