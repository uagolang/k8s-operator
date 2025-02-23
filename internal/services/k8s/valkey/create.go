package valkey

import (
	"context"
	"encoding/base64"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/cmd/cli/utils"
	"github.com/uagolang/k8s-operator/internal/lib/validator"
	"github.com/uagolang/k8s-operator/internal/services/k8s"
)

type CreateRequest struct {
	CrdName   string            `json:"crd_name" validate:"required"`
	Namespace string            `json:"namespace" validate:"required"`
	Image     string            `json:"image" validate:"required"`
	User      string            `json:"user" validate:"required"`
	Password  string            `json:"password" validate:"required"`
	Replicas  int32             `json:"replicas" validate:"required"`
	Volume    v1alpha1.Volume   `json:"volume" validate:"required"`
	Resource  v1alpha1.Resource `json:"resource" validate:"required"`
}

func (s *valkeyService) Create(ctx context.Context, i *CreateRequest) (*v1alpha1.Valkey, error) {
	if err := validator.Validate(ctx, i); err != nil {
		return nil, err
	}

	_, err := s.createSecret(ctx, i)
	if err != nil {
		return nil, err
	}

	_, err = s.createDeployment(ctx, i)
	if err != nil {
		return nil, err
	}

	_, err = s.createService(ctx, i)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *valkeyService) createSecret(ctx context.Context, i *CreateRequest) (*corev1.Secret, error) {
	res, err := s.k8sClient.CoreV1().Secrets(i.Namespace).Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.CrdName,
			Namespace: i.Namespace,
		},
		StringData: map[string]string{
			k8s.SecretKeyPassword: base64.StdEncoding.EncodeToString([]byte(i.Password)),
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *valkeyService) createDeployment(ctx context.Context, i *CreateRequest) (*appsv1.Deployment, error) {
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume
	if i.Volume.Enabled {
		volumeMounts = []corev1.VolumeMount{{
			Name:      "data",
			MountPath: "/data",
		}}
		volumes = []corev1.Volume{{
			Name: "data",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "valkey-pvc",
				},
			},
		}}
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.CrdName,
			Namespace: i.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Pointer(i.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": i.CrdName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": i.CrdName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "valkey",
						Image: i.Image,
						Env: []corev1.EnvVar{
							{
								Name:  "VALKEY_USER",
								Value: i.User,
							},
							{
								Name: "VALKEY_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: i.CrdName,
										},
										Key: "password",
									},
								},
							},
						},
						Ports:        []corev1.ContainerPort{{ContainerPort: 6379}},
						VolumeMounts: volumeMounts,
					}},
					Volumes: volumes,
				},
			},
		},
	}
	res, err := s.k8sClient.AppsV1().Deployments(i.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *valkeyService) createService(ctx context.Context, i *CreateRequest) (*corev1.Service, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.CrdName,
			Namespace: i.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "valkey"},
			Ports: []corev1.ServicePort{
				{
					Port:       6379,
					TargetPort: intstr.FromInt32(6379),
				},
			},
			ClusterIP: "None",
		},
	}
	res, err := s.k8sClient.CoreV1().Services(i.Namespace).Create(ctx, service, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}
