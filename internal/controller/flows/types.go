package flows

import (
	"context"
	"errors"
)

var (
	ErrInvalidInputType  = errors.New("invalid input type")
	ErrInvalidOutputType = errors.New("invalid output type")
)

type Flow interface {
	Run(ctx context.Context, item any) (status any, finalizers []string, err error)
}
