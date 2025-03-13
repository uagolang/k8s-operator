package valkey_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	databasev1alpha1 "github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/internal/controller/flows"
	"github.com/uagolang/k8s-operator/internal/controller/flows/valkey"
	"github.com/uagolang/k8s-operator/mocks"
)

func TestFlowRun(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const (
		resourceName     = "test-resource"
		defaultNamespace = "default"
	)

	mockErr := errors.New("mock error")
	mockK8sClient := mocks.NewMockK8sClient(ctrl)
	mockValkeySvc := mocks.NewMockValkeyService(ctrl)

	flow := valkey.NewFlow(
		valkey.WithK8sClient(mockK8sClient),
		valkey.WithValkeySvc(mockValkeySvc),
	)

	t.Run("create resources error", func(t *testing.T) {
		mockValkeySvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(mockErr)

		status, finalizers, err := flow.Run(ctx, databasev1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: defaultNamespace,
			},
		})
		require.Nil(t, status)
		require.Nil(t, finalizers)
		require.Error(t, err)
	})

	t.Run("create resources with finalizers", func(t *testing.T) {
		mockValkeySvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		_, finalizers, err := flow.Run(ctx, databasev1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: defaultNamespace,
			},
		})
		require.NoError(t, err)
		require.Len(t, finalizers, 1)
		require.Equal(t, finalizers[0], valkey.Finalizer)
	})

	t.Run("invalid resource type", func(t *testing.T) {
		status, finalizers, err := flow.Run(ctx, &databasev1alpha1.ValkeyStatus{})
		require.Nil(t, status)
		require.Nil(t, finalizers)
		require.Error(t, err)
		require.Equal(t, flows.ErrInvalidInputType, err)
	})

	t.Run("update resources error", func(t *testing.T) {
		mockValkeySvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(mockErr)

		status, finalizers, err := flow.Run(ctx, databasev1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:       resourceName,
				Namespace:  defaultNamespace,
				Finalizers: []string{valkey.Finalizer},
			},
		})
		require.Nil(t, status)
		require.Nil(t, finalizers)
		require.Error(t, err)
	})

	t.Run("healthcheck error", func(t *testing.T) {
		mockValkeySvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockValkeySvc.EXPECT().IsReady(gomock.Any(), gomock.Any()).Return(false, int32(0), mockErr)

		status, finalizers, err := flow.Run(ctx, databasev1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:       resourceName,
				Namespace:  defaultNamespace,
				Finalizers: []string{valkey.Finalizer},
			},
		})
		require.Nil(t, status)
		require.Nil(t, finalizers)
		require.Error(t, err)
	})

	t.Run("not ready", func(t *testing.T) {
		mockValkeySvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockValkeySvc.EXPECT().IsReady(gomock.Any(), gomock.Any()).Return(false, int32(0), nil)

		status, finalizers, err := flow.Run(ctx, databasev1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:       resourceName,
				Namespace:  defaultNamespace,
				Finalizers: []string{valkey.Finalizer},
			},
		})
		require.NotNil(t, status)
		require.Len(t, finalizers, 0)
		require.NoError(t, err)
	})

	t.Run("success reconcile", func(t *testing.T) {
		mockValkeySvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockValkeySvc.EXPECT().IsReady(gomock.Any(), gomock.Any()).Return(true, int32(1), nil)

		status, finalizers, err := flow.Run(ctx, databasev1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:       resourceName,
				Namespace:  defaultNamespace,
				Finalizers: []string{valkey.Finalizer},
			},
		})
		require.Equal(t, &databasev1alpha1.ValkeyStatus{
			Status:        databasev1alpha1.TypeStatusHealthy,
			ReadyReplicas: 1,
		}, status)
		require.Len(t, finalizers, 1)
		require.Equal(t, finalizers[0], valkey.Finalizer)
		require.NoError(t, err)
	})

	t.Run("delete resource", func(t *testing.T) {
		mockValkeySvc.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)

		status, finalizers, err := flow.Run(ctx, databasev1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:              resourceName,
				Namespace:         defaultNamespace,
				Finalizers:        []string{valkey.Finalizer},
				DeletionTimestamp: &metav1.Time{Time: time.Now()},
			},
		})
		require.NotNil(t, status)
		require.Len(t, finalizers, 0)
		require.NoError(t, err)
	})

	t.Run("delete resource error", func(t *testing.T) {
		mockValkeySvc.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(mockErr)

		status, finalizers, err := flow.Run(ctx, databasev1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:              resourceName,
				Namespace:         defaultNamespace,
				Finalizers:        []string{valkey.Finalizer},
				DeletionTimestamp: &metav1.Time{Time: time.Now()},
			},
		})
		require.Nil(t, status)
		require.Nil(t, finalizers)
		require.Error(t, err)
	})
}
