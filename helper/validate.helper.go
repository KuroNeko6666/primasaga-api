package helper

import "github.com/go-playground/validator/v10"

func Validate() *validator.Validate {
	return validator.New()
}
