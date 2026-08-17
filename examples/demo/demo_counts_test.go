package main

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

// TestHeroCountsMatchFeatures guards the demo hero stat cards against drift:
// the hardcoded component/package counts must equal the canonical totals in
// FEATURES.md (which itself is drift-guard verified). The icon count is
// computed dynamically from icons.AllIconNames() and is asserted against
// FEATURES.md so both stay honest.
func TestHeroCountsMatchFeatures(t *testing.T) {
	features, err := os.ReadFile(filepath.Join("..", "..", "FEATURES.md"))
	if err != nil {
		t.Skipf("FEATURES.md not readable: %v", err)
	}

	totals := regexp.MustCompile(`\*\*Totals:\*\* (\d+) templ components.*?, (\d+) icon names`).
		FindStringSubmatch(string(features))
	if totals == nil {
		t.Fatal("could not parse FEATURES.md Totals line")
	}

	if componentCount != totals[1] {
		t.Errorf("demo componentCount = %s, FEATURES.md totals = %s (update the const in demo.templ)",
			componentCount, totals[1])
	}

	wantIcons, err := strconv.Atoi(totals[2])
	if err != nil {
		t.Fatalf("parse icon count: %v", err)
	}
	if got := len(icons.AllIconNames()); got != wantIcons {
		t.Errorf("icons.AllIconNames() = %d icons, FEATURES.md = %d", got, wantIcons)
	}

	// The rendered hero must carry the real library version, not a stale one.
	var buf strings.Builder
	if err := demoHero().Render(context.Background(), &buf); err != nil {
		t.Fatalf("render demoHero: %v", err)
	}
	if !strings.Contains(buf.String(), "v"+utils.Version) {
		t.Errorf("hero does not render current utils.Version (%s)", utils.Version)
	}
}
