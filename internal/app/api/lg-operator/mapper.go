package lg_operator

import (
	"github.com/spirt-t/lg-operator/internal/config"
	"github.com/spirt-t/lg-operator/internal/model"
	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
)

const (
	defaultResourcesKey = "default_resources"
)

// ResourceMapper ...
type ResourceMapper struct {
	cfg config.Manager
}

// NewResourceMapper - constructor for ResourceMapper.
func NewResourceMapper(cfg config.Manager) *ResourceMapper {
	return &ResourceMapper{
		cfg: cfg,
	}
}

// PBToModel ...
func (m *ResourceMapper) PBToModel(resources *desc.Resources) (model.Resources, error) {
	var r model.Resources

	if err := m.cfg.UnmarshalKey(defaultResourcesKey, &r); err != nil {
		return model.Resources{}, err
	}

	if resources != nil {
		if resources.Cpu != nil {
			if resources.Cpu.Request != "" {
				r.CPU.Request = resources.Cpu.Request
			}

			if resources.Cpu.Limit != "" {
				r.CPU.Limit = resources.Cpu.Limit
			}
		}

		if resources.Memory != nil {
			if resources.Memory.Request != "" {
				r.Memory.Request = resources.Memory.Request
			}

			if resources.Memory.Limit != "" {
				r.Memory.Limit = resources.Memory.Limit
			}
		}
	}

	return r, nil
}

// GeneratorMapper ...
type GeneratorMapper struct{}

// ModelToPB - map generator model to proto-message.
func (gm GeneratorMapper) ModelToPB(generator model.LoadGenerator) *desc.LoadGenerator {
	return &desc.LoadGenerator{
		Name:       generator.Name,
		ClusterIp:  generator.ClusterIP,
		ExternalIp: generator.ExternalIP,
		Port:       generator.Port,
		Status:     string(generator.Status),
	}
}

// ModelToPBMany - map generators to proto-message.
func (gm GeneratorMapper) ModelToPBMany(generators []model.LoadGenerator) []*desc.LoadGenerator {
	list := make([]*desc.LoadGenerator, 0, len(generators))
	for _, generator := range generators {
		list = append(list, gm.ModelToPB(generator))
	}

	return list
}

// EnvVarMapper ...
type EnvVarMapper struct{}

// PBToModel - map environment variable maode to proto-message.
func (em EnvVarMapper) PBToModel(envVar *desc.EnvVar) *model.EnvVar {
	if envVar == nil {
		return nil
	}

	return &model.EnvVar{
		Name:  envVar.Name,
		Value: envVar.Val,
	}
}

// PbToModelMany ...
func (em EnvVarMapper) PbToModelMany(envVars []*desc.EnvVar) []model.EnvVar {
	listEnvVars := make([]model.EnvVar, 0, len(envVars))

	for _, envVar := range envVars {
		if newVar := em.PBToModel(envVar); newVar != nil {
			listEnvVars = append(listEnvVars, *newVar)
		}
	}

	return listEnvVars
}
