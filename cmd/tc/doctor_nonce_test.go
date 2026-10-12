package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// writeFile writes content to path, creating parent directories.
func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// TestNonceConsumerDocsSync pins docs/nonce-consumers.md to the doctor's
// nonceConsumers scan table: every table entry must appear in the docs and
// vice versa. The docs are the human-readable twin of the scan list (F-058,
// issue #407) — the 2026-10-12 fleet audit shipped three repos of dead
// scripts because nobody knew which components consume a nonce.
func TestNonceConsumerDocsSync(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile(filepath.Join("..", "..", "docs", "nonce-consumers.md"))
	if err != nil {
		t.Fatalf("read docs/nonce-consumers.md: %v", err)
	}

	docs := string(data)

	if len(nonceConsumers) == 0 {
		t.Fatal("nonceConsumers table is empty — the doctor scan would be a no-op")
	}

	for _, consumer := range nonceConsumers {
		name := consumer.Qualified
		if !strings.Contains(docs, "`"+name+"`") {
			t.Errorf(
				"nonceConsumers entry %q missing from docs/nonce-consumers.md — update both in the same commit",
				name,
			)
		}
	}

	// Every backticked table row in the docs must be a scan entry (or the
	// documented prose exception), so the docs never advertise a component
	// the doctor does not scan.
	rowRe := regexp.MustCompile(`(?m)^\| ` + "`" + `([^` + "`" + `]+)` + "`")
	for _, m := range rowRe.FindAllStringSubmatch(docs, -1) {
		name := m[1]
		if name == "Nonce" || name == "tc doctor" || name == "TableRow.Href" ||
			name == "Dismissible" || name == "Eager" || name == "Content-Security-Policy" {
			continue
		}

		found := false

		for _, consumer := range nonceConsumers {
			if consumer.Qualified == name {
				found = true

				break
			}
		}

		if !found {
			t.Errorf("docs row %q is not in the nonceConsumers scan table — update both in the same commit", name)
		}
	}
}

// TestCheckNoncelessScripts exercises the doctor's nonce scan against a
// fixture module: empty-props-literal and fake-positional-literal call sites
// are detected; clean sources pass.
func TestCheckNoncelessScripts(t *testing.T) {
	t.Parallel()

	t.Run("detects empty nonce call sites", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/app\n\ngo 1.27\n")
		writeFile(t, filepath.Join(dir, "page.templ"), "package ui\n\n"+`
templ page() {
	@components.TCErrorHandling()
	@feedback.ToastContainer("")
	@layout.ThemeToggle("Toggle theme", "")
	@display.CopyButton(display.CopyButtonProps{ Nonce: "" })
}
`)

		findings := scanNonceFindings(dir)
		if len(findings) < 3 {
			t.Fatalf("expected >=3 findings, got %d: %v", len(findings), findings)
		}

		for _, f := range findings {
			if !strings.Contains(f, "page.templ:") {
				t.Errorf("finding not file-anchored: %q", f)
			}
		}
	})

	t.Run("clean sources pass", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "go.mod"), "module example.com/app\n\ngo 1.27\n")
		writeFile(t, filepath.Join(dir, "page.templ"), "package ui\n\n"+`
templ page(nonce string) {
	@feedback.ToastContainer(nonce)
	@layout.ThemeToggle("Toggle theme", nonce)
}
`)

		if findings := scanNonceFindings(dir); len(findings) != 0 {
			t.Fatalf("expected 0 findings, got %d: %v", len(findings), findings)
		}
	})

	t.Run("library package dirs are skipped", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "go.mod"), "module github.com/larsartmann/templ-components\n\ngo 1.27\n")
		writeFile(t, filepath.Join(dir, "display", "x.templ"), "package display\n\nNonce: \"\"\n")

		if findings := scanNonceFindings(dir); len(findings) != 0 {
			t.Fatalf("library dirs must be skipped, got %d findings", len(findings))
		}
	})
}
