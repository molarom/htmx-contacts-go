package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni      ut.Translator
	validate *validator.Validate
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	uni, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	_ = en_translations.RegisterDefaultTranslations(validate, uni)
}

type Errs map[string]string

func (errs Errs) Error() string {
	sl := make([]string, 0, len(errs))
	for f, e := range errs {
		sl = append(sl, fmt.Sprintf("%s: %s", f, e))
	}

	return strings.Join(sl, ", ")
}

func Verify(val any) error {
	if err := validate.Struct(val); err != nil {

		var ve validator.ValidationErrors
		if ok := errors.As(err, &ve); !ok {
			return err
		}

		errs := make(Errs, len(ve))
		for _, e := range ve {
			errs[e.Field()] = e.Translate(uni)
		}

		return errs
	}

	return nil
}
