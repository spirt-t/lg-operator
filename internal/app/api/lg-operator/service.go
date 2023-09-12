package lg_operator

import (
	"context"

	"github.com/spirt-t/lg-operator/internal/config"
	"github.com/spirt-t/lg-operator/internal/k8s"
	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
	"go.uber.org/zap"
)

// Service - load-generator service implementation.
type Service struct {
	desc.UnimplementedLoadGeneratorOperatorServiceServer
	k8s            k8s.Manager
	config         config.Manager
	logger         *zap.Logger
	resourceMapper *ResourceMapper
	cleaners       []Cleaner
}

//go:generate mockgen -source=./service.go -destination=./mock/service.go

// Cleaner - clean generators.
type Cleaner interface {
	Run(ctx context.Context) error
}

// NewService - constructor for Service.
func NewService(k8s k8s.Manager, config config.Manager, lg *zap.Logger, cleaners []Cleaner) *Service {
	return &Service{
		k8s:            k8s,
		config:         config,
		logger:         lg,
		resourceMapper: NewResourceMapper(config),
		cleaners:       cleaners,
	}
}

// RunCleaning ...
func (s *Service) RunCleaning(ctx context.Context) {
	for _, cleaner := range s.cleaners {
		go func(cleaner Cleaner) {
			if err := cleaner.Run(ctx); err != nil {
				s.logger.Error("fail to clean generators", zap.Error(err))
			}
		}(cleaner)
	}
}
