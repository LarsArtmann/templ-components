package datastar

// RetryMode controls automatic reconnection of Datastar backend actions
// when a stream ends or a request fails. Mirrors the retry argument the
// runtime accepts on @get/@post/... action expressions.
//
// Runtime behaviour matrix (behaviour verified against datastar.js v1.0.2;
// pinned tokens re-verified present in the v1.0.3 bundle on 2026-09-05):
//
//   - clean stream EOF (e.g. server restart): only RetryAlways reconnects.
//     The default auto disposes the connection, leaving the region stale.
//   - HTTP status >= 400: RetryAlways and RetryError reconnect.
//   - thrown network errors: every mode reconnects.
//   - the failure counter resets on every successful (200) connect, so
//     RetryAlways self-heals indefinitely across individual restarts.
type RetryMode string

const (
	// RetryAuto reconnects on network errors only (runtime default).
	RetryAuto RetryMode = "auto"
	// RetryAlways reconnects after clean stream EOF, HTTP errors, and
	// network errors. Use for long-lived streams that must self-heal
	// across server restarts.
	RetryAlways RetryMode = "always"
	// RetryError reconnects on HTTP (>=400) and network errors, but not
	// after a clean stream EOF.
	RetryError RetryMode = "error"
	// RetryNever never reconnects.
	RetryNever RetryMode = "never"
)

// validRetryModes are the allowed retry argument values.
//
//nolint:gochecknoglobals // Package-level validation set
var validRetryModes = map[RetryMode]bool{
	RetryAuto: true, RetryAlways: true, RetryError: true, RetryNever: true,
}

// retryModeValue returns v if valid, otherwise falls back to RetryAuto.
func retryModeValue(v RetryMode) RetryMode {
	if validRetryModes[v] {
		return v
	}

	return RetryAuto
}
