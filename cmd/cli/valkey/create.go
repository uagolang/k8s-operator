package valkey

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/uagolang/k8s-operator/api/v1alpha1"
	"github.com/uagolang/k8s-operator/cmd/cli/utils"
)

// CreateCmd represents the cli/modules/resources command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new ValKey using CRD",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		k8sClient := utils.CtxClient(cmd.Context())

		crdName := cmd.Flag("name").Value.String()
		ns := cmd.Flag("namespace").Value.String()

		image := cmd.Flag("image").Value.String()
		replicasStr := cmd.Flag("replicas").Value.String()
		replicas, err := strconv.Atoi(replicasStr)
		if err != nil {
			log.Fatal(err)
		}

		dbUser := cmd.Flag("user").Value.String()
		dbPass := cmd.Flag("pass").Value.String()

		cpu := cmd.Flag("cpu").Value.String()
		memory := cmd.Flag("memory").Value.String()

		volumeEnabled, _ := strconv.ParseBool(cmd.Flag("volume_enabled").Value.String())
		storage := cmd.Flag("storage").Value.String()

		err = k8sClient.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns,
			},
		})
		if err != nil && !errors.IsAlreadyExists(err) {
			log.Fatal(err)
		}

		// create CRD
		err = k8sClient.Create(ctx, &v1alpha1.Valkey{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crdName,
				Namespace: ns,
			},
			Spec: v1alpha1.ValkeySpec{
				Image:    image,
				Replicas: int32(replicas),
				User:     dbUser,
				Password: dbPass,
				Volume: v1alpha1.Volume{
					Enabled: volumeEnabled,
					Storage: storage,
				},
				Resource: v1alpha1.Resource{
					CPU:     cpu,
					Memory:  memory,
					Storage: storage,
				},
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("new Valkey resource created")
	},
}
