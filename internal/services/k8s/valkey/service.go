package valkey

import (
	"context"

	"k8s.io/client-go/kubernetes"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
)

type Service interface {
	Create(ctx context.Context, i *CreateRequest) (*v1alpha1.Valkey, error)
	Update(ctx context.Context, i *UpdateRequest) error
	Delete(ctx context.Context, name string) error
}

type valkeyService struct {
	k8sClient *kubernetes.Clientset
}

type Option func(s *valkeyService)

func WithK8sClient(c *kubernetes.Clientset) Option {
	return func(s *valkeyService) {
		s.k8sClient = c
	}
}

func NewPostgresService(opts ...Option) Service {
	s := new(valkeyService)
	for _, opt := range opts {
		opt(s)
	}

	return s
}
