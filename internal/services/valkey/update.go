package valkey

import (
	"context"
	"encoding/base64"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/internal/lib/validator"
)

type UpdateRequest struct {
	CrdName   string             `json:"crd_name" validate:"required"`
	Namespace string             `json:"namespace" validate:"required"`
	Image     *string            `json:"image,omitempty" validate:"omitempty"`
	User      *string            `json:"user,omitempty" validate:"omitempty"`
	Password  *string            `json:"password,omitempty" validate:"omitempty"`
	Replicas  *int32             `json:"replicas,omitempty" validate:"omitempty"`
	Volume    *v1alpha1.Volume   `json:"volume,omitempty" validate:"omitempty"`
	Resource  *v1alpha1.Resource `json:"resource,omitempty" validate:"omitempty"`
}

func (s *valkeyService) Update(ctx context.Context, i *UpdateRequest) error {
	if err := validator.Validate(ctx, i); err != nil {
		return err
	}

	if i.Password != nil && *i.Password != "" {
		err := s.updateSecret(ctx, i)
		if err != nil {
			return err
		}
	}

	err := s.updateDeployment(ctx, i)
	if err != nil {
		return err
	}

	err = s.updateService(ctx, i)
	if err != nil {
		return err
	}

	return nil
}

func (s *valkeyService) getSecret(ctx context.Context, i types.NamespacedName) (*corev1.Secret, error) {
	res := new(corev1.Secret)
	err := s.k8sClient.Get(ctx, types.NamespacedName{
		Name:      i.Name,
		Namespace: i.Namespace,
	}, res)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}

func (s *valkeyService) updateSecret(ctx context.Context, i *UpdateRequest) error {
	res, err := s.getSecret(ctx, types.NamespacedName{
		Name:      i.CrdName,
		Namespace: i.Namespace,
	})
	if err != nil {
		return err
	}
	if res == nil {
		return nil
	}

	res.Data[secretKeyPassword] = []byte(base64.StdEncoding.EncodeToString([]byte(*i.Password)))

	err = s.k8sClient.Update(ctx, res)
	if err != nil {
		return err
	}

	return nil
}

func (s *valkeyService) getDeployment(ctx context.Context, i types.NamespacedName) (*appsv1.Deployment, error) {
	res := new(appsv1.Deployment)
	err := s.k8sClient.Get(ctx, i, res)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}

func (s *valkeyService) updateDeployment(ctx context.Context, i *UpdateRequest) error {
	res, err := s.getDeployment(ctx, types.NamespacedName{
		Name:      i.CrdName,
		Namespace: i.Namespace,
	})
	if err != nil {
		return err
	}
	if res == nil {
		return nil
	}

	var shouldUpdate bool
	if res.Spec.Replicas != nil && i.Replicas != nil && *res.Spec.Replicas != *i.Replicas {
		shouldUpdate = true
		res.Spec.Replicas = i.Replicas
	}

	if len(res.Spec.Template.Spec.Containers) > 0 {
		container := &res.Spec.Template.Spec.Containers[0]
		if i.Image != nil {
			shouldUpdate = true
			container.Image = *i.Image
		}

		if i.Resource != nil {
			shouldUpdate = true

			resourceList := corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse(i.Resource.CPU),
				corev1.ResourceMemory: resource.MustParse(i.Resource.Memory),
			}

			container.Resources = corev1.ResourceRequirements{
				Requests: resourceList,
				Limits:   resourceList,
			}
		}
	}

	if shouldUpdate {
		err = s.k8sClient.Update(ctx, res)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *valkeyService) updateService(ctx context.Context, i *UpdateRequest) error {
	res := new(corev1.Service)
	err := s.k8sClient.Get(ctx, types.NamespacedName{
		Name:      i.CrdName,
		Namespace: i.Namespace,
	}, res)
	if err != nil {
		return err
	}

	err = s.k8sClient.Update(ctx, res)
	if err != nil {
		return err
	}

	return nil
}
