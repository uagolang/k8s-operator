package valkey

const (
	Namespace = "default"
)

func init() {
	CreateCmd.Flags().String("name", "app-db", "crd name")
	CreateCmd.Flags().String("namespace", Namespace, "namespace for resources")
	CreateCmd.Flags().String("image", "valkey/valkey", "service image")
	CreateCmd.Flags().String("user", "root", "admin user")
	CreateCmd.Flags().String("pass", "root", "admin password")
	CreateCmd.Flags().String("replicas", "1", "number of replicas")
	CreateCmd.Flags().String("volume_enabled", "true", "should have persistent volume")
	CreateCmd.Flags().String("cpu", "200m", "resource cpu")
	CreateCmd.Flags().String("memory", "512Mi", "resource memory")
	CreateCmd.Flags().String("storage", "512Mi", "resource storage")

	RootCmd.AddCommand(CreateCmd)
}
