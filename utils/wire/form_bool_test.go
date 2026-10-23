package wire

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

// TestSetBoolFormFieldRejectsUnparseableInput proves malformed checkbox values
// surface as an error (naming the field) instead of silently binding false.
func TestSetBoolFormFieldRejectsUnparseableInput(t *testing.T) {
	t.Parallel()

	field := reflect.New(reflect.TypeOf(false)).Elem()

	err := setFormValue(field, "subscribe", "definitely-not-a-bool")
	if err == nil {
		t.Fatal("setFormValue accepted an unparseable bool")
	}

	if !errors.Is(err, strconv.ErrSyntax) {
		t.Errorf("error %v does not wrap strconv.ErrSyntax", err)
	}

	if field.Bool() {
		t.Error("failed parse mutated the field to true")
	}
}
