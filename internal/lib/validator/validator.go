package validator

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	v     = validator.New(validator.WithRequiredStructEnabled())
	trans ut.Translator
)

func init() {
	v = validator.New(validator.WithRequiredStructEnabled())
	trans, _ = ut.New(en.New()).GetTranslator("en")

	_ = entranslations.RegisterDefaultTranslations(v, trans)

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func Validate(ctx context.Context, i interface{}) error {
	return v.StructCtx(ctx, i)
}
