package datastar

// RequestCancellation controls what happens to an in-flight @get stream
// request when the element carrying data-init is removed from the DOM —
// typically an HTMX swap that replaces the region with a new instance
// pointing at a different URL.
type RequestCancellation string

const (
	// CancellationNone is the default: the option is omitted from the
	// action expression and the in-flight request survives element
	// removal. This is the runtime's own default behaviour.
	CancellationNone RequestCancellation = ""
	// CancellationCleanup aborts the request when its initiating element
	// is removed from the DOM. Without it, a swapped-out region's old
	// stream keeps running and re-patches the region with stale content
	// that races the replacement's stream.
	CancellationCleanup RequestCancellation = "cleanup"
)

// validRequestCancellation are the allowed requestCancellation option values.
//
//nolint:gochecknoglobals // Package-level validation set
var validRequestCancellation = map[RequestCancellation]bool{
	CancellationNone: true, CancellationCleanup: true,
}

// requestCancellationValue returns c if valid, otherwise falls back to none.
func requestCancellationValue(c RequestCancellation) RequestCancellation {
	if validRequestCancellation[c] {
		return c
	}

	return CancellationNone
}
