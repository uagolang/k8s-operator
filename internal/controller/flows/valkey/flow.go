package valkey

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/internal/controller/flows"
	valkeysvc "github.com/uagolang/k8s-operator/internal/services/valkey"
)

type FlowImpl struct {
	k8sClient client.Client
	valkeySvc valkeysvc.Service
}

type ImplOption func(r *FlowImpl)

func NewFlow(opts ...ImplOption) flows.Flow {
	res := new(FlowImpl)
	for _, opt := range opts {
		opt(res)
	}

	return res
}

func WithK8sClient(v client.Client) ImplOption {
	return func(r *FlowImpl) {
		r.k8sClient = v
	}
}

func WithValkeySvc(v valkeysvc.Service) ImplOption {
	return func(r *FlowImpl) {
		r.valkeySvc = v
	}
}

func (r *FlowImpl) Run(ctx context.Context, input any) (any, []string, error) {
	item, ok := input.(v1alpha1.Valkey)
	if !ok {
		return nil, nil, flows.ErrInvalidInputType
	}

	logger := log.FromContext(ctx).WithValues("flow", "valkey", "crd_name", item.Name)
	log.IntoContext(ctx, logger)

	res := new(v1alpha1.ValkeyStatus)

	if !item.DeletionTimestamp.IsZero() { // should be deleted
		return r.delete(ctx, &item)
	}

	if len(item.Finalizers) == 0 { // save finalizers
		var err error
		logger.Info("finalizer was added to valkey")

		_, err = r.valkeySvc.Create(ctx, &valkeysvc.CreateRequest{
			CrdName:   item.Name,
			Namespace: item.Namespace,
			Image:     item.Spec.Image,
			User:      item.Spec.User,
			Password:  item.Spec.Password,
			Replicas:  item.Spec.Replicas,
			Volume:    item.Spec.Volume,
			Resource:  item.Spec.Resource,
		})
		if err != nil {
			return nil, nil, err
		}

		res.Status = StatusInProgress

		return res, []string{Finalizer}, nil
	}

	err := r.valkeySvc.Update(ctx, &valkeysvc.UpdateRequest{
		CrdName:   item.Name,
		Namespace: item.Namespace,
		Image:     &item.Spec.Image,
		User:      &item.Spec.User,
		Password:  &item.Spec.Password,
		Replicas:  &item.Spec.Replicas,
		Volume:    &item.Spec.Volume,
		Resource:  &item.Spec.Resource,
	})
	if err != nil {
		return nil, nil, err
	}

	ready, readyReplicas, err := r.valkeySvc.IsReady(ctx, &item)
	if err != nil {
		return nil, nil, err
	}
	if !ready || readyReplicas == 0 {
		res.Status = StatusStopped
		return res, []string{}, nil
	}

	res.ReadyReplicas = readyReplicas
	res.Status = StatusRunning

	return res, item.Finalizers, nil
}

func (r *FlowImpl) delete(ctx context.Context, item *v1alpha1.Valkey) (status *v1alpha1.ValkeyStatus, finalizers []string, err error) {
	logger := log.FromContext(ctx)

	if len(item.Finalizers) > 0 {
		err = r.valkeySvc.Delete(ctx, types.NamespacedName{
			Name:      item.Name,
			Namespace: item.Namespace,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	logger.Info("valkey resources were successfully deleted")

	return status, []string{}, nil
}
