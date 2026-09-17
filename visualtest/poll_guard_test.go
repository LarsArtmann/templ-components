package visualtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoRawChromedpPoll bans raw chromedp.Poll calls in this module (TODO
// #240): a bool-valued expression polled into a *string always fails JSON
// unmarshal and the caller retries forever — the mistake is only writable
// through the raw API. Every poll must go through pollBool / pollTrue /
// pollText, which force the correct expression/capture pairing (poll.go is
// the one allowed implementation site).
func TestNoRawChromedpPoll(t *testing.T) {
	t.Parallel()

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if name := d.Name(); name == "testdata" || name == ".git" {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(path, ".go") || path == "poll.go" {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(string(src), "chromedp.Poll"+"(") {
			t.Errorf(
				"%s uses raw chromedp.Poll — use pollBool/pollTrue (booleans) or pollText (strings); see poll.go and the TODO #193 lesson",
				path,
			)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk visualtest module: %v", err)
	}
}
