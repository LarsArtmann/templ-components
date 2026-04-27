// Package layout provides layout components like Base, Theme, and SRI utilities.
package layout

func htmxIntegrityHash(version string) string {
	if h, ok := getHtmxSRIHashes()[version]; ok {
		return h.Main
	}
	return ""
}

func htmxResponseTargetsIntegrityHash(version string) string {
	if h, ok := getHtmxSRIHashes()[version]; ok {
		return h.ResponseTargets
	}
	return ""
}

func getHtmxSRIHashes() map[string]struct{ Main, ResponseTargets string } {
	return map[string]struct{ Main, ResponseTargets string }{
		"2.0.6": {
			Main:            "sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm",
			ResponseTargets: "sha384-FcXXcaqsB+SLXujBqU9KJ7E84XV/wxvVAMAGam/W56Y4g0mE9pgh4HG+A4IlfbNd",
		},
	}
}
