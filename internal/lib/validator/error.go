package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

// Error represents a validation error for responses.
type Error struct {
	Field   string `json:"field" example:"password"`
	Message string `json:"message" example:"Password is incorrect"`
}

func (e Error) Error() string {
	return e.Message
}

type Errors []Error

func (e Errors) Error() string {
	messages := lo.Map(e, func(err Error, _ int) string {
		return err.Error()
	})

	return strings.Join(messages, ", ")
}

func GetErrors(err error) Errors {
	errs := make(Errors, 0)

	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		for _, verr := range vErrs {
			errs = append(errs, newError(verr))
		}
	}

	var ferrs Errors
	if errors.As(err, &ferrs) {
		errs = append(errs, ferrs...)
	}

	var ferr Error
	if errors.As(err, &ferr) {
		errs = append(errs, ferr)
	}

	return errs
}

func newError(err validator.FieldError) Error {
	return Error{
		Field:   err.Field(),
		Message: err.Translate(trans),
	}
}
