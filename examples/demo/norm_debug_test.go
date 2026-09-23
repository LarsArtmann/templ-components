package main

import (
	"net/http"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDebugNorm(t *testing.T) {
	dir := t.TempDir()
	if err := prerender(dir); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(newMux())
	defer srv.Close()
	pre, _ := os.ReadFile(filepath.Join(dir, "index.html"))
	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	live, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	a := normalizePrerender(string(pre))
	b := normalizePrerender(string(live))
	for d := range firstDiffs(a, b, 3) {
		lo := max(0, d-80)
		t.Logf("diff@%d\nA:...%q\nB:...%q", d, a[lo:min(d+80, len(a))], b[lo:min(d+80, len(b))])
	}
	if a != b {
		t.Fail()
	}
}

func firstDiffs(a, b string, n int) map[int]bool {
	out := map[int]bool{}
	found := 0
	for i := 0; i < len(a) && i < len(b) && found < n; i++ {
		if a[i] != b[i] {
			out[i] = true
			found++
		}
	}
	return out
}
