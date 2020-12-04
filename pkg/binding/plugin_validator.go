package binding

import (
	validator "github.com/go-playground/validator/v10"
	"github.com/system18188/jupiter-plugin/pkg/check"
	"regexp"
)

var (
	dateRegex = regexp.MustCompile(`^2\d{3}-\d{2}-\d{2}$`)
	timeRegex = regexp.MustCompile(`^2\d{3}-\d{2}-\d{2} \d{2}\:\d{2}\:\d{2}$`)
)

func IsNoHTML(fl validator.FieldLevel) bool {
	return !check.IsHTML(fl.Field().String())
}

func IsColumn(fl validator.FieldLevel) bool {
	return check.IsColumn(fl.Field().String())
}
