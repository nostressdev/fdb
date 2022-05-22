package utils

import "github.com/nostressdev/fdb/errors"

func Validate(expression bool, text string) {
	if !expression {
		panic(errors.ValidationError.New(text))
	}
}

func Validatef(expression bool, format string, args ...interface{}) {
	if expression {
		panic(errors.ValidationError.Newf(format, args...))
	}
}