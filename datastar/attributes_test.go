package datastar

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// TestSignalNameIsValid pins the attribute-suffix-safe signal grammar
// (lowercase segments, dotted paths, `_`-local prefix, `__` banned — the
// modifier delimiter). Same-commit rule: no IsValid without a test.
func TestSignalNameIsValid(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		valid bool
	}{
		{"users", true},
		{"metrics.users", true},
		{"_local_draft", true},
		{"metrics._transient", true},
		{"user2.name_2", true},
		{"", false},
		{"Users", false},               // attribute names lowercase — uppercase would silently mismatch
		{"camelCase", false},           // use VALUE-borne data-signals JSON for camelCase keys
		{"metrics__deep", false},       // __ is the modifier delimiter
		{"a..b", false},                // empty segment
		{".users", false},              // leading empty segment
		{"users.", false},              // trailing empty segment
		{"has space", false},           // no spaces
		{"has-dash", false},            // dashes are the attribute key delimiter
		{"javascript:alert(1)", false}, // scheme injection
	} {
		if got := SignalName(tt.name).IsValid(); got != tt.valid {
			t.Errorf("SignalName(%q).IsValid() = %v, want %v", tt.name, got, tt.valid)
		}
	}
}

// TestBuildersDegradeToNil pins the graceful-degradation rule: a builder
// without its inputs renders nothing instead of a broken binding.
func TestBuildersDegradeToNil(t *testing.T) {
	t.Parallel()

	for name, attrs := range map[string]templ.Attributes{
		"Show":        Show(""),
		"Text":        Text(""),
		"Bind":        Bind(""),
		"BindInvalid": Bind("Users"),
		"Computed":    Computed("total", ""),
		"Signals":     Signals(nil),
		"On":          On("", "$count++").Attrs(),
		"OnExpr":      On("click", "").Attrs(),
		"OnBadEvent":  On(`click'`, "$count++").Attrs(),
		"OnNoMods":    On("click", "$x").Attrs(),
	} {
		if name == "OnNoMods" {
			if attrs == nil {
				t.Errorf("On with valid inputs must render")

				continue
			}

			continue
		}

		if attrs != nil {
			t.Errorf("%s with empty/invalid inputs must degrade to nil, got %#v", name, attrs)
		}
	}
}

// TestBuilderSpelling pins the pinned-runtime attribute spellings: plugin
// keys, modifier groups after __ (never a dotted suffix — the
// data-on:click.prevent never-fires trap), and the debounce duration unit.
func TestBuilderSpelling(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		attrs     templ.Attributes
		wantKey   string
		wantValue string
	}{
		{"show", Show("$visible"), "data-show", "$visible"},
		{"text", Text("$user.name"), "data-text", "$user.name"},
		{"bind", Bind("email"), "data-bind:email", ""},
		{"computed", Computed("total", "$price * $qty"), "data-computed:total", "$price * $qty"},
		{"on-plain", On("click", "$count++").Attrs(), "data-on:click", "$count++"},
		{"on-prevent", On("click", "@get('/x')").Prevent().Attrs(), "data-on:click__prevent", "@get('/x')"},
		{"on-debounce", On("input", "@get('/s')").Debounce(300).Attrs(), "data-on:input__debounce.300ms", "@get('/s')"},
		{"on-chain", On("submit", "@post('/f')").Prevent().Debounce(120).Throttle(80).Attrs(),
			"data-on:submit__prevent__debounce.120ms__throttle.80ms", "@post('/f')"},
	}

	for _, tt := range cases {
		value, ok := tt.attrs[tt.wantKey]
		if !ok {
			t.Errorf("%s: missing key %q in %#v", tt.name, tt.wantKey, tt.attrs)

			continue
		}

		if value != tt.wantValue {
			t.Errorf("%s: %q = %q, want %q", tt.name, tt.wantKey, value, tt.wantValue)
		}

		if len(tt.attrs) != 1 {
			t.Errorf("%s: want exactly one attribute, got %#v", tt.name, tt.attrs)
		}
	}
}

// TestSignalsBuilder pins the JSON-object form (value-borne keys, so
// camelCase is safe here) and the JSON-marshallable contract.
func TestSignalsBuilder(t *testing.T) {
	t.Parallel()

	attrs := Signals(map[string]any{"metrics": map[string]any{"users": 5, "camelsRock": true}})
	raw, ok := attrs["data-signals"]
	if !ok {
		t.Fatalf("data-signals missing: %#v", attrs)
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(raw.(string)), &decoded); err != nil {
		t.Fatalf("data-signals must be valid JSON: %v (%q)", err, raw)
	}

	if !strings.Contains(raw.(string), `"camelsRock":true`) {
		t.Errorf("camelCase keys must survive verbatim: %q", raw)
	}

	// Non-marshalable payloads degrade to nil (channels, funcs).
	if Signals(map[string]any{"bad": make(chan int)}) != nil {
		t.Error("non-marshalable Signals payload must degrade to nil")
	}
}

// TestNormalizeExpr pins the expression-value contract: CR/LF collapse to
// spaces (one statement line), while quotes and angle brackets pass through
// VERBATIM — they are expression code (string literals like @get('/x'),
// comparisons) that the runtime parses; the HTML attribute layer is escaped
// by templ at render time. Escaping quotes here would MANGLE every
// datastar.Get-constructed expression.
func TestNormalizeExpr(t *testing.T) {
	t.Parallel()

	expr := Get("/api/x")

	if got := normalizeExpr(expr); got != expr {
		t.Errorf("normalizeExpr must leave expression code verbatim: got %q", got)
	}

	if got := normalizeExpr("a\nb\rc"); got != "a b c" {
		t.Errorf("normalizeExpr must collapse CR/LF to spaces: got %q", got)
	}
}
