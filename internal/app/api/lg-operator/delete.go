package lg_operator

import (
	"context"

	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
	"go.uber.org/multierr"
)

// DeleteGenerators - delete generator pod, service and ingress by generator name.
func (s *Service) DeleteGenerators(ctx context.Context, in *desc.DeleteGeneratorsRequest) (*desc.DeleteGeneratorsResponse, error) {
	n := len(in.Names)
	errs := make(chan error, len(in.Names))
	defer close(errs)

	for _, name := range in.Names {
		go func(name string) {
			errs <- s.k8s.Delete(ctx, name)
		}(name)
	}

	var err error
	for i := 0; i < n; i++ {
		err = multierr.Append(err, <-errs)
	}

	return &desc.DeleteGeneratorsResponse{}, err
}
