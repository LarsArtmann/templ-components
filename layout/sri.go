// Package layout pins the HTMX CDN scripts and their Sub-Resource Integrity hashes.
package layout

// HTMXVersion is a pinned HTMX main-script version. Use the exported constants
// (e.g. HTMXVersion2_0_10) for compile-time safety; custom versions can be
// constructed via HTMXVersion("x.y.z") but will fall back to the default SRI hash.
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

const (
	htmxCDNBase = "https://cdn.jsdelivr.net/npm/htmx.org@"
	// responseTargetsCDNURL is the full URL of the v2 response-targets minified
	// build. It is not derived from htmxCDNBase because the extension is its own
	// package now (see the comment on the version block above).
	responseTargetsCDNURL = "https://cdn.jsdelivr.net/npm/htmx-ext-response-targets@" +
		responseTargetsVersion + "/dist/response-targets.min.js"
)

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
// htmxMainSRIDefault is the sha384 SRI hash for the pinned default htmx version.
const htmxMainSRIDefault = "sha384-H5SrcfygHmAuTDZphMHqBJLc3FhssKjG7w/CeCpFReSfwBWDTKpkzPP8c+cLsK+V"

//nolint:gochecknoglobals // Version-keyed SRI lookup table.
var sriHTMXMainByVersion = map[HTMXVersion]string{
	HTMXVersion2_0_10: htmxMainSRIDefault,
}

// htmxScriptURL returns the CDN URL for the htmx main script at the given version.
func htmxScriptURL(version HTMXVersion) string {
	return htmxCDNBase + string(version)
}

// htmxMainSRI returns the SRI hash for the htmx main script at the given
// version. If the version is not pinned, it falls back to the default version's
// SRI so HTMXUseSRI never silently renders a script without an integrity hash
// (an unknown version bump shouldn't silently drop SRI).
func htmxMainSRI(version HTMXVersion) string {
	if sri, ok := sriHTMXMainByVersion[version]; ok {
		return sri
	}
	return sriHTMXMainByVersion[HTMXVersion2_0_10]
}
