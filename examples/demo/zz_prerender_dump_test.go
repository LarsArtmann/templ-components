package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestZZDumpPrerenderDiff(t *testing.T) {
	dir := t.TempDir()

	if err := prerender(dir); err != nil {
		t.Fatalf("prerender: %v", err)
	}

	srv := httptest.NewServer(newDemoHandler())
	t.Cleanup(srv.Close)

	jar, jarErr := cookiejar.New(nil)
	if jarErr != nil {
		t.Fatalf("cookie jar: %v", jarErr)
	}
	client := &http.Client{Jar: jar}

	resp, err := client.Get(srv.URL + "/") //nolint:noctx // throwaway diagnostic
	if err != nil {
		t.Fatalf("fetch: %v", err)
	}
	live, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	pre, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		t.Fatalf("read prerender: %v", err)
	}

	outDir := "/tmp/prerender-diff"
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "pre.html"), []byte(normalizePrerender(string(pre))), 0o644); err != nil {
		t.Fatalf("write pre: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "live.html"), []byte(normalizePrerender(string(live))), 0o644); err != nil {
		t.Fatalf("write live: %v", err)
	}
	t.Logf("dumped to %s", outDir)
}
