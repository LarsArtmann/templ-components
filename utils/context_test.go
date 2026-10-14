package utils

import (
	"context"
	"testing"
)

// TestNoncePrecedence pins the ADR-0047 contract: explicit props ALWAYS win,
// context only fills empty, and both-empty yields the zero value. A violation
// here would make struct literals render differently depending on invisible
// request state.
func TestNoncePrecedence(t *testing.T) {
	t.Parallel()

	ctx := ContextWithNonce(context.Background(), "ctx-nonce")

	tests := []struct {
		name     string
		withCtx  bool
		explicit string
		want     string
	}{
		{name: "explicit wins over context", withCtx: true, explicit: "prop-nonce", want: "prop-nonce"},
		{name: "context fills empty", withCtx: true, explicit: "", want: "ctx-nonce"},
		{name: "both empty yields zero value", withCtx: false, explicit: "", want: ""},
		{name: "empty context value treated as absent", withCtx: false, explicit: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testCtx := context.Background()
			if tt.withCtx {
				testCtx = ctx
			}

			if got := ResolveNonce(testCtx, tt.explicit); got != tt.want {
				t.Errorf("ResolveNonce() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestNonceFromContextZeroFallback proves the getter never panics and returns
// the zero value for foreign contexts and wrong value types.
func TestNonceFromContextZeroFallback(t *testing.T) {
	t.Parallel()

	if got := NonceFromContext(context.Background()); got != "" {
		t.Errorf("bare context: got %q, want empty", got)
	}

	if got := NonceFromContext(context.WithValue(context.Background(), nonceKey{}, 42)); got != "" {
		t.Errorf("wrong value type: got %q, want empty (type assertion fails safe)", got)
	}

	ctx := ContextWithNonce(context.Background(), "abc")

	if got := NonceFromContext(ctx); got != "abc" {
		t.Errorf("round trip: got %q, want abc", got)
	}
}
