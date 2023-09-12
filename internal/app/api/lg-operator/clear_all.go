package lg_operator

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

// ClearAll - delete all generator's pods, services and ingresses.
func (s *Service) ClearAll(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.k8s.DeleteAll(ctx)

	return &emptypb.Empty{}, err
}
