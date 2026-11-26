package layout

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// TestStylesheetInputsEmbedTemplates pins the vendor carrier: `go mod vendor`
// ships layout/styles/ because embed.go embeds it, so those files have to be
// byte-identical twins of templates/ and the embed has to cover exactly them.
// A drifted twin reaches a consumer's Tailwind build from vendor/ as a wrong
// stylesheet; a missing one arrives as an unresolved import.
func TestStylesheetInputsEmbedTemplates(t *testing.T) {
	t.Parallel()

	embedded := map[string]bool{}

	err := fs.WalkDir(stylesheetInputs, "styles", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}

		embedded[entry.Name()] = true

		return nil
	})
	if err != nil {
		t.Fatalf("walk the embedded stylesheet inputs: %v", err)
	}

	for _, name := range []string{"templ-components.css", "custom.css", "templ-components-theme.css", "safelist.css"} {
		twin, readErr := stylesheetInputs.ReadFile("styles/" + name)
		if readErr != nil {
			t.Errorf("layout/styles/%s is not embedded, so go mod vendor does not ship it: %v", name, readErr)

			continue
		}

		template, templateErr := os.ReadFile(filepath.Join("..", "templates", name))
		if templateErr != nil {
			t.Errorf("templates/%s unreadable: %v", name, templateErr)

			continue
		}

		if !bytes.Equal(twin, template) {
			t.Errorf("layout/styles/%s differs from templates/%s; re-copy it (go mod vendor ships these bytes)",
				name, name)
		}

		delete(embedded, name)
	}

	for stray := range embedded {
		t.Errorf("layout/styles/%s is embedded but is not one of the library's stylesheet inputs", stray)
	}
}
