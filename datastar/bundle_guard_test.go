package datastar

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
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
		// Reactive plugins consumed by SignalChip (the event, class, attr and
		// signals attributes). Their STRUCTURE is pinned by
		// TestPinnedRuntimeSignalChipContract below; these literals only prove
		// the plugins are still registered under these names.
		`name:"on"`,
		`name:"class"`,
		`name:"attr"`,
		`name:"signals"`,
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

// bundleWhitespaceRe collapses the bundle's line breaks and indentation so a
// structural pattern can match across statements.
//

var bundleWhitespaceRe = regexp.MustCompile(`\s+`)

// bundlePluginWindow returns the minified source of one plugin registration
// (from its name literal up to the next registration), whitespace-collapsed.
func bundlePluginWindow(t *testing.T, bundle, plugin string) string {
	t.Helper()

	marker := `name:"` + plugin + `"`

	start := strings.Index(bundle, marker)

	if start < 0 {
		t.Fatalf("pinned runtime bundle v%s no longer registers the %q plugin", static.Version, plugin)
	}

	const next = `m({name:`

	rel := strings.Index(bundle[start+len(marker):], next)

	if rel < 0 {
		t.Fatalf("no plugin registration follows %s in the bundle", marker)
	}

	return bundleWhitespaceRe.ReplaceAllString(bundle[start:start+len(marker)+rel], " ")
}

// TestPinnedRuntimeSignalChipContract pins, in the pinned bundle's own bytes,
// the four runtime facts SignalChip depends on. String assertions on the
// component's output cannot tell whether the runtime still parses that
// output; a pin bump that changes any structure below fails here first.
// Identifier names are minified and may change; every pattern matches them
// as placeholders and pins only the shape.
func TestPinnedRuntimeSignalChipContract(t *testing.T) {
	t.Parallel()

	bundle := string(static.Bytes())

	collapsed := bundleWhitespaceRe.ReplaceAllString(bundle, " ")

	tests := []struct {
		name    string
		subject string
		pattern string
		fact    string
	}{
		{
			name:    "attribute dispatcher",
			subject: collapsed,
			pattern: `pluginName:\w+,key:\w+,mods:\w+\}=\w+\(\w+\),\w+=\w+\.get\(\w+\);if\(\(!\w+\|\|\w+\.has\(\w+\)\)&&!!\w+\)\{`,
			fact:    "an attribute name is split into pluginName/key and the plugin is applied only when a registration is found, with no else branch: data-on-click (an unknown plugin name) is dropped silently, which is why the chip must spell data-on:click",
		},
		{
			name:    "class plugin",
			subject: bundlePluginWindow(t, bundle, "class"),
			pattern: `for\(let \w+ in \w+\)\{let \w+=\w+\.split\(/\\s\+/\)`,
			fact:    "the key-less form iterates the evaluated value with for..in and splits each key on whitespace: data-class takes an OBJECT of class names to booleans",
		},
		{
			name:    "class plugin truthiness",
			subject: bundlePluginWindow(t, bundle, "class"),
			pattern: `if\(\w+\[\w+\]\)for\(let \w+ of \w+\)\w+\.classList\.contains\(\w+\)\|\|\w+\.classList\.add\(\w+\);else for\(let \w+ of \w+\)\w+\.classList\.contains\(\w+\)&&\w+\.classList\.remove\(\w+\)`,
			fact:    "a truthy key adds its classes and a falsy one removes them, so both looks can travel in one data-class object",
		},
		{
			name:    "attr plugin true and false",
			subject: bundlePluginWindow(t, bundle, "attr"),
			pattern: `===""\|\|\w+===!0\?\w+\.setAttribute\(\w+,""\):\w+===!1\|\|\w+==null\?\w+\.removeAttribute\(\w+\):typeof \w+=="string"\?\w+\.setAttribute\(\w+,\w+\)`,
			fact:    "a boolean value sets an EMPTY attribute or REMOVES it, only a string sets a real value: aria-pressed must be driven by a string expression",
		},
		{
			name:    "signals plugin merges the object",
			subject: bundlePluginWindow(t, bundle, "signals"),
			pattern: `Object\.assign\(\{\},\w+\?\.\(\)\);\w+\(\w+,\{ifMissing:\w+\}\)`,
			fact:    "a key-less data-signals value is evaluated and merged into the signal store at init, which is how a selected chip declares the signal's initial value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if !regexp.MustCompile(tt.pattern).MatchString(tt.subject) {
				t.Errorf(
					"pinned runtime bundle v%s changed the %s structure this fact rests on: %s — re-audit SignalChip against the new bundle",
					static.Version,
					tt.name,
					tt.fact,
				)
			}
		})
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
