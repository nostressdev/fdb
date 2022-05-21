package tests

import (
	"fmt"
	"reflect"
	"testing"
)

func AssertError(t *testing.T, err error) {
	if err != nil {
		//debug.PrintStack()
		t.Fatalf("unexpected error %v", err)
	}
}

func AssertEqual(t testing.TB, expected interface{}, actual interface{}, description ...string) {
	if reflect.DeepEqual(expected, actual) {
		return
	}
	errText := ""
	if len(description) == 0 {
		errText = fmt.Sprintf("unexpected value %v, expected %v", actual, expected)
	} else {
		errText = fmt.Sprintf("%v. unexpected value %v, expected %v", description[0], actual, expected)
	}
	t.Fatalf(errText)
}

func AssertNotEqual(t testing.TB, expected interface{}, actual interface{}, description ...string) {
	if !reflect.DeepEqual(expected, actual) {
		return
	}
	errText := ""
	if len(description) == 0 {
		errText = fmt.Sprintf("unexpected value %v, expected %v", actual, expected)
	} else {
		errText = fmt.Sprintf("%v. unexpected value %v, expected %v", description[0], actual, expected)
	}
	t.Fatalf(errText)
}
