package visualtest

import (
	"testing"
)

// RetryOnce runs a flake-prone browser flow once and retries it exactly
// once on failure (docs/testing/flake-policy.md, M19/F088). Both attempts
// failing fails the test; the first failure is logged, never swallowed.
//
// The body reports failure by returning an error (not via t) so the retry
// starts from a clean slate without mutating outer test state.
//
// Use ONLY for browser-timing flakes whose root cause is documented in
// AGENTS.md. If a RetryOnce-wrapped test fails twice, that is class 3 in
// the policy (real bug): remove the wrapper and fix the cause.
func RetryOnce(t *testing.T, flow string, body func() error) {
	t.Helper()

	err := body()
	if err == nil {
		return
	}

	t.Logf("%s: first attempt failed (%v) — retrying once per flake policy", flow, err)

	if err := body(); err != nil {
		t.Errorf("%s: failed on retry too — not a flake (flake-policy class 3): %v", flow, err)
	}
}
