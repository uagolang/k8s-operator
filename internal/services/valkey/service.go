package valkey

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Service interface {
	Create(ctx context.Context, i *CreateRequest) error
	Update(ctx context.Context, i *UpdateRequest) error
	IsReady(ctx context.Context, i *IsReadyRequest) (bool, int32, error)
	Delete(ctx context.Context, i *DeleteRequest) error
}

type valkeyService struct {
	k8sClient client.Client
}

type Option func(s *valkeyService)

func WithK8sClient(v client.Client) Option {
	return func(s *valkeyService) {
		s.k8sClient = v
	}
}

func NewValkeyService(opts ...Option) Service {
	s := new(valkeyService)
	for _, opt := range opts {
		opt(s)
	}

	return s
}
