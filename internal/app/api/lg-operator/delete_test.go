package lg_operator

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spirt-t/lg-operator/internal/config"
	mock_k8s "github.com/spirt-t/lg-operator/internal/k8s/mock"
	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestService_DeleteGenerator(t *testing.T) {
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

	t.Run("ok", func(t *testing.T) {
		k8sManager.EXPECT().Delete(ctx, "test-generator-name").Return(nil)

		_, err = s.DeleteGenerators(ctx, &desc.DeleteGeneratorsRequest{Names: []string{"test-generator-name"}})
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		er := errors.New("some error")
		k8sManager.EXPECT().Delete(ctx, "test-generator-name").Return(er)

		_, err = s.DeleteGenerators(ctx, &desc.DeleteGeneratorsRequest{Names: []string{"test-generator-name"}})
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, er))
	})
}
