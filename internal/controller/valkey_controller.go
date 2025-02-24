/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	databasev1alpha1 "github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/internal/controller/flows"
	"github.com/uagolang/k8s-operator/internal/controller/flows/valkey"
	"github.com/uagolang/k8s-operator/internal/utils"
)

// ValkeyReconciler reconciles a Valkey object
type ValkeyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Flow   flows.Flow
}

//+kubebuilder:rbac:groups=database.kuberly.io,resources=valkeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.kuberly.io,resources=valkeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=database.kuberly.io,resources=valkeys/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ValkeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var emptyResp ctrl.Result
	requeueRes := ctrl.Result{RequeueAfter: 10 * time.Second}

	item := new(databasev1alpha1.Valkey)
	if err := r.Get(ctx, req.NamespacedName, item); err != nil {
		if k8serrors.IsNotFound(err) {
			return emptyResp, reconcile.TerminalError(err)
		} else {
			return emptyResp, err
		}
	}

	status := new(databasev1alpha1.ValkeyStatus)
	statusItem, finalizers, err := r.Flow.Run(ctx, *item)
	if err == nil {
		var ok bool
		status, ok = statusItem.(*databasev1alpha1.ValkeyStatus)
		if !ok {
			return emptyResp, flows.ErrInvalidOutputType
		}
	} else {
		status = &databasev1alpha1.ValkeyStatus{
			Status:          valkey.StatusFailed,
			LastReconcileAt: utils.Pointer(metav1.Now()),
			Error:           err.Error(),
		}
	}

	shouldUpdateFinalizers := !utils.SlicesEqualSorted(item.Finalizers, finalizers)
	if err == nil && shouldUpdateFinalizers {
		item.Finalizers = finalizers
		if err = r.Update(ctx, item); err != nil {
			if k8serrors.IsNotFound(err) {
				return emptyResp, reconcile.TerminalError(err)
			}

			return emptyResp, err
		}

		return emptyResp, nil
	}

	changed := item.Status.IsChanged(status)
	if !changed {
		return requeueRes, nil
	}

	item.Status = *status
	item.Status.LastReconcileAt = utils.Pointer(metav1.Now())
	err = r.Status().Update(ctx, item)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return emptyResp, reconcile.TerminalError(err)
		}

		return emptyResp, err
	}

	return emptyResp, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ValkeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.Valkey{}).
		Complete(r)
}
