package valkey

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
)

type Service interface {
	Create(ctx context.Context, i *CreateRequest) (*v1alpha1.Valkey, error)
	Update(ctx context.Context, i *UpdateRequest) error
	IsReady(ctx context.Context, item *v1alpha1.Valkey) (bool, int32, error)
	Delete(ctx context.Context, i types.NamespacedName) error
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
