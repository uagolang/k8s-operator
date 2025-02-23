package valkey

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
)

func (s *valkeyService) IsReady(ctx context.Context, item *v1alpha1.Valkey) (bool, int32, error) {
	deployment := &appsv1.Deployment{}
	err := s.k8sClient.Get(ctx, types.NamespacedName{
		Name:      item.Name,
		Namespace: item.Namespace,
	}, deployment)
	if err != nil {
		return false, 0, err
	}

	if deployment.Status.ReadyReplicas > 0 {
		return true, deployment.Status.ReadyReplicas, nil
	}

	return false, 0, nil
}
