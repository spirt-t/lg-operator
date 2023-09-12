package lg_operator

import (
	"context"

	"github.com/spirt-t/lg-operator/internal/k8s"
	"github.com/spirt-t/lg-operator/internal/model"
	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
	"golang.org/x/sync/errgroup"
)

// CreateGenerators ...
func (s *Service) CreateGenerators(ctx context.Context, in *desc.CreateGeneratorsRequest) (*desc.CreateGeneratorsResponse, error) {
	generators := make([]model.LoadGenerator, len(in.Parameters))

	g, ctxg := errgroup.WithContext(ctx)
	for i, inParams := range in.Parameters {
		params := inParams
		i := i
		g.Go(func() error {
			generator, err := s.createGenerator(ctxg, params)
			if err == nil {
				generators[i] = *generator
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &desc.CreateGeneratorsResponse{
		LoadGenerators: GeneratorMapper{}.ModelToPBMany(generators),
	}, nil
}

func (s *Service) createGenerator(ctx context.Context, in *desc.CreateGeneratorsParams) (*model.LoadGenerator, error) {
	resources, err := s.resourceMapper.PBToModel(in.Resources)
	if err != nil {
		return nil, err
	}

	envs := EnvVarMapper{}.PbToModelMany(in.AdditionalEnvs)

	generator, err := s.k8s.Create(ctx, k8s.CreationConfig{
		Image:            in.Image,
		Resources:        resources,
		Envs:             envs,
		Commands:         in.Commands,
		ExposeExternalIP: in.ExposeExternalIp,
	})

	return generator, err
}
