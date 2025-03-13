package valkey_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
	validatorlib "github.com/uagolang/k8s-operator/internal/lib/validator"
	"github.com/uagolang/k8s-operator/internal/services/valkey"
	"github.com/uagolang/k8s-operator/internal/utils"
	"github.com/uagolang/k8s-operator/mocks"
)

func TestService(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockErr := errors.New("mock error")
	storage := "300Mi"
	k8sClient := mocks.NewMockK8sClient(ctrl)
	s := valkey.NewValkeyService(valkey.WithK8sClient(k8sClient))

	createRequest := &valkey.CreateRequest{
		CrdName:   "valkey",
		Namespace: "default",
		Image:     "nesymno/k8s-operator:latest",
		User:      "user",
		Password:  "password",
		Replicas:  1,
		Volume: v1alpha1.Volume{
			Enabled: true,
			Storage: storage,
		},
		Resource: v1alpha1.Resource{
			CPU:     "100m",
			Memory:  "200Mi",
			Storage: storage,
		},
	}

	isReadyRequest := &valkey.IsReadyRequest{
		Name:      createRequest.CrdName,
		Namespace: createRequest.Namespace,
	}

	updateRequest := &valkey.UpdateRequest{
		CrdName:   createRequest.CrdName,
		Namespace: createRequest.Namespace,
	}

	deleteRequest := &valkey.DeleteRequest{
		Name:      createRequest.CrdName,
		Namespace: createRequest.Namespace,
	}

	t.Run("create", func(t *testing.T) {

		t.Run("success", func(t *testing.T) {
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			err := s.Create(ctx, createRequest)
			require.NoError(t, err)
		})

		t.Run("with validation errors", func(t *testing.T) {
			// make copy of object, don't use pointer
			// because value will be updated in createRequest
			req := *createRequest
			// invalidate some request data
			req.CrdName = ""
			req.Namespace = ""
			req.Image = ""

			// pointer to copied object
			err := s.Create(ctx, &req)
			require.Error(t, err)

			// get validator errors
			errs := validatorlib.GetErrors(err)
			if len(errs) == 0 {
				require.Errorf(t, err, "expected validation errors")
			}

			require.Len(t, errs, 3)
			require.Equal(t, validatorlib.Errors{
				{
					Field:   "crd_name",
					Message: "crd_name is a required field",
				},
				{
					Field:   "namespace",
					Message: "namespace is a required field",
				},
				{
					Field:   "image",
					Message: "image is a required field",
				},
			}, errs)
		})

		t.Run("create secret failed", func(t *testing.T) {
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Create(ctx, createRequest)
			require.Error(t, err)
		})

		t.Run("wait secret failed", func(t *testing.T) {
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			notFoundErr := k8serrors.NewNotFound(schema.GroupResource{
				Group:    "",
				Resource: "secrets",
			}, createRequest.CrdName)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(notFoundErr)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Create(ctx, createRequest)
			require.Error(t, err)
		})

		t.Run("create deployment failed", func(t *testing.T) {
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Create(ctx, createRequest)
			require.Error(t, err)
		})

		t.Run("create deployment failed", func(t *testing.T) {
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Create(ctx, createRequest)
			require.Error(t, err)
		})

		t.Run("wait deployment failed", func(t *testing.T) {
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			notFoundErr := k8serrors.NewNotFound(schema.GroupResource{
				Group:    "",
				Resource: "deployments",
			}, createRequest.CrdName)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(notFoundErr)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Create(ctx, createRequest)
			require.Error(t, err)
		})

		t.Run("create service failed", func(t *testing.T) {
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Create(ctx, createRequest)
			require.Error(t, err)
		})

		t.Run("wait service failed", func(t *testing.T) {
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			notFoundErr := k8serrors.NewNotFound(schema.GroupResource{
				Group:    "",
				Resource: "services",
			}, createRequest.CrdName)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(notFoundErr)
			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Create(ctx, createRequest)
			require.Error(t, err)
		})
	})

	t.Run("is_ready", func(t *testing.T) {

		t.Run("success", func(t *testing.T) {
			deployment := &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{
					ReadyReplicas: 1,
				},
			}

			k8sClient.EXPECT().Get(ctx, types.NamespacedName{
				Name:      createRequest.CrdName,
				Namespace: createRequest.Namespace,
			}, gomock.AssignableToTypeOf(deployment)).DoAndReturn(
				func(_ context.Context, _ types.NamespacedName, obj runtimeclient.Object, _ ...runtimeclient.GetOption) error {
					*obj.(*appsv1.Deployment) = *(deployment)
					return nil
				})
		})

		ready, readyReplicas, err := s.IsReady(ctx, isReadyRequest)
		require.NoError(t, err)
		require.True(t, ready)
		require.Equal(t, int32(1), readyReplicas)
	})

	t.Run("with validation errors", func(t *testing.T) {
		req := *isReadyRequest
		req.Name = ""
		req.Namespace = ""

		ready, readyReplicas, err := s.IsReady(ctx, &req)
		require.Error(t, err)
		require.False(t, ready)
		require.Equal(t, int32(0), readyReplicas)

		// get validator errors
		errs := validatorlib.GetErrors(err)
		if len(errs) == 0 {
			require.Errorf(t, err, "expected validation errors")
		}

		require.Len(t, errs, 2)
		require.Equal(t, validatorlib.Errors{
			{
				Field:   "name",
				Message: "name is a required field",
			},
			{
				Field:   "namespace",
				Message: "namespace is a required field",
			},
		}, errs)
	})

	t.Run("get crd failed", func(t *testing.T) {
		deployment := &appsv1.Deployment{
			Status: appsv1.DeploymentStatus{
				ReadyReplicas: 0,
			},
		}

		k8sClient.EXPECT().Get(ctx, types.NamespacedName{
			Name:      createRequest.CrdName,
			Namespace: createRequest.Namespace,
		}, gomock.AssignableToTypeOf(deployment)).DoAndReturn(
			func(_ context.Context, _ types.NamespacedName, obj runtimeclient.Object, _ ...runtimeclient.GetOption) error {
				*obj.(*appsv1.Deployment) = *(deployment)
				return nil
			})
		k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

		ready, readyReplicas, err := s.IsReady(ctx, isReadyRequest)
		require.NoError(t, err)
		require.False(t, ready)
		require.Equal(t, int32(0), readyReplicas)

		ready, readyReplicas, err = s.IsReady(ctx, isReadyRequest)
		require.Error(t, err)
		require.False(t, ready)
		require.Equal(t, int32(0), readyReplicas)
	})

	t.Run("update", func(t *testing.T) {
		secret := &v1.Secret{}
		deployment := &appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name:  "valkey",
								Image: "myimage",
							},
						},
					},
				},
				Replicas: utils.Pointer(int32(1)),
			},
			Status: appsv1.DeploymentStatus{
				ReadyReplicas: 2,
			},
		}
		service := &v1.Service{
			Status: v1.ServiceStatus{},
		}

		t.Run("success", func(t *testing.T) {
			k8sClient.EXPECT().Get(ctx, types.NamespacedName{
				Name:      createRequest.CrdName,
				Namespace: createRequest.Namespace,
			}, gomock.AssignableToTypeOf(secret)).DoAndReturn(
				func(_ context.Context, _ types.NamespacedName, obj runtimeclient.Object, _ ...runtimeclient.GetOption) error {
					*obj.(*v1.Secret) = *(secret)
					return nil
				})
			k8sClient.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			k8sClient.EXPECT().Get(ctx, types.NamespacedName{
				Name:      createRequest.CrdName,
				Namespace: createRequest.Namespace,
			}, gomock.AssignableToTypeOf(deployment)).DoAndReturn(
				func(_ context.Context, _ types.NamespacedName, obj runtimeclient.Object, _ ...runtimeclient.GetOption) error {
					*obj.(*appsv1.Deployment) = *(deployment)
					return nil
				})
			k8sClient.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			k8sClient.EXPECT().Get(ctx, types.NamespacedName{
				Name:      createRequest.CrdName,
				Namespace: createRequest.Namespace,
			}, gomock.AssignableToTypeOf(service)).DoAndReturn(
				func(_ context.Context, _ types.NamespacedName, obj runtimeclient.Object, _ ...runtimeclient.GetOption) error {
					*obj.(*v1.Service) = *(service)
					return nil
				})
			k8sClient.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			err := s.Update(ctx, &valkey.UpdateRequest{
				CrdName:   createRequest.CrdName,
				Namespace: createRequest.Namespace,
				Image:     &createRequest.Image,
				Password:  &createRequest.Password,
				Replicas:  utils.Pointer(int32(2)),
				Resource: &v1alpha1.Resource{
					CPU:     "50m",
					Memory:  "64Mi",
					Storage: "0",
				},
			})
			require.NoError(t, err)
		})

		t.Run("with validation errors", func(t *testing.T) {
			req := *updateRequest
			req.CrdName = ""
			req.Namespace = ""

			err := s.Update(ctx, &req)
			require.Error(t, err)

			// get validator errors
			errs := validatorlib.GetErrors(err)
			if len(errs) == 0 {
				require.Errorf(t, err, "expected validation errors")
			}

			require.Len(t, errs, 2)
			require.Equal(t, validatorlib.Errors{
				{
					Field:   "crd_name",
					Message: "crd_name is a required field",
				},
				{
					Field:   "namespace",
					Message: "namespace is a required field",
				},
			}, errs)
		})

		t.Run("update secret failed", func(t *testing.T) {
			req := *updateRequest
			req.Password = utils.Pointer("password2")
			req.Image = utils.Pointer("myimage")
			req.Replicas = utils.Pointer(int32(3))

			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Update(ctx, &req)
			require.Error(t, err)
		})

		t.Run("update deployment failed", func(t *testing.T) {
			req := *updateRequest
			// do not update Password to skip update secret
			req.Image = utils.Pointer("myimage")
			req.Replicas = utils.Pointer(int32(2))

			k8sClient.EXPECT().Get(ctx, types.NamespacedName{
				Name:      createRequest.CrdName,
				Namespace: createRequest.Namespace,
			}, gomock.AssignableToTypeOf(deployment)).DoAndReturn(
				func(_ context.Context, _ types.NamespacedName, obj runtimeclient.Object, _ ...runtimeclient.GetOption) error {
					*obj.(*appsv1.Deployment) = *(deployment)
					return nil
				})
			k8sClient.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Update(ctx, &req)
			require.Error(t, err)
		})

		t.Run("update service failed", func(t *testing.T) {
			req := *updateRequest
			// do not update Password to skip update secret
			req.Image = utils.Pointer("myimage")
			req.Replicas = utils.Pointer(int32(2))

			k8sClient.EXPECT().Get(ctx, types.NamespacedName{
				Name:      createRequest.CrdName,
				Namespace: createRequest.Namespace,
			}, gomock.AssignableToTypeOf(deployment)).DoAndReturn(
				func(_ context.Context, _ types.NamespacedName, obj runtimeclient.Object, _ ...runtimeclient.GetOption) error {
					*obj.(*appsv1.Deployment) = *(deployment)
					return nil
				})
			k8sClient.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			k8sClient.EXPECT().Get(ctx, types.NamespacedName{
				Name:      createRequest.CrdName,
				Namespace: createRequest.Namespace,
			}, gomock.AssignableToTypeOf(service)).DoAndReturn(
				func(_ context.Context, _ types.NamespacedName, obj runtimeclient.Object, _ ...runtimeclient.GetOption) error {
					*obj.(*v1.Service) = *(service)
					return nil
				})
			k8sClient.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Update(ctx, &req)
			require.Error(t, err)
		})

		t.Run("not found resources", func(t *testing.T) {
			req := *updateRequest
			req.Password = utils.Pointer("password2")

			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(k8serrors.NewNotFound(schema.GroupResource{
					Group:    "",
					Resource: "secrets",
				}, updateRequest.CrdName))

			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(k8serrors.NewNotFound(schema.GroupResource{
					Group:    "",
					Resource: "deployments",
				}, updateRequest.CrdName))

			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(k8serrors.NewNotFound(schema.GroupResource{
					Group:    "",
					Resource: "services",
				}, updateRequest.CrdName))

			err := s.Update(ctx, &req)
			require.NoError(t, err)
		})

		t.Run("get secret failed", func(t *testing.T) {
			req := *updateRequest
			req.Password = utils.Pointer("password2")

			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Update(ctx, &req)
			require.Error(t, err)
		})

		t.Run("get deployment failed", func(t *testing.T) {
			req := *updateRequest

			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Update(ctx, &req)
			require.Error(t, err)
		})

		t.Run("get service failed", func(t *testing.T) {
			req := *updateRequest

			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(k8serrors.NewNotFound(schema.GroupResource{
					Group:    "",
					Resource: "deployments",
				}, updateRequest.CrdName))

			k8sClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Update(ctx, &req)
			require.Error(t, err)
		})

	})

	t.Run("delete", func(t *testing.T) {

		t.Run("success", func(t *testing.T) {
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			err := s.Delete(ctx, deleteRequest)
			require.NoError(t, err)
		})

		t.Run("with validation errors", func(t *testing.T) {
			req := *deleteRequest
			req.Name = ""
			req.Namespace = ""

			err := s.Delete(ctx, &req)
			require.Error(t, err)

			// get validator errors
			errs := validatorlib.GetErrors(err)
			if len(errs) == 0 {
				require.Errorf(t, err, "expected validation errors")
			}

			require.Len(t, errs, 2)
			require.Equal(t, validatorlib.Errors{
				{
					Field:   "name",
					Message: "name is a required field",
				},
				{
					Field:   "namespace",
					Message: "namespace is a required field",
				},
			}, errs)
		})

		t.Run("delete secret failed", func(t *testing.T) {
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Delete(ctx, deleteRequest)
			require.Error(t, err)
		})

		t.Run("delete deployment failed", func(t *testing.T) {
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Delete(ctx, deleteRequest)
			require.Error(t, err)
		})

		t.Run("delete pvc failed", func(t *testing.T) {
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Delete(ctx, deleteRequest)
			require.Error(t, err)
		})

		t.Run("delete service failed", func(t *testing.T) {
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockErr)

			err := s.Delete(ctx, deleteRequest)
			require.Error(t, err)
		})

		t.Run("idempotent delete", func(t *testing.T) {
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(k8serrors.NewNotFound(schema.GroupResource{
				Group:    "",
				Resource: "secrets",
			}, createRequest.CrdName))
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(k8serrors.NewNotFound(schema.GroupResource{
				Group:    "",
				Resource: "deployments",
			}, createRequest.CrdName))
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(k8serrors.NewNotFound(schema.GroupResource{
				Group:    "",
				Resource: "persistentvolumeclaims",
			}, createRequest.CrdName))
			k8sClient.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(k8serrors.NewNotFound(schema.GroupResource{
				Group:    "",
				Resource: "services",
			}, createRequest.CrdName))

			err := s.Delete(ctx, deleteRequest)
			require.NoError(t, err)
		})

	})
}
