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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	databasev1alpha1 "github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/internal/controller/flows/valkey"
)

const defaultNamespace = "default"

var _ = Describe("Valkey Controller", func() {
	Context("Resource reconcile process", func() {
		const resourceName = "test-resource"
		const valkeyImage = "valkey/valkey:latest"
		const storage = "1Gi"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: defaultNamespace,
		}
		resourceObjectMeta := metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: defaultNamespace,
		}

		resource := &databasev1alpha1.Valkey{
			ObjectMeta: resourceObjectMeta,
			Spec: databasev1alpha1.ValkeySpec{
				Image:    valkeyImage,
				Replicas: 1,
				User:     "user",
				Password: "password",
				Volume: databasev1alpha1.Volume{
					Enabled: true,
					Storage: storage,
				},
				Resource: databasev1alpha1.Resource{
					CPU:     "200m", // 0.2 CPU
					Memory:  "256Mi",
					Storage: storage,
				},
			},
		}

		BeforeEach(func() {
			By("beforeEach: create Valkey")
			resource.ResourceVersion = ""
			Expect(controllerValkey.Client.Create(ctx, resource)).To(Succeed())
		})

		AfterEach(func() {
			resource := &databasev1alpha1.Valkey{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("afterEach: cleanup Valkey")
			if len(resource.GetFinalizers()) > 0 {
				resource.Finalizers = []string{}
				Expect(controllerValkey.Client.Update(ctx, resource)).To(Succeed())
			}
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully reconcile the resource", func() {
			mockFlow.EXPECT().Run(gomock.Any(), gomock.Any()).Return(&databasev1alpha1.ValkeyStatus{
				Status: "healthy",
			}, []string{valkey.Finalizer}, nil)

			_, err := controllerValkey.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
