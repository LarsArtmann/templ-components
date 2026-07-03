package layout

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
	"net/http"
	"testing"
	"time"
)

// TestPinnedSRIMatchesCDN fetches each pinned CDN script and verifies that
// the live bytes match the hardcoded SRI hash. This catches version bumps
// that forget to update the hash, and supply-chain drift on the CDN.
//
// Skipped under -short and on network errors (offline, CDN down, etc.) so
// it never causes spurious CI failures — it is a safety net, not a gate.
func TestPinnedSRIMatchesCDN(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping CDN SRI verification in -short mode")
	}

	cases := []struct {
		name string
		url  string
		want string // full SRI string, e.g. "sha384-..."
	}{
		{
			name: "htmx main 2.0.10",
			url:  htmxScriptURL(defaultHTMXVersion, ""),
			want: htmxMainSRIDefault,
		},
		{
			name: "response-targets 2.0.4",
			url:  responseTargetsURL(""),
			want: sriResponseTargets,
		},
	}

	client := &http.Client{Timeout: 15 * time.Second}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Get(tc.url)
			if err != nil {
				t.Skipf("network error fetching %s: %v", tc.url, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Skipf("non-200 status (%d) fetching %s", resp.StatusCode, tc.url)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Skipf("error reading body from %s: %v", tc.url, err)
			}

			// SRI uses sha384 (part of sha512 family in Go's crypto).
			hash := sha512.Sum384(body)
			got := "sha384-" + base64.StdEncoding.EncodeToString(hash[:])

			if got != tc.want {
				t.Errorf("SRI mismatch for %s\n  URL:  %s\n  want: %s\n  got:  %s", tc.name, tc.url, tc.want, got)
			}
		})
	}
}
