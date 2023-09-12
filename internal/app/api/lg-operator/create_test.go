package lg_operator

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spirt-t/lg-operator/internal/config"
	"github.com/spirt-t/lg-operator/internal/k8s"
	mock_k8s "github.com/spirt-t/lg-operator/internal/k8s/mock"
	"github.com/spirt-t/lg-operator/internal/model"
	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestService_CreateGenerator(t *testing.T) {
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
		lg := model.LoadGenerator{
			Name:       "testname",
			ClusterIP:  "127.1.2.3",
			ExternalIP: "0.1.2.3",
			Port:       8888,
			Status:     "Running",
		}

		k8sManager.EXPECT().Create(gomock.Any(), k8s.CreationConfig{
			Image: "testimage",
			Resources: model.Resources{
				Memory: model.Resource{
					Limit:   "5Gi",
					Request: "4Gi",
				},
				CPU: model.Resource{
					Limit:   "4",
					Request: "3",
				},
			},
			Envs: []model.EnvVar{
				{
					Name:  "testname",
					Value: "testval",
				},
			},
			Commands: []string{"run"},
		}).Return(&lg, nil)

		res, err := s.CreateGenerators(ctx, &desc.CreateGeneratorsRequest{
			Parameters: []*desc.CreateGeneratorsParams{
				{
					Image: "testimage",
					Resources: &desc.Resources{
						Memory: &desc.Resource{
							Limit:   "5Gi",
							Request: "4Gi",
						},
						Cpu: &desc.Resource{
							Limit:   "4",
							Request: "3",
						},
					},
					AdditionalEnvs: []*desc.EnvVar{
						{
							Name: "testname",
							Val:  "testval",
						},
					},
					Commands: []string{"run"},
				},
			},
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotNil(t, res.LoadGenerators[0])
		assert.Equal(t, lg.Name, res.LoadGenerators[0].Name)
		assert.Equal(t, string(lg.Status), res.LoadGenerators[0].Status)
		assert.Equal(t, lg.Port, res.LoadGenerators[0].Port)
		assert.Equal(t, lg.ClusterIP, res.LoadGenerators[0].ClusterIp)
		assert.Equal(t, lg.ExternalIP, res.LoadGenerators[0].ExternalIp)
	})

	t.Run("error", func(t *testing.T) {
		er := errors.New("some error")
		k8sManager.EXPECT().Create(gomock.Any(), k8s.CreationConfig{
			Image: "testimage",
			Resources: model.Resources{
				Memory: model.Resource{
					Limit:   "5Gi",
					Request: "4Gi",
				},
				CPU: model.Resource{
					Limit:   "4",
					Request: "3",
				},
			},
			Envs: []model.EnvVar{
				{
					Name:  "testname",
					Value: "testval",
				},
			},
			Commands: []string{"run"},
		}).Return(nil, er)

		res, err := s.CreateGenerators(ctx, &desc.CreateGeneratorsRequest{
			Parameters: []*desc.CreateGeneratorsParams{
				{
					Image: "testimage",
					Resources: &desc.Resources{
						Memory: &desc.Resource{
							Limit:   "5Gi",
							Request: "4Gi",
						},
						Cpu: &desc.Resource{
							Limit:   "4",
							Request: "3",
						},
					},
					AdditionalEnvs: []*desc.EnvVar{
						{
							Name: "testname",
							Val:  "testval",
						},
					},
					Commands: []string{"run"},
				},
			},
		})
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}
