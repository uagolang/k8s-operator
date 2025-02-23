package valkey

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the cli/modules/resources command
var RootCmd = &cobra.Command{
	Use:   "valkey",
	Short: "Valkey K8s Resource Manager",
	Long:  "Manage Valkey inside k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}
