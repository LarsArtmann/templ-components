package wire

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// ErrUnsupportedFormField reports a DecodeForm target field whose kind the
// decoder cannot fill (only string, int, int64, and bool are supported).
var ErrUnsupportedFormField = errors.New("unsupported form field kind")

// DecodeForm decodes a submitted HTML form into T. It is the server-side
// counterpart of every wire.Action that submits form-encoded data: fields opt
// in via `form:"name"` struct tags (the same names the wired component sends),
// untagged fields keep their zero values, and absent or empty fields also keep
// their zero values. Both POST bodies and GET query parameters are read
// (http.Request.ParseForm semantics: body values take precedence).
//
// A field that is present but malformed (e.g. a non-numeric value for an int
// field) fails with an error naming the field; a field of an unsupported kind
// fails with ErrUnsupportedFormField. Domain validation (required fields,
// ranges, negative indices) stays with the caller — DecodeForm only decodes.
//
//	display.ParseKanbanMove is the in-repo example: it decodes via
//	DecodeForm, then applies the board's own rules.
func DecodeForm[T any](r *http.Request) (T, error) {
	var out T

	if err := r.ParseForm(); err != nil {
		return out, fmt.Errorf("parse form: %w", err)
	}

	if err := decodeFormInto(r, reflect.ValueOf(&out).Elem()); err != nil {
		return out, fmt.Errorf("decode form: %w", err)
	}

	return out, nil
}

func decodeFormInto(r *http.Request, target reflect.Value) error {
	structType := target.Type()

	for fieldIndex := range structType.NumField() {
		tag, tagged := structType.Field(fieldIndex).Tag.Lookup("form")
		if !tagged {
			continue
		}

		name, _, _ := strings.Cut(tag, ",")
		if name == "" || name == "-" {
			continue
		}

		raw := r.Form.Get(name)
		if raw == "" {
			continue
		}

		if err := setFormValue(target.Field(fieldIndex), name, raw); err != nil {
			return err
		}
	}

	return nil
}

func setFormValue(field reflect.Value, name, raw string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)

	case reflect.Int, reflect.Int64:
		number, err := strconv.ParseInt(raw, 10, field.Type().Bits())
		if err != nil {
			return fmt.Errorf("field %q: parse int: %w", name, err)
		}

		field.SetInt(number)

	case reflect.Bool:
		// HTML checkboxes submit "on" when checked; everything else follows
		// strconv.ParseBool.
		parsed := raw == "on"

		if !parsed {
			var err error

			if parsed, err = strconv.ParseBool(raw); err != nil {
				return fmt.Errorf("field %q: parse bool: %w", name, err)
			}
		}

		field.SetBool(parsed)

	default:
		return fmt.Errorf("%w: %s is a %s", ErrUnsupportedFormField, name, field.Kind())
	}

	return nil
}
