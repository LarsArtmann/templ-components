package utils

import "context"

// nonceKey is the unexported context key for the request-scoped CSP nonce.
// Zero-size struct: collisions impossible across packages.
type nonceKey struct{}

// ContextWithNonce returns a child context carrying the request-scoped CSP
// nonce. Call it once in middleware; every context-aware component resolves
// the nonce from it (see ResolveNonce).
func ContextWithNonce(ctx context.Context, nonce string) context.Context {
	return context.WithValue(ctx, nonceKey{}, nonce)
}

// NonceFromContext returns the request-scoped nonce, or "" when absent —
// the zero value, never a panic and never a sentinel.
func NonceFromContext(ctx context.Context) string {
	nonce, _ := ctx.Value(nonceKey{}).(string)

	return nonce
}

// ResolveNonce applies the library's context precedence (ADR-0047):
// explicit props > context > zero value. An explicit prop always wins;
// context only fills empty; both empty yields "" and the component's
// omit-empty rendering applies.
func ResolveNonce(ctx context.Context, explicit string) string {
	if explicit != "" {
		return explicit
	}

	return NonceFromContext(ctx)
}
