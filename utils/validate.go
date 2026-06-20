package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// ValidateID returns an error if id is empty, using componentName in the message.
func ValidateID(componentName, id string) error {
	if id == "" {
		return fmt.Errorf("%s: id=%q cannot be empty", componentName, id)
	}
	return nil
}

// EnsureID returns id if non-empty, otherwise generates a unique ID with the
// given prefix using crypto/rand for collision safety across HTMX page loads.
// Format: tc-<prefix>-<16 hex chars> (e.g. tc-modal-a1b2c3d4e5f6a7b8).
func EnsureID(prefix, id string) string {
	if id != "" {
		return id
	}
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("tc-%s-%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("tc-%s-%x", prefix, b[:])
}
