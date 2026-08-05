// Package cdn provides shared CDN URL helpers for runtime libraries
// (HTMX, Datastar, etc.). Both layout and datastar need the same logic:
// default a base URL when empty, trim trailing slashes, and extract the
// scheme+host (origin) for <link rel="preconnect">. Centralizing the
// helpers prevents drift between the two packages.
package cdn

import "strings"

// ResolveBase returns cdnBase if non-empty, otherwise defaultBase. A
// trailing slash is trimmed so consumers can pass "https://unpkg.com/"
// safely. The fallback default is the package default (e.g. jsDelivr).
func ResolveBase(cdnBase, defaultBase string) string {
	if cdnBase == "" {
		cdnBase = defaultBase
	}

	return strings.TrimRight(cdnBase, "/")
}

// Origin extracts the scheme+host (origin) from the CDN base URL for use
// in <link rel="preconnect">. Returns "" if the input doesn't parse as
// an absolute URL (e.g. a relative self-hosted path like "/assets").
//
// The scheme prefix is fixed at 8 characters to match both "http://"
// and "https://" — both are 7 chars plus the colon-slash-slash.
//
//	cdn.Origin("https://cdn.jsdelivr.net/npm") -> "https://cdn.jsdelivr.net"
//	cdn.Origin("https://unpkg.com/foo")        -> "https://unpkg.com"
//	cdn.Origin("/assets")                      -> ""
func Origin(cdnBase, defaultBase string) string {
	base := ResolveBase(cdnBase, defaultBase)
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		return ""
	}

	idx := strings.Index(base[8:], "/")
	if idx < 0 {
		return base
	}

	return base[:8+idx]
}
