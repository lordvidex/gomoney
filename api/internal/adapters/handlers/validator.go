package handlers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/lordvidex/gomoney/api/internal/adapters/handlers/response"
	g "github.com/lordvidex/gomoney/pkg/gomoney"
	"reflect"
	"strings"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
}

func validateStruct(obj interface{}) []response.Error {
	var errors []response.Error
	err := validate.Struct(obj)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := response.Error{
				Error:   g.ErrInvalidInput.String(),
				Code:    response.C(g.ErrInvalidInput),
				Message: fmt.Sprintf("validation error on field %s, condition: %s", err.Field(), err.Tag())}
			errors = append(errors, element)
		}
		return errors
	}
	return nil
}
