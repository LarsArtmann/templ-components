package datastar

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/larsartmann/go-datastar/static"
)

// pinnedBundleSHA256 is the SHA-256 of the embedded runtime bundle
// (static.Bytes()) at go-datastar/static v0.5.0 → v0.6.1 — Datastar 1.0.3,
// 33538 bytes. Re-audited 2026-09-05 on the v0.4.0 → v0.5.0 bump (the bundle
// shrank, 56330 → 33538 bytes, upstream minification refactor; the retry
// machinery retryInterval/retryMax/retryMaxCount + auto/error/never/always
// literals survived) and re-verified 2026-10-02 on the v0.5.0 → v0.6.1 bump:
// BYTE-IDENTICAL (same sha256), so every pinned token carried over.
// Re-audited 2026-11-17 on the v0.6.1 → v0.7.0 bump (Datastar 1.0.3 → 1.0.4,
// 33553 bytes): minifier identifier renames throughout, the error formatter
// rewritten (multi-line `More info:`/`Context:` template), the CSP machinery
// now assembles the data-nonce attribute name at runtime (j("nonce") — the
// bare literal is no longer byte-present), the empty-nonce throw became the
// typed NonceRequired error, and expression compilation gained an explicit
// "Blocked by CSP." throw. Mechanism unchanged; every functional token
// re-verified present.
const pinnedBundleSHA256 = "727844adfc825ee651fb93c544a2a739986f9a21820a94524b35f0cac470cf91"

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
		// Form-serialization machinery consumed by wire.ContentTypeForm
		// (forms.Form Wires dual-transport submission): the fetch option
		// selecting form-field serialization over the signals-JSON default,
		// and the two failure paths (no enclosing form / unknown type).
		"contentType",
		"FetchFormNotFound",
		"FetchInvalidContentType",
		// Interval-polling plugin consumed by PolledRegion (ADR-0042): the
		// plugin name (attribute data-on-interval), the duration modifier,
		// and the leading flag (an element of the duration mod set, not a
		// separate group — PolledRegion deliberately never emits it).
		`name:"on-interval"`,
		`name:"on-intersect"`,
		"duration",
		"leading",
		// CSP-mode machinery (decoded 2026-11-07, see runtime-facts §CSP;
		// re-audited 2026-11-17 on the 1.0.4 bump): the runtime reads
		// data-nonce from documentElement ONCE at init — as of 1.0.4 the
		// attribute name is ASSEMBLED at runtime (j("nonce") → "data-nonce"),
		// so the byte-provable literals are the "nonce" argument, the typed
		// NonceRequired error (was a literal message), and the explicit
		// "Blocked by CSP." throw on expression-compile rejection. When
		// data-nonce is present and nonempty, expressions execute as nonce'd
		// inline scripts instead of Function() eval, so script-src 'nonce-…'
		// suffices without 'unsafe-eval'. trustedTypes is the self-registered
		// identity policy used under Trusted Types enforcement. A pin bump
		// that drops any of these fails here before a consumer CSP flips.
		`"nonce"`,
		"NonceRequired",
		"Blocked by CSP.",
		"trustedTypes",
		"createPolicy",
		// Verified present but UNMODELED surface (runtime-facts §"Verified
		// present, unmodeled"): the signal-patch listener plugin (attribute
		// data-on-signal-patch; plugin registrations are name:"…" literals,
		// the data- prefixes assemble at runtime, so only this literal is
		// byte-provable), the view-transition-selector dataline, and the
		// fetch options / datalines below.
		`name:"on-signal-patch"`,
		"viewTransitionSelector",
		"filterSignals",
		"openWhenHidden",
		"onlyIfMissing",
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
// DatastarVersion1_0_4. Its value is derived from static.Version, so a
// go-datastar/static bump silently re-points the value while the NAME still
// claims 1.0.4 — a name that lies about its value is worse than none. On
// failure: re-audit the new runtime, rename the constant, update references.
func TestDatastarVersionConstantNameMatchesValue(t *testing.T) {
	t.Parallel()

	if got := string(DatastarVersion1_0_4); got != "1.0.4" {
		t.Errorf(
			"DatastarVersion1_0_4 = %q — static.Version moved and the constant name no longer tells the truth; re-audit the new bundle, rename the constant, and update references",
			got,
		)
	}

	// The compat aliases are deliberate LITERALS (never derived), so their
	// names keep telling the truth even when the pin moves again.
	if got := string(DatastarVersion1_0_3); got != "1.0.3" {
		t.Errorf("DatastarVersion1_0_3 compat alias = %q, want literal %q",
			got, "1.0.3")
	}

	if got := string(DatastarVersion1_0_2); got != "1.0.2" {
		t.Errorf("DatastarVersion1_0_2 compat alias = %q, want literal %q",
			got, "1.0.2")
	}
}
