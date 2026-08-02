package datastar

import "strings"

// DatastarVersion is a pinned Datastar runtime version. Use the exported
// constants (e.g. DatastarVersion1_0_2) for compile-time safety; custom
// versions can be constructed via DatastarVersion("x.y.z") but will render
// without SRI (no integrity hash — Datastar does not publish official SRI
// hashes; self-host and compute your own if integrity is required).
type DatastarVersion string

const (
	// DatastarVersion1_0_2 is the Datastar runtime version this package pins.
	DatastarVersion1_0_2 DatastarVersion = "1.0.2"

	// defaultDatastarVersion is the internal default; always equals the latest
	// exported DatastarVersion constant.
	defaultDatastarVersion DatastarVersion = DatastarVersion1_0_2
)

// Datastar CDN paths. Datastar is served from jsDelivr's GitHub path
// (not npm), so the URL structure differs from HTMX's npm-style path.
const (
	// defaultDatastarCDNBase is the CDN base URL. Consumers override via
	// SDKScriptProps.CDN (e.g. "https://unpkg.com" or a self-hosted origin).
	defaultDatastarCDNBase = "https://cdn.jsdelivr.net/gh"

	// datastarRepoPath is the GitHub org/repo segment in the jsDelivr URL.
	datastarRepoPath = "starfederation/datastar"

	// datastarBundlePath is the path to the bundled runtime inside the repo.
	datastarBundlePath = "bundles/datastar.js"
)

// resolveDatastarCDN returns cdnBase if non-empty, otherwise defaultDatastarCDNBase.
// A trailing slash is trimmed so consumers can pass "https://unpkg.com/" safely.
func resolveDatastarCDN(cdnBase string) string {
	if cdnBase == "" {
		cdnBase = defaultDatastarCDNBase
	}

	return strings.TrimRight(cdnBase, "/")
}

// datastarScriptURL returns the CDN URL for the Datastar runtime at the given
// version. If cdnBase is empty, defaults to defaultDatastarCDNBase.
//
// Example: https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.2/bundles/datastar.js.
func datastarScriptURL(version DatastarVersion, cdnBase string) string {
	v := version
	if v == "" {
		v = defaultDatastarVersion
	}

	return resolveDatastarCDN(cdnBase) + "/" + datastarRepoPath + "@" + string(v) + "/" + datastarBundlePath
}

// datastarCDNOrigin extracts the scheme+host (origin) from the CDN base URL for
// use in <link rel="preconnect">. Returns "" if the input doesn't parse as
// an absolute URL (e.g. a relative self-hosted path like "/assets").
func datastarCDNOrigin(cdnBase string) string {
	base := resolveDatastarCDN(cdnBase)
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		return ""
	}

	idx := strings.Index(base[8:], "/")
	if idx < 0 {
		return base
	}

	return base[:8+idx]
}
