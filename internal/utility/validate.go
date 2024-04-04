package utility

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func IsValidStructField[T any](s T, fieldName string) bool {
	t := reflect.TypeOf(s)
	_, found := t.FieldByName(fieldName)
	if !found {
		return false
	}
	err := validate.StructPartial(s, fieldName)
	return err == nil
}
