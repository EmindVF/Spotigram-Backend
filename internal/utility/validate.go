package utility

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validates a struct field by its field name.
// Returns false on failed validation tags checks or the lack of the field.
func IsValidStructField[T any](s T, fieldName string) bool {
	t := reflect.TypeOf(s)
	_, found := t.FieldByName(fieldName)
	if !found {
		return false
	}
	err := validate.StructPartial(s, fieldName)
	return err == nil
}
