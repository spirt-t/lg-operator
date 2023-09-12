package cleaner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spirt-t/lg-operator/internal/config"
	"github.com/spirt-t/lg-operator/internal/k8s"
	"github.com/spirt-t/lg-operator/internal/model"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	outdatedLGCleaningEnabledKey = "cleaning.outdated.enabled"
	outdatedLGCleaningTTLKey     = "cleaning.outdated.ttl"
)

// OutdatedLGCleaner - delete completed generators at interval specified in the config.
type OutdatedLGCleaner struct {
	config config.Manager
	k8s    k8s.Manager
	logger *zap.Logger
}

// NewOutdatedLGCleaner constructor for RegularCleaner.
func NewOutdatedLGCleaner(config config.Manager, k8s k8s.Manager, logger *zap.Logger) *OutdatedLGCleaner {
	return &OutdatedLGCleaner{
		config: config,
		k8s:    k8s,
		logger: logger,
	}
}

// Run - clean old generators regular.
func (oc *OutdatedLGCleaner) Run(ctx context.Context) error {
	for {
		interval, err := oc.ttl()
		if err != nil {
			return err
		}

		if oc.enabled() {
			if err = oc.regularCleaning(ctx, interval); err != nil {
				oc.logger.Error("failed to clean old generators", zap.Error(err))
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}

func (oc *OutdatedLGCleaner) enabled() bool {
	var enabled bool
	_ = oc.config.UnmarshalKey(outdatedLGCleaningEnabledKey, &enabled)

	return enabled
}

func (oc *OutdatedLGCleaner) ttl() (time.Duration, error) {
	var intervalStr string
	if err := oc.config.UnmarshalKey(outdatedLGCleaningTTLKey, &intervalStr); err != nil {
		return 0, fmt.Errorf("failed to define interval: %w", err)
	}

	return time.ParseDuration(intervalStr)
}

func (oc *OutdatedLGCleaner) regularCleaning(ctx context.Context, ttl time.Duration) error {
	var namesToDelete []string

	oc.logger.Info("Start cleaning old generators")
	defer func() {
		oc.logger.Info("Deleted old generators", zap.String("generators", strings.Join(namesToDelete, ",")))
	}()

	generators, err := oc.k8s.List(ctx)
	if err != nil {
		return err
	}

	namesToDelete = oc.namesToDelete(generators, ttl)
	if len(namesToDelete) == 0 {
		return nil
	}

	ch := make(chan error)
	for _, name := range namesToDelete {
		go func(name string) {
			ch <- oc.k8s.Delete(ctx, name)
		}(name)
	}

	for i := 0; i < len(namesToDelete); i++ {
		err = multierr.Append(err, <-ch)
	}

	return err
}

func (oc *OutdatedLGCleaner) namesToDelete(list []model.LoadGenerator, ttl time.Duration) []string {
	maxOldTime := time.Now().UTC().Add(-ttl)

	var namesToDelete []string
	for _, generator := range list {
		if generator.CreatedAt.Before(maxOldTime) {
			namesToDelete = append(namesToDelete, generator.Name)
		}
	}

	return namesToDelete
}
