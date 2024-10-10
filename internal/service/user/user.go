package user

import (
	"context"
	"errors"
	"log/slog"
	"practice/internal/repository/postgres/user"

	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Options struct {
	fx.In
	Logger *slog.Logger

	UserRepository user.RepositoryUser
}

type Service struct {
	logger   *slog.Logger
	repoUser user.RepositoryUser
}

func New(opts Options) ServiceUser {
	return &Service{
		logger:   opts.Logger,
		repoUser: opts.UserRepository,
	}
}

type ServiceUser interface {
	Create(ctx context.Context, user *user.User) (*user.User, error)
	Read(ctx context.Context, userID string) (*user.User, error)
	Update(ctx context.Context, user *user.User) (string, error)
	Delete(ctx context.Context, userID string) (string, error)
}

func (s *Service) Create(ctx context.Context, user *user.User) (*user.User, error) {
	if !s.validUser(user) {
		return nil, errors.New("invalid user")
	}

	return s.repoUser.Create(ctx, user)
}

func (s *Service) Read(ctx context.Context, userID string) (*user.User, error) {
	if userID == "" {
		return nil, errors.New("userID not exists")
	}

	return s.repoUser.Read(ctx, userID)
}

func (s *Service) Update(ctx context.Context, user *user.User) (string, error) {
	if !s.validUser(user) {
		return "", errors.New("invalid user")
	}

	return s.repoUser.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", errors.New("userID not exists")
	}

	return s.repoUser.Delete(ctx, userID)
}

func (s *Service) validUser(user *user.User) bool {
	if user == nil {
		s.logger.Error("user is nil")
		return false
	}

	if user.ID == "" {
		s.logger.Error("user.ID not exists")
		return false
	}

	if user.Name == "" {
		s.logger.Error("user.Name not exists")
		return false
	}

	if user.Age < 1 {
		s.logger.Error("user.Age not exists")
		return false
	}

	if user.Email == "" {
		s.logger.Error("user.Email not exists")
		return false
	}

	return true
}
