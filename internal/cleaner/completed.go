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
	coreV1 "k8s.io/api/core/v1"
)

const (
	completedLGCleaningEnabledKey  = "cleaning.completed.enabled"
	completedLGCleaningIntervalKey = "cleaning.completed.interval"
)

var (
	terminalStatuses = []coreV1.PodPhase{
		coreV1.PodFailed,
		coreV1.PodSucceeded,
	}
)

// CompletedLGCleaner - delete completed generators at interval specified in the config.
type CompletedLGCleaner struct {
	config config.Manager
	k8s    k8s.Manager
	logger *zap.Logger
}

// NewCompletedLGCleaner constructor for RegularCleaner.
func NewCompletedLGCleaner(config config.Manager, k8s k8s.Manager, logger *zap.Logger) *CompletedLGCleaner {
	return &CompletedLGCleaner{
		config: config,
		k8s:    k8s,
		logger: logger,
	}
}

// Run - clean completed generators regular.
func (rc *CompletedLGCleaner) Run(ctx context.Context) error {
	for {
		interval, err := rc.interval()
		if err != nil {
			return err
		}

		if rc.enabled() {
			if err = rc.regularCleaning(ctx); err != nil {
				rc.logger.Error("failed to clean completed generators", zap.Error(err))
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}

func (rc *CompletedLGCleaner) enabled() bool {
	var enabled bool
	_ = rc.config.UnmarshalKey(completedLGCleaningEnabledKey, &enabled)

	return enabled
}

func (rc *CompletedLGCleaner) interval() (time.Duration, error) {
	var intervalStr string
	if err := rc.config.UnmarshalKey(completedLGCleaningIntervalKey, &intervalStr); err != nil {
		return 0, fmt.Errorf("failed to define interval: %w", err)
	}

	return time.ParseDuration(intervalStr)
}

func (rc *CompletedLGCleaner) regularCleaning(ctx context.Context) error {
	var namesToDelete []string

	rc.logger.Info("Start cleaning completed generators")
	defer func() {
		rc.logger.Info("Deleted completed generators", zap.String("generators", strings.Join(namesToDelete, ",")))
	}()

	generators, err := rc.k8s.List(ctx)
	if err != nil {
		return err
	}

	namesToDelete = rc.namesToDelete(generators)
	if len(namesToDelete) == 0 {
		return nil
	}

	ch := make(chan error)
	for _, name := range namesToDelete {
		go func(name string) {
			ch <- rc.k8s.Delete(ctx, name)
		}(name)
	}

	for i := 0; i < len(namesToDelete); i++ {
		err = multierr.Append(err, <-ch)
	}

	return err
}

func (rc *CompletedLGCleaner) namesToDelete(list []model.LoadGenerator) []string {
	var namesToDelete []string
	for _, generator := range list {
		for _, terminalStatus := range terminalStatuses {
			if terminalStatus == generator.Status {
				namesToDelete = append(namesToDelete, generator.Name)
				break
			}
		}
	}

	return namesToDelete
}
