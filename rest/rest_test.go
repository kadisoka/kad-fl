package rest

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, expected interface{}, actual interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\tExpected: %#v\n\tActual: %#v", expected, actual)
	}
}
