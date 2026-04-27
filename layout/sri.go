package layout

var htmxSRIHashes = map[string]struct{ Main, ResponseTargets string }{
	"2.0.6": {
		Main:            "sha384-Akqfrbj/HpNVo8k11SXBb6TlBWmXXlYQrCSqEWmyKJe+hDm3Z/B2WVG4smwBkRVm",
		ResponseTargets: "sha384-FcXXcaqsB+SLXujBqU9KJ7E84XV/wxvVAMAGam/W56Y4g0mE9pgh4HG+A4IlfbNd",
	},
}

func htmxIntegrityHash(version string) string {
	if h, ok := htmxSRIHashes[version]; ok {
		return h.Main
	}
	return ""
}

func htmxResponseTargetsIntegrityHash(version string) string {
	if h, ok := htmxSRIHashes[version]; ok {
		return h.ResponseTargets
	}
	return ""
}
