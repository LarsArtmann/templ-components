package visualtest

import (
	"os"
	"path/filepath"
	"testing"
)

// TestVisualCoverage reports the ratio of golden PNGs to source components.
// This is informational — it surfaces coverage gaps without failing CI.
// Threshold can be tightened as coverage grows.
func TestVisualCoverage(t *testing.T) {
	t.Parallel()

	goldenCount := countGoldens(t, "testdata")
	componentCount := countComponents(t, "..")

	ratio := float64(goldenCount) / float64(componentCount)

	t.Logf(
		"visual coverage: %d goldens / %d components = %.1f%%",
		goldenCount,
		componentCount,
		ratio*100,
	)

	if ratio < 0.05 {
		t.Errorf(
			"visual coverage %.1f%% is below 5%% — add more goldens",
			ratio*100,
		)
	}
}

func countGoldens(t *testing.T, root string) int {
	t.Helper()

	count := 0

	err := filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".png" {
			count++
		}

		return nil
	})
	if err != nil {
		t.Logf("walk goldens: %v", err)
	}

	return count
}

func countComponents(t *testing.T, root string) int {
	t.Helper()

	packages := []string{
		"display", "errorpage", "feedback", "forms",
		"htmx", "layout", "navigation",
	}

	count := 0

	for _, pkg := range packages {
		files, err := filepath.Glob(filepath.Join(root, pkg, "*.templ"))
		if err != nil {
			continue
		}

		count += len(files)
	}

	return count
}
