package cleaner

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/spirt-t/lg-operator/internal/config"
	mock_k8s "github.com/spirt-t/lg-operator/internal/k8s/mock"
	"github.com/spirt-t/lg-operator/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	coreV1 "k8s.io/api/core/v1"
)

func Test_outdated_namesToDelete(t *testing.T) {
	t.Parallel()
	rc := OutdatedLGCleaner{}

	t.Run("with completed pods", func(t *testing.T) {
		t.Parallel()

		names := rc.namesToDelete([]model.LoadGenerator{
			{
				Name:       "lg-1",
				ClusterIP:  "127.0.0.1",
				ExternalIP: "127.0.0.2",
				Port:       1000,
				Status:     coreV1.PodFailed,
				CreatedAt:  time.Now().Add(-time.Hour),
			},
			{
				Name:       "lg-2",
				ClusterIP:  "127.0.1.1",
				ExternalIP: "127.0.1.2",
				Port:       1001,
				Status:     coreV1.PodRunning,
				CreatedAt:  time.Now().Add(-22 * time.Hour),
			},
			{
				Name:       "lg-3",
				ClusterIP:  "127.0.2.1",
				ExternalIP: "127.0.2.2",
				Port:       1002,
				Status:     coreV1.PodPending,
				CreatedAt:  time.Now().Add(-25 * time.Hour),
			},
			{
				Name:       "lg-4",
				ClusterIP:  "127.0.3.1",
				ExternalIP: "127.0.3.2",
				Port:       1003,
				Status:     coreV1.PodSucceeded,
				CreatedAt:  time.Now().Add(-30 * time.Hour),
			},
		}, 24*time.Hour)
		assert.Equal(t, 2, len(names))
		assert.Equal(t, []string{"lg-3", "lg-4"}, names)
	})

	t.Run("without completed pods", func(t *testing.T) {
		t.Parallel()

		names := rc.namesToDelete([]model.LoadGenerator{
			{
				Name:       "lg-1",
				ClusterIP:  "127.0.0.1",
				ExternalIP: "127.0.0.2",
				Port:       1000,
				Status:     coreV1.PodRunning,
				CreatedAt:  time.Now().Add(-time.Minute * 30),
			},
			{
				Name:       "lg-2",
				ClusterIP:  "127.0.1.1",
				ExternalIP: "127.0.1.2",
				Port:       1001,
				Status:     coreV1.PodRunning,
				CreatedAt:  time.Now().Add(-time.Hour),
			},
			{
				Name:       "lg-3",
				ClusterIP:  "127.0.2.1",
				ExternalIP: "127.0.2.2",
				Port:       1002,
				Status:     coreV1.PodPending,
				CreatedAt:  time.Now().Add(-2 * time.Hour),
			},
			{
				Name:       "lg-4",
				ClusterIP:  "127.0.3.1",
				ExternalIP: "127.0.3.2",
				Port:       1003,
				Status:     coreV1.PodRunning,
				CreatedAt:  time.Now().Add(-23 * time.Hour),
			},
		}, 24*time.Hour)
		assert.Equal(t, 0, len(names))
	})

	t.Run("nil list", func(t *testing.T) {
		t.Parallel()

		names := rc.namesToDelete(nil, time.Hour)
		assert.Equal(t, 0, len(names))
	})
}

func Test_outdated_regularCleaning(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	k8sManager := mock_k8s.NewMockManager(ctrl)
	mngr, err := config.NewManager("../../testfiles/test_config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	l := zaptest.NewLogger(t)

	rc := OutdatedLGCleaner{mngr, k8sManager, l}

	t.Run("with completed pods", func(t *testing.T) {
		k8sManager.EXPECT().List(ctx).Return([]model.LoadGenerator{
			{
				Name:       "lg-1",
				ClusterIP:  "127.0.0.1",
				ExternalIP: "127.0.0.2",
				Port:       1000,
				Status:     coreV1.PodFailed,
				CreatedAt:  time.Now().Add(-25 * time.Hour),
			},
			{
				Name:       "lg-2",
				ClusterIP:  "127.0.1.1",
				ExternalIP: "127.0.1.2",
				Port:       1001,
				Status:     coreV1.PodRunning,
				CreatedAt:  time.Now().Add(-23 * time.Hour),
			},
			{
				Name:       "lg-3",
				ClusterIP:  "127.0.2.1",
				ExternalIP: "127.0.2.2",
				Port:       1002,
				Status:     coreV1.PodPending,
				CreatedAt:  time.Now().Add(-2 * time.Hour),
			},
			{
				Name:       "lg-4",
				ClusterIP:  "127.0.3.1",
				ExternalIP: "127.0.3.2",
				Port:       1003,
				Status:     coreV1.PodSucceeded,
				CreatedAt:  time.Now().Add(-30 * time.Hour),
			},
		}, nil)
		k8sManager.EXPECT().Delete(ctx, "lg-1").Return(nil)
		k8sManager.EXPECT().Delete(ctx, "lg-4").Return(nil)

		err = rc.regularCleaning(ctx, 24*time.Hour)
		assert.NoError(t, err)
	})

	t.Run("without completed pods", func(t *testing.T) {
		k8sManager.EXPECT().List(ctx).Return([]model.LoadGenerator{
			{
				Name:       "lg-1",
				ClusterIP:  "127.0.0.1",
				ExternalIP: "127.0.0.2",
				Port:       1000,
				Status:     coreV1.PodRunning,
				CreatedAt:  time.Now().Add(-time.Hour),
			},
			{
				Name:       "lg-2",
				ClusterIP:  "127.0.1.1",
				ExternalIP: "127.0.1.2",
				Port:       1001,
				Status:     coreV1.PodRunning,
				CreatedAt:  time.Now().Add(-12 * time.Hour),
			},
			{
				Name:       "lg-3",
				ClusterIP:  "127.0.2.1",
				ExternalIP: "127.0.2.2",
				Port:       1002,
				Status:     coreV1.PodPending,
				CreatedAt:  time.Now().Add(-time.Minute * 30),
			},
			{
				Name:       "lg-4",
				ClusterIP:  "127.0.3.1",
				ExternalIP: "127.0.3.2",
				Port:       1003,
				Status:     coreV1.PodRunning,
				CreatedAt:  time.Now().Add(-2 * time.Hour),
			},
		}, nil)

		err = rc.regularCleaning(ctx, 24*time.Hour)
		assert.NoError(t, err)
	})

	t.Run("nil pods list", func(t *testing.T) {
		k8sManager.EXPECT().List(ctx).Return(nil, nil)

		err = rc.regularCleaning(ctx, 24*time.Hour)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		k8sManager.EXPECT().List(ctx).Return(nil, errors.New("some error"))

		err = rc.regularCleaning(ctx, 24*time.Hour)
		assert.NotNil(t, err)
	})
}
