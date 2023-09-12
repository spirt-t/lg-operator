package lg_operator

import (
	coreV1 "k8s.io/api/core/v1"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spirt-t/lg-operator/internal/config"
	"github.com/spirt-t/lg-operator/internal/model"
	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
	"github.com/stretchr/testify/assert"
)

func TestResourceMapper_PBToModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mngr, err := config.NewManager("../../../../testfiles/test_config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	mapper := NewResourceMapper(mngr)

	t.Run("ok", func(t *testing.T) {
		res, err := mapper.PBToModel(&desc.Resources{
			Memory: &desc.Resource{
				Limit:   "5Gi",
				Request: "4Gi",
			},
			Cpu: &desc.Resource{
				Limit:   "4",
				Request: "3",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, model.Resources{
			Memory: model.Resource{
				Limit:   "5Gi",
				Request: "4Gi",
			},
			CPU: model.Resource{
				Limit:   "4",
				Request: "3",
			},
		}, res)
	})

	t.Run("nil resource", func(t *testing.T) {
		res, err := mapper.PBToModel(nil)
		assert.NoError(t, err)
		assert.Equal(t, model.Resources{
			Memory: model.Resource{
				Limit:   "2Gi",
				Request: "1Gi",
			},
			CPU: model.Resource{
				Limit:   "2",
				Request: "1",
			},
		}, res)
	})

	t.Run("memory only", func(t *testing.T) {
		res, err := mapper.PBToModel(&desc.Resources{
			Memory: &desc.Resource{
				Limit:   "5Gi",
				Request: "4Gi",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, model.Resources{
			Memory: model.Resource{
				Limit:   "5Gi",
				Request: "4Gi",
			},
			CPU: model.Resource{
				Limit:   "2",
				Request: "1",
			},
		}, res)
	})

	t.Run("cpu only", func(t *testing.T) {
		res, err := mapper.PBToModel(&desc.Resources{
			Cpu: &desc.Resource{
				Limit:   "4",
				Request: "3",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, model.Resources{
			Memory: model.Resource{
				Limit:   "2Gi",
				Request: "1Gi",
			},
			CPU: model.Resource{
				Limit:   "4",
				Request: "3",
			},
		}, res)
	})

	t.Run("cpu limit only", func(t *testing.T) {
		res, err := mapper.PBToModel(&desc.Resources{
			Cpu: &desc.Resource{
				Limit: "4",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, model.Resources{
			Memory: model.Resource{
				Limit:   "2Gi",
				Request: "1Gi",
			},
			CPU: model.Resource{
				Limit:   "4",
				Request: "1",
			},
		}, res)
	})

	t.Run("cpu request only", func(t *testing.T) {
		res, err := mapper.PBToModel(&desc.Resources{
			Cpu: &desc.Resource{
				Request: "3",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, model.Resources{
			Memory: model.Resource{
				Limit:   "2Gi",
				Request: "1Gi",
			},
			CPU: model.Resource{
				Limit:   "2",
				Request: "3",
			},
		}, res)
	})

	t.Run("memory limit only", func(t *testing.T) {
		res, err := mapper.PBToModel(&desc.Resources{
			Memory: &desc.Resource{
				Limit: "5Gi",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, model.Resources{
			Memory: model.Resource{
				Limit:   "5Gi",
				Request: "1Gi",
			},
			CPU: model.Resource{
				Limit:   "2",
				Request: "1",
			},
		}, res)
	})

	t.Run("memory request only", func(t *testing.T) {
		res, err := mapper.PBToModel(&desc.Resources{
			Memory: &desc.Resource{
				Request: "4Gi",
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, model.Resources{
			Memory: model.Resource{
				Limit:   "2Gi",
				Request: "4Gi",
			},
			CPU: model.Resource{
				Limit:   "2",
				Request: "1",
			},
		}, res)
	})
}

func TestGeneratorMapper(t *testing.T) {
	mapper := GeneratorMapper{}

	t.Run("ModelToPB full info", func(t *testing.T) {
		lg := model.LoadGenerator{
			Name:       "testname",
			ClusterIP:  "127.0.0.1",
			ExternalIP: "0.0.0.1",
			Port:       8888,
			Status:     coreV1.PodRunning,
		}
		res := mapper.ModelToPB(lg)
		assert.NotNil(t, res)
		assert.Equal(t, lg.ExternalIP, res.ExternalIp)
		assert.Equal(t, lg.ClusterIP, res.ClusterIp)
		assert.Equal(t, lg.Name, res.Name)
		assert.Equal(t, lg.Port, res.Port)
		assert.Equal(t, string(lg.Status), res.Status)
	})

	t.Run("ModelToPB part info", func(t *testing.T) {
		lg := model.LoadGenerator{
			Name:      "testname",
			ClusterIP: "127.0.0.1",
			Status:    coreV1.PodRunning,
		}
		res := mapper.ModelToPB(lg)
		assert.NotNil(t, res)
		assert.Equal(t, lg.ClusterIP, res.ClusterIp)
		assert.Equal(t, lg.Name, res.Name)
		assert.Equal(t, lg.Port, res.Port)
		assert.Equal(t, string(lg.Status), res.Status)
	})

	t.Run("ModelToPBMany full info", func(t *testing.T) {
		lgs := []model.LoadGenerator{
			{
				Name:       "testname-1",
				ClusterIP:  "127.0.0.1",
				ExternalIP: "0.0.0.1",
				Port:       8888,
				Status:     coreV1.PodRunning,
			},
			{
				Name:       "testname-2",
				ClusterIP:  "127.0.0.2",
				ExternalIP: "0.0.0.2",
				Port:       8888,
				Status:     coreV1.PodRunning,
			},
		}
		res := mapper.ModelToPBMany(lgs)
		assert.NotNil(t, res)

		for i, resLg := range res {
			assert.Equal(t, lgs[i].ClusterIP, resLg.ClusterIp)
			assert.Equal(t, lgs[i].ExternalIP, resLg.ExternalIp)
			assert.Equal(t, lgs[i].Name, resLg.Name)
			assert.Equal(t, lgs[i].Port, resLg.Port)
			assert.Equal(t, string(lgs[i].Status), resLg.Status)
		}
	})
}

func TestEnvVarMapper(t *testing.T) {
	mapper := EnvVarMapper{}

	t.Run("PBToModel ok", func(t *testing.T) {
		res := mapper.PBToModel(&desc.EnvVar{
			Name: "testname",
			Val:  "testval",
		})
		assert.Equal(t, model.EnvVar{
			Name:  "testname",
			Value: "testval",
		}, *res)
	})

	t.Run("PBToModel nil", func(t *testing.T) {
		res := mapper.PBToModel(nil)
		assert.Equal(t, (*model.EnvVar)(nil), res)
	})

	t.Run("PbToModelMany ok", func(t *testing.T) {
		envs := []*desc.EnvVar{
			{
				Name: "testname-1",
				Val:  "testval-1",
			},
			{
				Name: "testname-2",
				Val:  "testval-2",
			},
		}
		res := mapper.PbToModelMany(envs)
		assert.Equal(t, 2, len(res))
		for i, env := range res {
			assert.Equal(t, model.EnvVar{
				Name:  envs[i].Name,
				Value: envs[i].Val,
			}, env)
		}
	})

	t.Run("PbToModelMany nil", func(t *testing.T) {
		res := mapper.PbToModelMany(nil)
		assert.Equal(t, 0, len(res))
	})
}
