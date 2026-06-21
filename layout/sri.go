// Package layout pins the HTMX CDN scripts and their Sub-Resource Integrity hashes.
package layout

// HTMX 2.x split its extensions into standalone npm packages with independent
// version numbers. The response-targets extension is no longer bundled inside
// the htmx.org package: loading htmx.org@<v>/dist/ext/response-targets.js
// emits the console warning "You are using an htmx 1 extension with htmx 2.0.x"
// because that file is the legacy htmx 1 build. The v2-compatible build ships
// as the separate package htmx-ext-response-targets, pinned below.
const (
	// defaultHTMXVersion is the htmx main-script version DefaultPageProps pins.
	defaultHTMXVersion = "2.0.10"
	// responseTargetsVersion is the pinned version of the separate
	// htmx-ext-response-targets package. It is independent of defaultHTMXVersion.
	responseTargetsVersion = "2.0.4"
)

const (
	htmxCDNBase = "https://unpkg.com/htmx.org@"
	// responseTargetsCDNURL is the full URL of the v2 response-targets minified
	// build. It is not derived from htmxCDNBase because the extension is its own
	// package now (see the comment on the version block above).
	responseTargetsCDNURL = "https://unpkg.com/htmx-ext-response-targets@" +
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
//	curl -sL https://unpkg.com/htmx.org@<v> | openssl dgst -sha384 -binary | openssl base64 -A
//
//nolint:gochecknoglobals // Version-keyed SRI lookup table.
var sriHTMXMainByVersion = map[string]string{
	defaultHTMXVersion: "sha384-H5SrcfygHmAuTDZphMHqBJLc3FhssKjG7w/CeCpFReSfwBWDTKpkzPP8c+cLsK+V",
}

// htmxScriptURL returns the CDN URL for the htmx main script at the given version.
func htmxScriptURL(version string) string {
	return htmxCDNBase + version
}

// htmxMainSRI returns the SRI hash for the htmx main script at the given
// version, or "" if the version is not pinned in sriHTMXMainByVersion.
func htmxMainSRI(version string) string {
	return sriHTMXMainByVersion[version]
}
