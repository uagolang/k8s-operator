package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/cmd/cli/utils"
	"github.com/uagolang/k8s-operator/cmd/cli/valkey"
)

var scheme = runtime.NewScheme()

var rootCmd = &cobra.Command{
	Use:   "uagolang",
	Short: "UA Golang CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	rootCmd.AddCommand(valkey.RootCmd)
}

func main() {
	_ = godotenv.Load("./cmd/cli/.env")
	var err error
	ctx := context.Background()

	k8sClient, err := createK8sClient(ctx, os.Getenv("K8S_KUBECONFIG_PATH"))
	if err != nil {
		log.Fatal("failed to create k8s client: ", err)
	}

	rootCmd.SetContext(context.WithValue(ctx, utils.CtxClientKey, k8sClient))

	if err = rootCmd.Execute(); err != nil {
		log.Fatal("failed command execution: ", err)
	}
}

func createK8sClient(ctx context.Context, path string) (runtimeclient.Client, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}

	client, err := runtimeclient.New(cfg, runtimeclient.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}

	err = client.List(ctx, &corev1.NamespaceList{}, &runtimeclient.ListOptions{})
	if err != nil {
		return nil, err
	}

	return client, nil
}
