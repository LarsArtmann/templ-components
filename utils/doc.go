// Package utils provides shared types, Tailwind class merging, and generic helpers
// used across all templ-components packages.
//
// BaseProps is embedded by every component props struct for consistent ID, Class,
// Attrs, AriaLabel, and Nonce propagation. Class() merges Tailwind classes with
// conflict resolution via tailwind-merge-go.
package utils
