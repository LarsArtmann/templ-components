package visualtest

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// compiledCSS caches the project's compiled Tailwind CSS so every test renders
// with the exact same stylesheet consumers see. It is read once from
// examples/demo/static/app.css (the demo's compiled output) relative to the
// repo root.
var (
	cssOnce sync.Once
	cssData []byte
	cssErr  error
)

// loadCSS returns the compiled Tailwind stylesheet. It resolves the path from
// this file's location so it works regardless of the test working directory.
func loadCSS() ([]byte, error) {
	cssOnce.Do(func() {
		path := cssPath()

		cssData, cssErr = os.ReadFile(path)
		if cssErr == nil {
			return
		}

		cssErr = &cssLoadError{path: path, err: cssErr}
	})

	return cssData, cssErr
}

// cssPath locates examples/demo/static/app.css relative to this source file.
// runtime.Caller gives the absolute path of this file in the module tree, so
// the lookup is immune to the process working directory.
func cssPath() string {
	_, file, _, _ := runtime.Caller(0) // visualtest/css.go
	// file = <repo>/visualtest/css.go
	repoRoot := filepath.Dir(file)

	return filepath.Join(repoRoot, "..", "examples", "demo", "static", "app.css")
}

type cssLoadError struct {
	path string
	err  error
}

func (e *cssLoadError) Error() string {
	return "visualtest: read compiled CSS at " + e.path + ": " + e.err.Error()
}

func (e *cssLoadError) Unwrap() error { return e.err }
