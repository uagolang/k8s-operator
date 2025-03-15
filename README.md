# Resource Operator for Kubernetes

Examples for the [articles cycle](https://t.me/uagolang/9) about managing k8s resources.
It should work on local environment with `minikube` as a local k8s cluster.

## Get started

### Dependencies

- Go 1.24
- Docker Engine
- Minikube
- kubectl

### Local environment

Firstly, you need to clone this repo:

```shell
# firstly
git clone https://github.com/uagolang/k8s-operator.git
# then
cd k8s-operator
go mod tidy
```

Then start minikube, check it status and use its context:

```shell
minikube start --nodes 2
minikube status
kubectl config use-context minikube
```

After that you need to execute `docker login` and then build Docker image 
and push it to the Docker Hub using command:

```shell
make docker-build docker-push
```

After that you need to generate all needed resources and deploy operator:

```shell
make generate
make manifests
make deploy

# use `make undeploy` to undeploy operator
# use `make uninstall` to uninstall CRDs
```

### Create new CRD using CLI tool

```shell
go run cmd/cli/main.go valkey create --name=test --namespace=test --image=valkey/valkey --user=root --pass=root --replicas=1 --volume_enabled=true --cpu=200m --memory=512Mi --storage=512Mi
```
