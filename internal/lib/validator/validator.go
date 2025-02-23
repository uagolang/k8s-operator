package validator

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var v = validator.New(validator.WithRequiredStructEnabled())

func Validate(ctx context.Context, i interface{}) error {
	return v.StructCtx(ctx, i)
}
