package lg_operator

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spirt-t/lg-operator/internal/config"
	mock_k8s "github.com/spirt-t/lg-operator/internal/k8s/mock"
	"github.com/spirt-t/lg-operator/internal/model"
	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestService_GeneratorsList(t *testing.T) {
	l := zaptest.NewLogger(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	k8sManager := mock_k8s.NewMockManager(ctrl)
	mngr, err := config.NewManager("../../../../testfiles/test_config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := NewService(k8sManager, mngr, l, nil)

	t.Run("empty list", func(t *testing.T) {
		k8sManager.EXPECT().List(ctx).Return(nil, nil)

		res, err := s.GeneratorsList(ctx, &desc.GeneratorsListRequest{})
		assert.NoError(t, err)
		assert.Equal(t, 0, len(res.LoadGenerators))
	})

	t.Run("ok", func(t *testing.T) {
		k8sManager.EXPECT().List(ctx).Return([]model.LoadGenerator{
			{
				Name:      "generator-1",
				ClusterIP: "127.0.0.1",
				Port:      8888,
				Status:    "Running",
			},
			{
				Name:      "generator-2",
				ClusterIP: "127.0.0.2",
				Port:      8888,
				Status:    "Running",
			},
		}, nil)

		res, err := s.GeneratorsList(ctx, &desc.GeneratorsListRequest{})
		assert.NoError(t, err)

		assert.Equal(t, 2, len(res.LoadGenerators))
		assert.Equal(t, "generator-1", res.LoadGenerators[0].Name)
		assert.Equal(t, "127.0.0.1", res.LoadGenerators[0].ClusterIp)
		assert.Equal(t, int32(8888), res.LoadGenerators[0].Port)
		assert.Equal(t, "Running", res.LoadGenerators[0].Status)

		assert.Equal(t, "generator-2", res.LoadGenerators[1].Name)
		assert.Equal(t, "127.0.0.2", res.LoadGenerators[1].ClusterIp)
		assert.Equal(t, int32(8888), res.LoadGenerators[1].Port)
		assert.Equal(t, "Running", res.LoadGenerators[1].Status)
	})

	t.Run("error", func(t *testing.T) {
		er := errors.New("some error")
		k8sManager.EXPECT().List(ctx).Return(nil, er)

		res, err := s.GeneratorsList(ctx, &desc.GeneratorsListRequest{})
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, er))
	})
}
