package navigation

import (
	"strings"
	"testing"
)

// TestBreadcrumbJSONLDNoScriptBreakout pins the JSON-LD break-out closure:
// item names are operator-controlled text injected into a CSP-exempt
// `<script type="application/ld+json">` block, so a literal "</" in a name
// would end the block early. The class is closed by breadcrumbJSONLD's
// explicit "</" → "<\\/" escape (json/v2 does NOT HTML-escape by default,
// unlike v1 Marshal); this pin keeps it closed against encoder changes.
func TestBreadcrumbJSONLDNoScriptBreakout(t *testing.T) {
	t.Parallel()

	malicious := []BreadcrumbItem{
		{Text: `Home</script><script>alert(1)</script>`, Href: "/"},
		{Text: `Detail & <b>bold</b>`, Href: "/detail"},
	}

	out := breadcrumbJSONLD(malicious, "https://example.com")

	// The break-out sequence is exactly "</": a bare "<" or "&" elsewhere is
	// inert JSON payload text. The escaped form is "<\/" (a valid JSON
	// forward-slash escape), so no "</" may survive into the payload.
	if strings.Contains(out, "</") {
		t.Fatalf("JSON-LD output contains a literal </ sequence (break-out class reopened): %s", out)
	}

	if strings.Contains(out, "</script>") {
		t.Fatalf("close tag survived into the JSON-LD payload: %s", out)
	}
}
