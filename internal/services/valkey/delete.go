package valkey

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (s *valkeyService) Delete(ctx context.Context, i types.NamespacedName) error {
	err := s.deleteService(ctx, types.NamespacedName{
		Name:      i.Name,
		Namespace: i.Namespace,
	})
	if err != nil {
		return err
	}

	err = s.deleteDeployment(ctx, types.NamespacedName{
		Name:      i.Name,
		Namespace: i.Namespace,
	})
	if err != nil {
		return err
	}

	err = s.deleteSecret(ctx, types.NamespacedName{
		Name:      i.Name,
		Namespace: i.Namespace,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *valkeyService) deleteSecret(ctx context.Context, i types.NamespacedName) error {
	err := s.k8sClient.Delete(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.Name,
			Namespace: i.Namespace,
		},
	})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}

		return err
	}

	return nil
}

func (s *valkeyService) deleteDeployment(ctx context.Context, i types.NamespacedName) error {
	err := s.k8sClient.Delete(ctx, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.Name,
			Namespace: i.Namespace,
		},
	})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}

		return err
	}

	err = s.k8sClient.Delete(ctx, &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.Name,
			Namespace: i.Namespace,
		},
	})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}

		return err
	}

	return nil
}

func (s *valkeyService) deleteService(ctx context.Context, i types.NamespacedName) error {
	err := s.k8sClient.Delete(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.Name,
			Namespace: i.Namespace,
		},
	})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}

		return err
	}

	return nil
}
