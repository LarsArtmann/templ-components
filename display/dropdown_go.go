package display

import (
	"errors"
	"strconv"
)

func dropdownSafeID(id string) string {
	return strconv.Quote(id)
}

func validateDropdownID(id string) error {
	if id == "" {
		return errors.New(
			"Dropdown requires a non-empty ID for ARIA attributes and JavaScript functionality",
		)
	}
	return nil
}
