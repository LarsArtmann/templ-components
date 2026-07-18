package utils

import (
	"crypto/rand"
	"fmt"
	"sync/atomic"
	"time"
)

// ValidateID returns an error if id is empty, using componentName in the message.
func ValidateID(componentName, id string) error {
	if id == "" {
		//nolint:err113 // dynamic: message embeds caller componentName + id
		return fmt.Errorf("%s: id=%q cannot be empty", componentName, id)
	}

	return nil
}

// ensureIDCounter is a process-wide atomic counter used as a fallback when
// crypto/rand fails (which should never happen on a healthy system). Combined
// with time.Now().UnixNano() it provides uniqueness without predictability.
//
//nolint:gochecknoglobals // Atomic counter for EnsureID fallback (process-wide)
var ensureIDCounter atomic.Uint64

// EnsureID returns id if non-empty, otherwise generates a unique ID with the
// given prefix using crypto/rand for collision safety across HTMX page loads.
// Format: tc-<prefix>-<16 hex chars> (e.g. tc-modal-a1b2c3d4e5f6a7b8).
// If crypto/rand fails (extremely unlikely), falls back to an atomic counter
// + nanosecond timestamp — never raw nanoseconds (predictable under concurrency).
func EnsureID(prefix, id string) string {
	if id != "" {
		return id
	}

	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		c := ensureIDCounter.Add(1)

		return fmt.Sprintf("tc-%s-%d-%d", prefix, time.Now().UnixNano(), c)
	}

	return fmt.Sprintf("tc-%s-%x", prefix, b[:])
}
