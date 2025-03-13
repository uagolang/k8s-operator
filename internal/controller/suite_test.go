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
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	databasev1alpha1 "github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/internal/controller/flows"
	"github.com/uagolang/k8s-operator/mocks"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg              *rest.Config
	k8sClient        client.Client
	testEnv          *envtest.Environment
	controllerValkey *ValkeyReconciler
	mockErr          = errors.New("mock error")
)

var (
	mockK8sClient       *mocks.MockK8sClient
	mockK8sStatusClient *mocks.MockK8sStatusClient
	mockFlow            *mocks.MockFlow
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}
}

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
			fmt.Sprintf("1.29.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = databasev1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	// use fake k8s client (in-memory)
	// instead of real cluster connection
	k8sClient = fake.NewClientBuilder().
		WithScheme(scheme.Scheme).
		// you can add some object for seed cluster
		// with needed resource objects
		WithRuntimeObjects(flows.FakeComponents...).
		// register Status sub-resources to have an ability
		// to Update resource Status object
		WithStatusSubresource(
			&databasev1alpha1.Valkey{},
		).
		Build()
	// just test that fake client was initialized
	Expect(k8sClient).NotTo(BeNil())

	mockCtrl := gomock.NewController(GinkgoT())

	// init mocks
	mockK8sClient = mocks.NewMockK8sClient(mockCtrl)
	mockFlow = mocks.NewMockFlow(mockCtrl)
	mockK8sStatusClient = mocks.NewMockK8sStatusClient(mockCtrl)

	// init crd controllers
	controllerValkey = &ValkeyReconciler{
		Client: k8sClient,
		Scheme: k8sClient.Scheme(),
		Flow:   mockFlow,
	}
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
