// Package layout provides Sub-Resource Integrity hashes for HTMX CDN scripts.
package layout

// htmxSRIEntry holds the SRI hashes for a single HTMX version
type htmxSRIEntry struct {
	Main            string
	ResponseTargets string
}

//nolint:gochecknoglobals // Package-level lookup table for HTMX SRI hashes
var htmxSRIHashes = map[string]htmxSRIEntry{
	"2.0.6": {
		Main:            "sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm",
		ResponseTargets: "sha384-FcXXcaqsB+SLXujBqU9KJ7E84XV/wxvVAMAGam/W56Y4g0mE9pgh4HG+A4IlfbNd",
	},
}

const htmxCDNBase = "https://unpkg.com/htmx.org@"

func htmxCDNURL(version, path string) string {
	url := htmxCDNBase + version
	if path != "" {
		url += path
	}
	return url
}

// htmxSRI returns the SRI hash for the given version and extension.
// ext should be "main" or "response-targets".
func htmxSRI(version, ext string) string {
	h, ok := htmxSRIHashes[version]
	if !ok {
		return ""
	}
	switch ext {
	case "response-targets":
		return h.ResponseTargets
	default:
		return h.Main
	}
}
