package display

import "fmt"

func validateDropdownID(id string) error {
	if id == "" {
		return fmt.Errorf("dropdown: id=%q cannot be empty", id)
	}
	return nil
}
