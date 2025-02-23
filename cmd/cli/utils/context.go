package utils

import (
	"context"

	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var CtxClientKey = "k8s-client"

func CtxClient(ctx context.Context) runtimeclient.Client {
	return ctx.Value(CtxClientKey).(runtimeclient.Client)
}
