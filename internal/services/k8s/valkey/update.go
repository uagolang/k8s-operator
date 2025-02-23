package valkey

import (
	"context"
	"encoding/base64"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/internal/lib/validator"
	"github.com/uagolang/k8s-operator/internal/services/k8s"
)

type UpdateRequest struct {
	CrdName   string            `json:"crd_name" validate:"required"`
	Namespace string            `json:"namespace" validate:"required"`
	Image     string            `json:"image" validate:"omitempty"`
	User      string            `json:"user" validate:"omitempty"`
	Password  *string           `json:"password" validate:"omitempty"`
	Replicas  int32             `json:"replicas" validate:"omitempty"`
	Volume    v1alpha1.Volume   `json:"volume" validate:"required"`
	Resource  v1alpha1.Resource `json:"resource" validate:"required"`
}

func (s *valkeyService) Update(ctx context.Context, i *UpdateRequest) error {
	if err := validator.Validate(ctx, i); err != nil {
		return err
	}

	if i.Password != nil && *i.Password != "" {
		_, err := s.updateSecret(ctx, i)
		if err != nil {
			return err
		}
	}

	_, err := s.updateDeployment(ctx, i)
	if err != nil {
		return err
	}

	_, err = s.updateService(ctx, i)
	if err != nil {
		return err
	}

	return nil
}

func (s *valkeyService) updateSecret(ctx context.Context, i *UpdateRequest) (*corev1.Secret, error) {
	sec, err := s.k8sClient.CoreV1().Secrets(i.Namespace).Get(ctx, i.CrdName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	sec.StringData[k8s.SecretKeyPassword] = base64.StdEncoding.EncodeToString([]byte(*i.Password))

	res, err := s.k8sClient.CoreV1().Secrets(sec.Namespace).Update(ctx, sec, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *valkeyService) updateDeployment(ctx context.Context, i *UpdateRequest) (*appsv1.Deployment, error) {
	deployment, err := s.k8sClient.AppsV1().Deployments(i.Namespace).Get(ctx, i.CrdName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	res, err := s.k8sClient.AppsV1().Deployments(i.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *valkeyService) updateService(ctx context.Context, i *UpdateRequest) (*corev1.Service, error) {
	service, err := s.k8sClient.CoreV1().Services(i.Namespace).Get(ctx, i.CrdName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	res, err := s.k8sClient.CoreV1().Services(service.Namespace).Update(ctx, service, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
