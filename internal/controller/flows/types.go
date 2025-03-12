package flows

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/runtime"
)

var (
	ErrInvalidInputType  = errors.New("invalid input type")
	ErrInvalidOutputType = errors.New("invalid output type")
)

type Flow interface {
	Run(ctx context.Context, item any) (status any, finalizers []string, err error)
}

var FakeComponents = []runtime.Object{}
