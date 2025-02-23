package valkey

import (
	"context"
	"fmt"
)

func (s *valkeyService) Delete(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	return nil
}

//
//func (s *valkeyService) updateSecret(ctx context.Context, i *UpdateRequest) (*corev1.Secret, error) {
//	sec, err := s.k8sClient.CoreV1().Secrets(i.Namespace).Get(ctx, i.CrdName, metav1.GetOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	sec.StringData[k8s.SecretKeyPassword] = base64.StdEncoding.EncodeToString([]byte(*i.Password))
//
//	res, err := s.k8sClient.CoreV1().Secrets(sec.Namespace).Update(ctx, sec, metav1.UpdateOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	return res, nil
//}
//
//func (s *valkeyService) updateDeployment(ctx context.Context, i *UpdateRequest) (*appsv1.Deployment, error) {
//	deployment, err := s.k8sClient.AppsV1().Deployments(i.Namespace).Get(ctx, i.CrdName, metav1.GetOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	res, err := s.k8sClient.AppsV1().Deployments(i.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	return res, nil
//}
//
//func (s *valkeyService) updateService(ctx context.Context, i *UpdateRequest) (*corev1.Service, error) {
//	service, err := s.k8sClient.CoreV1().Services(i.Namespace).Get(ctx, i.CrdName, metav1.GetOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	res, err := s.k8sClient.CoreV1().Services(service.Namespace).Update(ctx, service, metav1.UpdateOptions{})
//	if err != nil {
//		return nil, err
//	}
//
//	return res, nil
//}
