package apputils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Helper to convert validator.ValidationErrors to a readable map using json tag
func ValidationErrorsToMap(err error, obj any) map[string]string {
	errs := map[string]string{}
	if ve, ok := err.(validator.ValidationErrors); ok {
		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		for _, fe := range ve {
			field := fe.Field()
			jsonTag := ""
			if structField, found := t.FieldByName(fe.StructField()); found {
				jsonTag = structField.Tag.Get("json")
			}
			if jsonTag != "" && jsonTag != "-" {
				field = strings.Split(jsonTag, ",")[0]
			} else {
				field = strings.ToLower(field[:1]) + field[1:]
			}
			tag := fe.Tag()
			var msg string
			switch tag {
			case "required":
				msg = fmt.Sprintf("The %s field is required", field)
			case "uuid":
				msg = "Must be a valid UUID"
			case "min":
				msg = fmt.Sprintf("Minimum length is %s", fe.Param())
			case "eqfield":
				msg = "Must match " + fe.Param()
			default:
				msg = "Invalid value"
			}
			errs[field] = msg
		}
	} else {
		errs["error"] = err.Error()
	}
	return errs
}
