package datastar

import (
	"strings"
	"testing"
)

// FuzzGetActionExpr verifies getActionExpr never panics and always emits a
// structurally valid @get expression for arbitrary URL + retry + cancellation
// combinations. Mirrors FuzzWriteDatastarPatch (examples/demo/sse_test.go):
// these strings land in data-init/data-on attribute values, so a malformed
// expression would silently deaden the region at runtime.
func FuzzGetActionExpr(f *testing.F) {
	f.Add("/api/stream", "always", "cleanup")
	f.Add("", "", "")
	f.Add("/x?it's=1\\", "never", "none")
	f.Add("javascript:'}", "bogus", "bogus")

	f.Fuzz(func(t *testing.T, url, retry, cancellation string) {
		out := getActionExpr(url, RetryMode(retry), RequestCancellation(cancellation))

		if !strings.HasPrefix(out, "@get('") {
			t.Fatalf("missing @get(' prefix: %q", out)
		}

		// Suffix depends on option presence (deterministic from the validated
		// enums): ') with no options, }) with an options object.
		hasOpts := retryModeValue(RetryMode(retry)) != RetryAuto ||
			requestCancellationValue(RequestCancellation(cancellation)) == CancellationCleanup
		wantSuffix := "})"

		if !hasOpts {
			wantSuffix = "')"
		}

		if !strings.HasSuffix(out, wantSuffix) {
			t.Fatalf("missing %s suffix: %q", wantSuffix, out)
		}

		// The URL literal is always the escaped URL, immediately after the
		// opening quote (single quotes become backslash-escaped).
		inner := strings.TrimPrefix(out, "@get('")

		if escaped := strings.ReplaceAll(url, "'", "\\'"); !strings.HasPrefix(inner, escaped) {
			t.Fatalf("URL literal not the escaped URL:\n got prefix: %q\nwant: %q", inner, escaped)
		}

		// Quote accounting: 2 delimiters + one per URL quote (escaped ones
		// keep their quote) + 2 per emitted option value. Anything else means
		// an option or the URL broke the expression shape.
		wantQuotes := 2 + strings.Count(url, "'")

		if retryModeValue(RetryMode(retry)) != RetryAuto {
			wantQuotes += 2
		}

		if requestCancellationValue(RequestCancellation(cancellation)) == CancellationCleanup {
			wantQuotes += 2
		}

		if got := strings.Count(out, "'"); got != wantQuotes {
			t.Fatalf("quote count = %d, want %d (expression shape broken): %q", got, wantQuotes, out)
		}
	})
}

// FuzzActionExpr verifies actionExpr never panics and always emits
// @<method>('<escaped-url>') for arbitrary method + URL input. The method is
// always an internal constant (get/post/...) at every call site; fuzzing it
// anyway pins that nothing in the builder depends on method sanity.
func FuzzActionExpr(f *testing.F) {
	f.Add("get", "/api/x")
	f.Add("post", "/save?it's=1")
	f.Add("", "")
	f.Add("PUT", "\\ weird ' url")

	f.Fuzz(func(t *testing.T, method, url string) {
		out := actionExpr(method, url)

		if !strings.HasPrefix(out, "@") {
			t.Fatalf("missing @ prefix: %q", out)
		}

		if !strings.HasSuffix(out, "')") {
			t.Fatalf("missing ') suffix: %q", out)
		}

		// Quote accounting: exactly the two delimiters plus one quote per URL
		// quote and one per METHOD quote (escaped URL quotes keep their quote
		// character; the method is interpolated verbatim). With arbitrary
		// strings the URL-literal boundaries are ambiguous, so the count is
		// the unambiguous shape check.
		wantQuotes := 2 + strings.Count(url, "'") + strings.Count(method, "'")
		if got := strings.Count(out, "'"); got != wantQuotes {
			t.Fatalf("quote count = %d, want %d (expression shape broken): %q", got, wantQuotes, out)
		}
	})
}
