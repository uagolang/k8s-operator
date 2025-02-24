package valkey

import (
	"context"
	"encoding/base64"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/internal/lib/validator"
	"github.com/uagolang/k8s-operator/internal/utils"
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
	res := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      i.CrdName,
			Namespace: i.Namespace,
		},
		StringData: map[string]string{
			secretKeyPassword: base64.StdEncoding.EncodeToString([]byte(i.Password)),
		},
	}
	if err := s.k8sClient.Create(ctx, res); err != nil && !k8serrors.IsAlreadyExists(err) {
		return nil, err
	}

	return res, s.waitForSecret(res.Name, res.Namespace, 5*time.Second)
}

func (s *valkeyService) waitForSecret(name, namespace string, dur time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	return wait.PollUntilContextTimeout(ctx, pollInterval, dur, true, func(ctx context.Context) (bool, error) {
		_, err := s.getSecret(ctx, types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		})
		if err == nil {
			return true, nil
		}
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	})
}

func (s *valkeyService) createDeployment(ctx context.Context, i *CreateRequest) (*appsv1.Deployment, error) {
	var volumeMounts []corev1.VolumeMount
	var volumes []corev1.Volume

	pvcName := "valkey-pvc"

	if i.Volume.Enabled {
		volumeMounts = []corev1.VolumeMount{{
			Name:      "data",
			MountPath: "/data",
		}}
		volumes = []corev1.Volume{{
			Name: "data",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: pvcName,
				},
			},
		}}
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: i.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(i.Resource.Storage),
				},
			},
		},
	}
	err := s.k8sClient.Create(ctx, pvc)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return nil, err
	}

	resourceList := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(i.Resource.CPU),
		corev1.ResourceMemory: resource.MustParse(i.Resource.Memory),
	}

	res := &appsv1.Deployment{
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
					Containers: []corev1.Container{
						{
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
							Resources: corev1.ResourceRequirements{
								Requests: resourceList,
								Limits:   resourceList,
							},
						},
					},
					Volumes: volumes,
				},
			},
		},
	}
	err = s.k8sClient.Create(ctx, res)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return nil, err
	}

	return res, nil
}

func (s *valkeyService) waitForDeployment(name, namespace string, dur time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	return wait.PollUntilContextTimeout(ctx, pollInterval, dur, true, func(ctx context.Context) (bool, error) {
		_, err := s.getDeployment(ctx, types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		})
		if err == nil {
			return true, nil
		}
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	})
}

func (s *valkeyService) createService(ctx context.Context, i *CreateRequest) (*corev1.Service, error) {
	res := &corev1.Service{
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
	err := s.k8sClient.Create(ctx, res)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return nil, err
	}

	return res, nil
}

func (s *valkeyService) waitForService(i types.NamespacedName, dur time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()

	return wait.PollUntilContextTimeout(ctx, pollInterval, dur, true, func(ctx context.Context) (bool, error) {
		_, err := s.getService(ctx, i)
		if err == nil {
			return true, nil
		}
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	})
}
