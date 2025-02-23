# k8s-operator

Examples for the articles cycle about managing k8s resources

## CLI Commands

```shell
go run cmd/cli/main.go valkey create --name=test --namespace=test --image=valkey/valkey --user=root --pass=root --replicas=1 --volume_enabled=true --cpu=200m --memory=512Mi --storage=512Mi
```
