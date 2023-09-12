package lg_operator

import (
	"context"

	desc "github.com/spirt-t/lg-operator/pkg/lg-operator"
)

// Hello - just debug method.
func (s *Service) Hello(_ context.Context, _ *desc.HelloRequest) (*desc.HelloResponse, error) {
	return &desc.HelloResponse{Hello: "Hello!"}, nil
}
