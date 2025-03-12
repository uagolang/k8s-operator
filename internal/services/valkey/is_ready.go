package valkey

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/uagolang/k8s-operator/internal/lib/validator"
)

type IsReadyRequest struct {
	Name      string `json:"name" validate:"required"`
	Namespace string `json:"namespace" validate:"required"`
}

func (s *valkeyService) IsReady(ctx context.Context, i *IsReadyRequest) (bool, int32, error) {
	if err := validator.Validate(ctx, i); err != nil {
		return false, 0, err
	}

	deployment := &appsv1.Deployment{}
	err := s.k8sClient.Get(ctx, types.NamespacedName{
		Name:      i.Name,
		Namespace: i.Namespace,
	}, deployment)
	if err != nil {
		return false, 0, err
	}

	if deployment.Status.ReadyReplicas > 0 {
		return true, deployment.Status.ReadyReplicas, nil
	}

	return false, 0, nil
}
