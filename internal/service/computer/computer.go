package computer

import (
	"context"
	"errors"
	"log/slog"
	"practice/internal/repository/mongodb/computer"

	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Options struct {
	fx.In
	Logger *slog.Logger

	ComputerRepository computer.RepositoryComputer
}

type Service struct {
	logger       *slog.Logger
	repoComputer computer.RepositoryComputer
}

func New(opts Options) ServiceComputer {
	return &Service{
		logger:       opts.Logger,
		repoComputer: opts.ComputerRepository,
	}
}

type ServiceComputer interface {
	Create(ctx context.Context, computer *computer.Computer) (*computer.Computer, error)
	Read(ctx context.Context, compID string) (*computer.Computer, error)
	Update(ctx context.Context, computer *computer.Computer) (string, error)
	Delete(ctx context.Context, compID string) (string, error)
	GetAll(ctx context.Context) ([]*computer.Computer, error)
}

func (s *Service) Create(ctx context.Context, computer *computer.Computer) (*computer.Computer, error) {
	if !s.validComputer(computer) {
		return nil, errors.New("invalid computer")
	}

	return s.repoComputer.Create(ctx, computer)
}

func (s *Service) Read(ctx context.Context, compID string) (*computer.Computer, error) {
	if compID == "" {
		return nil, errors.New("computerID not exists")
	}

	return s.repoComputer.Read(ctx, compID)
}

func (s *Service) Update(ctx context.Context, computer *computer.Computer) (string, error) {
	if !s.validComputer(computer) || computer.ID == nil {
		return "", errors.New("invalid computer")
	}

	return s.repoComputer.Update(ctx, computer)
}

func (s *Service) Delete(ctx context.Context, compID string) (string, error) {
	if compID == "" {
		return "", errors.New("computerID not exists")
	}

	return s.repoComputer.Delete(ctx, compID)
}

func (s *Service) GetAll(ctx context.Context) ([]*computer.Computer, error) {
	return s.repoComputer.GetAll(ctx)
}

func (s *Service) validComputer(computer *computer.Computer) bool {
	if computer == nil {
		s.logger.Error("computer is nil")
		return false
	}

	if computer.IP == "" {
		s.logger.Error("computer.IP not exists")
		return false
	}

	if computer.Manufacturer == "" {
		s.logger.Error("computer.Manufacturer not exists")
		return false
	}

	if computer.CPU == "" {
		s.logger.Error("computer.CPU not exists")
		return false
	}

	if computer.RAM == "" {
		s.logger.Error("computer.RAM not exists")
		return false
	}

	if computer.HDD == "" {
		s.logger.Error("computer.HDD not exists")
		return false
	}

	if computer.GPU == "" {
		s.logger.Error("computer.GPU not exists")
		return false
	}

	if computer.OS == "" {
		s.logger.Error("computer.OS not exists")
		return false
	}

	return true
}
