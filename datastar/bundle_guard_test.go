package datastar

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-datastar/static"
)

// TestPinnedRuntimeBundleContract pins the runtime surface this library
// integrates against, as verified byte-for-byte in the embedded bundle.
// SSEErrorHandling was inert for its entire life because it listened for a
// "datastar-sse-error" event that the runtime never dispatched — string
// assertions against our own output cannot catch that class of bug.
//
// A go-datastar/static version bump that renames or removes any pinned
// token fails here, forcing a re-audit of the integration instead of a
// silent breakage in LiveRegion, SSEErrorHandling, or the demo wire format.
func TestPinnedRuntimeBundleContract(t *testing.T) {
	t.Parallel()

	bundle := string(static.Bytes())

	for _, token := range []string{
		// The only lifecycle observability surface: SSEErrorHandling
		// listens to this document CustomEvent (detail.type is
		// started/finished/error/retrying/retries-failed).
		"datastar-fetch",
		// The only SSE event names with registered handlers in v1.x —
		// the demo stream and go-datastar's PatchElements emit these.
		"datastar-patch-elements",
		"datastar-patch-signals",
		// Retry-mode literals consumed by LiveRegionProps.Retry
		// ("always" also appears in the retry comparison chain).
		`"always"`,
		`"error"`,
		`"never"`,
	} {
		if !strings.Contains(bundle, token) {
			t.Errorf(
				"pinned runtime bundle v%s no longer contains %q — the Datastar integration is broken; re-audit against the new bundle",
				static.Version,
				token,
			)
		}
	}

	// The event SSEErrorHandling originally listened for never existed in
	// v1.x. If a future bundle dispatches it, prefer listening to it again
	// over the datastar-fetch lifecycle approximation.
	if strings.Contains(bundle, "datastar-sse-error") {
		t.Log("bundle now dispatches datastar-sse-error — consider migrating SSEErrorHandling back to it")
	}
}
