package utils

// Motion class constants for consistent animation timing across all packages.
// Durations follow motion design research guidelines:
//   - 150ms: micro-interactions (hover, toggle, tooltip)
//   - 200ms: overlays (modal, drawer, dropdown)
//   - 300ms: panel transitions, large element movement
//
// Every transition includes motion-reduce fallbacks for vestibular safety.
// Use these instead of inline strings so timing is consistent and auditable.
const (
	TransitionFast      = "transition-all duration-150 ease-out motion-reduce:transition-none motion-reduce:duration-0"
	TransitionNormal    = "transition-all duration-200 ease-out motion-reduce:transition-none motion-reduce:duration-0"
	TransitionColors    = "transition-colors motion-reduce:transition-none motion-reduce:duration-0"
	TransitionTransform = "transition-transform motion-reduce:transition-none"
)
