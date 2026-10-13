package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestNonceConsumerIntegrationSync pins the doctor's nonceConsumers table to
// the integration nonce render tables (third leg of the #407 guard triad:
// doctor ⇔ docs is TestNonceConsumerDocsSync, doctor ⇔ integration is this).
//
// Direction 1 — every doctor entry's component name must appear in at least
// one integration render table (TestAllInlineScriptsHaveNonce's with-nonce
// direction or TestNoEmptyNonceAttribute's omit-empty direction), so no
// scannable component escapes the library's own nonce verification.
//
// Direction 2 — every integration table entry must map back to a doctor
// entry or a documented exception, so the render tables never verify a
// component the doctor cannot scan for.
//
// Integration entry names may carry variant suffixes (KanbanBoardWired,
// TableWithRowHref, PolledRegionEager); a doctor entry matches when the
// integration name starts with the doctor's short name.
func TestNonceConsumerIntegrationSync(t *testing.T) {
	t.Parallel()

	integrationNames := parseIntegrationRenderTableNames(t, filepath.Join("..", "..", "integration"))

	if len(integrationNames) == 0 {
		t.Fatal("no integration render-table entries found — the tables moved or were renamed")
	}

	// Direction 1: doctor → integration.
	for _, consumer := range nonceConsumers {
		short := consumer.Qualified[strings.LastIndex(consumer.Qualified, ".")+1:]

		if !slicesContainsPrefix(integrationNames, short) {
			t.Errorf(
				"nonceConsumers entry %q has no integration render-table coverage (expected an entry starting with %q in integration/csp_nonce_test.go or integration/nonce_omit_empty_test.go) — add it to both tables in the same commit",
				consumer.Qualified, short,
			)
		}
	}

	// Direction 2: integration → doctor.
	//
	// Documented exceptions: entries that are deliberately not doctor-scan
	// targets — variant renderings of a scanned component, cross-module
	// runtime emitters the doctor reaches through their component wrapper,
	// and positional-arg components covered in docs/nonce-consumers.md prose.
	exceptions := map[string]string{
		"DatastarSDKScript":  "runtime loader; consumers reach it via components/layout.templ (docs prose)",
		"EChartsSDKScript":   "runtime loader twin of the scanned echarts.EChart",
		"ViewTransitions":    "config script without a Nonce prop (htmx global setting)",
		"SSEErrorHandling":   "datastar module twin of the scanned htmx.GlobalErrorHandling",
		"DatastarLiveRegion": "datastar module; documented in docs/nonce-consumers.md prose",
		"MobileMenu":         "positional-arg component; documented prose exception (see nonceConsumers comment)",
		"ErrorAlert":         "renders via the scanned errorpage script path",
		"Script":             "external-script helper; nonce optional by design (src= needs no nonce)",
		"BaseSelfHostedHTMX": "layout.Base whole-page rendering (ThemeScript + self-host htmx inside)",
		"TagsInput":          "",
		"PolledRegionEager":  "",
		"TableWithRowHref":   "",
		"KanbanBoardWired":   "",
	}

	for _, name := range integrationNames {
		if slicesHasDoctorEntry(nonceConsumers, name) {
			continue
		}

		if _, ok := exceptions[name]; !ok {
			t.Errorf(
				"integration render-table entry %q has no nonceConsumers scan entry — add it to cmd/tc/doctor.go's table (and docs/nonce-consumers.md) or to this test's exception list with a reason",
				name,
			)
		}
	}
}

// parseIntegrationRenderTableNames extracts the entry names from the
// `{"Name", ...}` rows of both integration nonce render tables.
func parseIntegrationRenderTableNames(t *testing.T, integrationDir string) []string {
	t.Helper()

	files := []string{"csp_nonce_test.go", "nonce_omit_empty_test.go"}
	entryRe := regexp.MustCompile(`(?m)^\s*\{"([A-Za-z0-9]+)",`)
	seen := map[string]bool{}

	var names []string

	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(integrationDir, f))
		if err != nil {
			t.Fatalf("read integration/%s: %v", f, err)
		}

		for _, m := range entryRe.FindAllStringSubmatch(string(data), -1) {
			if !seen[m[1]] {
				seen[m[1]] = true
				names = append(names, m[1])
			}
		}
	}

	return names
}

// slicesHasDoctorEntry reports whether any doctor entry's short name is a
// prefix of name (integration entries may carry variant suffixes).
func slicesHasDoctorEntry(consumers []struct {
	Qualified string
	Note      string
}, name string) bool {
	for _, c := range consumers {
		short := c.Qualified[strings.LastIndex(c.Qualified, ".")+1:]
		if strings.HasPrefix(name, short) {
			return true
		}
	}

	return false
}

// slicesContainsPrefix reports whether any entry starts with prefix.
func slicesContainsPrefix(names []string, prefix string) bool {
	for _, n := range names {
		if strings.HasPrefix(n, prefix) {
			return true
		}
	}

	return false
}
