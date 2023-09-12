package lg_operator

import (
	"context"

	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
)

// GeneratorsList - list of launched generators.
func (s *Service) GeneratorsList(ctx context.Context, _ *desc.GeneratorsListRequest) (*desc.GeneratorsListResponse, error) {
	generators, err := s.k8s.List(ctx)
	if err != nil {
		return nil, err
	}

	list := GeneratorMapper{}.ModelToPBMany(generators)

	return &desc.GeneratorsListResponse{LoadGenerators: list}, nil
}
