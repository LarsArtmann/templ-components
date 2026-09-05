package datastar

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/larsartmann/go-datastar/static"
)

// pinnedBundleSHA256 is the SHA-256 of the embedded runtime bundle
// (static.Bytes()) at go-datastar/static v0.5.0 — Datastar 1.0.3, 33538 bytes.
// Re-audited 2026-09-05 on the v0.4.0 → v0.5.0 bump: every pinned token
// survived; the bundle shrank (56330 → 33538 bytes, upstream minification
// refactor) and the retry machinery (retryInterval/retryMax/retryMaxCount +
// auto/error/never/always literals) is intact.
const pinnedBundleSHA256 = "5d6b7794a50a83d82da962aec5e382f5ae83ac7afbc751f903f7a9c6bd433c65"

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
		// Request-cancellation option consumed by
		// LiveRegionProps.Cancellation — under cleanup mode the runtime
		// aborts the in-flight stream when the element leaves the DOM.
		"requestCancellation",
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

	// Byte-integrity pin: the token assertions above tolerate arbitrary new
	// bytes as long as every pinned token survives. This hash closes that
	// gap — any byte-level bundle change fails here and forces a conscious
	// re-audit (docs/datastar-runtime-facts.md) before the pin is updated.
	sum := sha256.Sum256(static.Bytes())
	if got := hex.EncodeToString(sum[:]); got != pinnedBundleSHA256 {
		t.Errorf(
			"embedded bundle sha256 changed: got %s, want %s (static v%s) — re-audit per docs/datastar-runtime-facts.md",
			got,
			pinnedBundleSHA256,
			static.Version,
		)
	}
}

// TestDatastarVersionConstantNameMatchesValue guards the literal behind
// DatastarVersion1_0_3. Its value is derived from static.Version, so a
// go-datastar/static bump silently re-points the value while the NAME still
// claims 1.0.3 — a name that lies about its value is worse than none. On
// failure: re-audit the new runtime, rename the constant, update references.
func TestDatastarVersionConstantNameMatchesValue(t *testing.T) {
	t.Parallel()

	if got := string(DatastarVersion1_0_3); got != "1.0.3" {
		t.Errorf(
			"DatastarVersion1_0_3 = %q — static.Version moved and the constant name no longer tells the truth; re-audit the new bundle, rename the constant, and update references",
			got,
		)
	}
}
