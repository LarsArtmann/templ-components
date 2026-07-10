package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoJSONv2Imports prevents encoding/json/v2 and encoding/json/jsontext
// from being reintroduced. These packages require GOEXPERIMENT=jsonv2
// (not stable until Go 1.27) and break builds. BuildFlow's go-auto-upgrade
// tool can silently rewrite encoding/json to v2 — this test catches that.
func TestNoJSONv2Imports(t *testing.T) {
	t.Parallel()

	jsonPkg := "encoding/json"
	forbidden := []string{
		`"` + jsonPkg + `/v2"`,
		`"` + jsonPkg + `/jsontext"`,
	}

	skipDirs := map[string]bool{
		".git":         true,
		"vendor":       true,
		"node_modules": true,
	}

	var violations []string

	err := filepath.Walk("..", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if skipDirs[filepath.Base(path)] {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		data, readErr := os.ReadFile(path) //nolint:gosec // test scans source files
		if readErr != nil {
			return fmt.Errorf("read file: %w", readErr)
		}

		content := string(data)
		for _, imp := range forbidden {
			if strings.Contains(content, imp) {
				violations = append(violations, path+" contains "+imp)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}

	for _, v := range violations {
		t.Errorf("forbidden import: %s", v)
	}
	if len(violations) > 0 {
		t.Error("encoding/json/v2 requires GOEXPERIMENT=jsonv2 (Go 1.27+). Use encoding/json (v1).")
	}
}
