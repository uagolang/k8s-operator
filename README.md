# Resource Operator for Kubernetes

Examples for the [articles cycle](https://t.me/uagolang/9) about managing k8s resources.
It should work on local environment with `minikube` as a local k8s cluster.

## Get started

### Dependencies

- Go 1.24
- Docker Engine
- Minikube
- kubectl

### Clone repo

Firstly, you need to clone this repo:

```shell
# firstly
git clone https://github.com/uagolang/k8s-operator.git
# then
go mod tidy
```

### Minikube

You should be sure that your kubectl context is correct:

```shell
minikube start
minikube status
kubectl config use-context minikube
```

### Add CRD to k8s cluster (minikube)

```shell
make generate
make manifests
make install
```

```shell
# result of make install
~/apps/uagolang/k8s-operator git:[main]
make install
/Users/nesymno/apps/uagolang/k8s-operator/bin/controller-gen-v0.14.0 rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/nesymno/apps/uagolang/k8s-operator/bin/kustomize-v5.3.0 build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/valkeys.database.kuberly.io created
```

### Create new CRD using CLI tool

```shell
go run cmd/cli/main.go valkey create --name=test --namespace=test --image=valkey/valkey --user=root --pass=root --replicas=1 --volume_enabled=true --cpu=200m --memory=512Mi --storage=512Mi
```

### Deploy operator to cluster

```shell
make deploy
```
