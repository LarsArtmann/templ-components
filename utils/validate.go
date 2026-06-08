package utils

import "fmt"

// ValidateID returns an error if id is empty, using componentName in the message.
func ValidateID(componentName, id string) error {
	if id == "" {
		return fmt.Errorf("%s: id=%q cannot be empty", componentName, id)
	}
	return nil
}
