package display

import "errors"

func validateDropdownID(id string) error {
	if id == "" {
		return errors.New(
			"Dropdown requires a non-empty ID for ARIA attributes and JavaScript functionality",
		)
	}
	return nil
}
