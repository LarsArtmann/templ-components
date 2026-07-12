// HTMX CDN version pinning and SRI hash lookup.
package layout

import "strings"

// HTMXVersion is a pinned HTMX main-script version. Use the exported constants
// (e.g. HTMXVersion2_0_10) for compile-time safety; custom versions can be
// constructed via HTMXVersion("x.y.z") but will render without SRI (no integrity hash).
type HTMXVersion string

const (
	// HTMXVersion2_0_10 is the htmx main-script version DefaultPageProps pins.
	HTMXVersion2_0_10 HTMXVersion = "2.0.10"

	// defaultHTMXVersion is the internal default; always equals the latest
	// exported HTMXVersion constant.
	defaultHTMXVersion HTMXVersion = HTMXVersion2_0_10
)

// HTMX 2.x split its extensions into standalone npm packages with independent
// version numbers. The response-targets extension is no longer bundled inside
// the htmx.org package: loading htmx.org@<v>/dist/ext/response-targets.js
// emits the console warning "You are using an htmx 1 extension with htmx 2.0.x"
// because that file is the legacy htmx 1 build. The v2-compatible build ships
// as the separate package htmx-ext-response-targets, pinned below.
const (
	// responseTargetsVersion is the pinned version of the separate
	// htmx-ext-response-targets package. It is independent of defaultHTMXVersion.
	responseTargetsVersion = "2.0.4"
)

// defaultCDNBase is the CDN base URL used when PageProps.HTMXCDN is empty.
// Both htmx.org and htmx-ext-response-targets are served from npm-style paths
// under this base. Consumers override via PageProps.HTMXCDN (e.g.
// "https://unpkg.com" or a self-hosted origin).
const defaultCDNBase = "https://cdn.jsdelivr.net/npm"

// sriResponseTargets is the sha384 (base64) SRI hash for the pinned
// response-targets extension build.
const sriResponseTargets = "sha384-T41oglUPvXLGBVyRdZsVRxNWnOOqCynaPubjUVjxhsjFTKrFJGEMm3/0KGmNQ+Pg"

// sriHTMXMainByVersion maps an htmx main-script version to its SRI hash.
// Only versions listed here render an integrity attribute; a caller who pins a
// different version via PageProps.HTMXVersion gets the script without SRI.
//
// Regenerate a hash after a bump with:
//
//	curl -sL https://cdn.jsdelivr.net/npm/htmx.org@<v> | openssl dgst -sha384 -binary | openssl base64 -A
//
// sriHTMXMainDefault is the sha384 SRI hash for the pinned default htmx version.
const sriHTMXMainDefault = "sha384-H5SrcfygHmAuTDZphMHqBJLc3FhssKjG7w/CeCpFReSfwBWDTKpkzPP8c+cLsK+V"

//nolint:gochecknoglobals // Version-keyed SRI lookup table.
var sriHTMXMainByVersion = map[HTMXVersion]string{
	HTMXVersion2_0_10: sriHTMXMainDefault,
}

// resolveCDNBase returns cdnBase if non-empty, otherwise defaultCDNBase.
// A trailing slash is trimmed so consumers can pass "https://unpkg.com/" safely.
func resolveCDNBase(cdnBase string) string {
	if cdnBase == "" {
		cdnBase = defaultCDNBase
	}
	return strings.TrimRight(cdnBase, "/")
}

// htmxCDNOrigin extracts the scheme+host (origin) from the CDN base URL for
// use in <link rel="preconnect">. Returns "" if the input doesn't parse as
// an absolute URL (e.g. a relative self-hosted path like "/assets").
func htmxCDNOrigin(cdnBase string) string {
	base := resolveCDNBase(cdnBase)
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		return ""
	}
	idx := strings.Index(base[8:], "/")
	if idx < 0 {
		return base
	}
	return base[:8+idx]
}

// htmxScriptURL returns the CDN URL for the htmx main script at the given
// version. If cdnBase is empty, defaults to defaultCDNBase.
func htmxScriptURL(version HTMXVersion, cdnBase string) string {
	return resolveCDNBase(cdnBase) + "/htmx.org@" + string(version)
}

// responseTargetsURL returns the CDN URL for the response-targets extension.
// If cdnBase is empty, defaults to defaultCDNBase.
func responseTargetsURL(cdnBase string) string {
	return resolveCDNBase(cdnBase) + "/htmx-ext-response-targets@" +
		responseTargetsVersion + "/dist/response-targets.min.js"
}

// htmxMainSRI returns the SRI hash for the htmx main script at the given
// version. Returns empty string for unpinned versions — the browser will
// load the script without SRI verification, which is safer than returning
// a wrong hash (the browser would block the script on hash mismatch).
func htmxMainSRI(version HTMXVersion) string {
	return sriHTMXMainByVersion[version]
}
