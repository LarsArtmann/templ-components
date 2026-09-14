package build

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInlineScriptHashes(t *testing.T) {
	t.Parallel()

	themeScript := "(function(){var theme;})()"
	jsonLD := `{"@context":"https://schema.org"}`

	pages := []RenderedPage{
		{Path: "index.html", HTML: `<html><head>` +
			`<script src="/assets/js/app.js" defer></script>` +
			`<script type="application/ld+json">` + jsonLD + `</script>` +
			`<script nonce="abc">` + themeScript + `</script>` +
			`</head></html>`},
		{Path: "docs.html", HTML: `<html><head>` +
			`<script type="application/ld+json">` + jsonLD + `</script>` +
			`<script nonce="def">` + themeScript + `</script>` +
			`</head></html>`},
	}

	hashes := InlineScriptHashes(pages)
	if len(hashes) != 2 {
		t.Fatalf(
			"expected 2 unique inline hashes (external src excluded, duplicates merged), got %d: %v",
			len(hashes),
			hashes,
		)
	}

	digest := sha256.Sum256([]byte(themeScript))
	wantTheme := "'sha256-" + base64.StdEncoding.EncodeToString(digest[:]) + "'"

	digest = sha256.Sum256([]byte(jsonLD))
	wantJSONLD := "'sha256-" + base64.StdEncoding.EncodeToString(digest[:]) + "'"

	if hashes[0] != wantJSONLD || hashes[1] != wantTheme {
		t.Errorf("unexpected hashes:\n got %v\nwant [%s %s]", hashes, wantJSONLD, wantTheme)
	}
}

func TestInlineScriptHashesStableAcrossNonces(t *testing.T) {
	t.Parallel()

	body := "if (!window.tcGoBackAttached) { history.back(); }"

	first := InlineScriptHashes([]RenderedPage{{Path: "a.html", HTML: `<script nonce="one">` + body + `</script>`}})
	second := InlineScriptHashes([]RenderedPage{{Path: "a.html", HTML: `<script nonce="two">` + body + `</script>`}})

	if len(first) != 1 || first[0] != second[0] {
		t.Fatalf("hash must be nonce-independent: %v vs %v", first, second)
	}
}

func TestCSPHeader(t *testing.T) {
	t.Parallel()

	empty := CSPHeader(nil)
	if !strings.Contains(empty, "script-src 'self';") || strings.Contains(empty, "  ") {
		t.Errorf("header without hashes malformed: %q", empty)
	}

	header := CSPHeader([]string{"'sha256-aaa='", "'sha256-bbb='"})

	for _, directive := range []string{
		"default-src 'self'",
		"script-src 'self' 'sha256-aaa=' 'sha256-bbb='",
		"style-src 'self'",
		"img-src 'self' data:",
		"font-src 'self'",
		"connect-src 'self'",
		"object-src 'none'",
		"base-uri 'self'",
		"form-action 'self'",
		"frame-ancestors 'none'",
	} {
		if !strings.Contains(header, directive) {
			t.Errorf("header %q missing directive %q", header, directive)
		}
	}
}

func TestSyncFirebaseCSP(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "firebase.json")

	fixture := `{
  "hosting": {
    "public": "dist",
    "cleanUrls": true,
    "headers": [
      {
        "source": "**/*.css",
        "headers": [{"key": "Cache-Control", "value": "immutable"}]
      },
      {
        "source": "**",
        "headers": [{"key": "X-Frame-Options", "value": "DENY"}]
      }
    ]
  }
}`
	if err := os.WriteFile(path, []byte(fixture), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	changed, err := SyncFirebaseCSP(path, "default-src 'self'")
	if err != nil {
		t.Fatalf("sync: %v", err)
	}

	if !changed {
		t.Error("first sync must report a change")
	}

	changed, err = SyncFirebaseCSP(path, "default-src 'self'")
	if err != nil {
		t.Fatalf("re-sync: %v", err)
	}

	if changed {
		t.Error("idempotent re-sync must not report a change")
	}

	if err := CheckFirebaseCSP(path, "default-src 'self'"); err != nil {
		t.Errorf("check after sync: %v", err)
	}
}

func TestCheckFirebaseCSPFailures(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "firebase.json")
	if err := os.WriteFile(
		path,
		[]byte(
			`{"hosting":{"headers":[{"source":"**","headers":[{"key":"Content-Security-Policy","value":"default-src 'self'"}]}]}}`,
		),
		0o600,
	); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := CheckFirebaseCSP(path, "other"); !errors.Is(err, errCSPHeaderStale) {
		t.Errorf("stale header must return errCSPHeaderStale, got %v", err)
	}

	missing := filepath.Join(t.TempDir(), "firebase.json")
	if err := os.WriteFile(missing, []byte(`{"hosting":{}}`), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	if err := CheckFirebaseCSP(missing, "x"); !errors.Is(err, errCSPHeaderMissing) {
		t.Errorf("missing header must return errCSPHeaderMissing, got %v", err)
	}

	if err := CheckFirebaseCSP(filepath.Join(t.TempDir(), "absent.json"), "x"); err == nil {
		t.Error("absent config file must error")
	}
}
