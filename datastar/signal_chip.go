package datastar

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/larsartmann/templ-components/utils"
)

// SignalChipItem is one chip in a SignalChipGroup.
type SignalChipItem struct {
	// Label is the visible chip text.
	Label string

	// Value is what the group's signal becomes when the chip is clicked, and
	// the value the chip is compared against to decide whether it is pressed.
	Value string

	// Count, when non-nil, renders a numeric badge after the label.
	Count *int
}

// SignalChipProps configures one signal chip: a toggle button that assigns a
// value to a Datastar signal on click and swaps its own styling and
// aria-pressed state as that signal changes. It is the signal-driven twin of
// the zero-JS forms.FilterChips for pages that filter client-side.
//
//	@datastar.SignalChip(datastar.SignalChipProps{
//	    Signal:   "stateFilter",
//	    Value:    "open",
//	    Label:    "Open",
//	    Selected: true,
//	})
type SignalChipProps struct {
	utils.BaseProps

	// Signal is the signal name without the $ prefix (for example
	// "stateFilter" or "filters.state"). Clicking the chip assigns Value to
	// it. An empty or invalid name (anything but dot-separated identifier
	// segments) renders a static chip with no data-* bindings.
	Signal string

	// Value is assigned to Signal on click and is the chip's comparison
	// target: the chip is pressed while the signal equals Value. An empty
	// Value is a valid "all" option.
	Value string

	// Label is the visible chip text.
	Label string

	// Count, when non-nil, renders a numeric badge after the label.
	Count *int

	// Selected marks the chip as pressed at first paint, before and without
	// JavaScript. A selected chip also declares the signal's initial value
	// (data-signals) so the runtime starts in the state the server rendered;
	// an undeclared signal reads as empty and the runtime would strip the
	// pressed look on load.
	Selected bool
}

// DefaultSignalChipProps returns sensible defaults for a signal chip.
func DefaultSignalChipProps() SignalChipProps {
	return SignalChipProps{} //nolint:exhaustruct_v5 // intentionally minimal defaults
}

// SignalChipGroupProps configures a single-select row of signal chips that
// all write the same signal.
//
//	@datastar.SignalChipGroup(datastar.SignalChipGroupProps{
//	    AriaLabel: "Filter by state",
//	    Signal:    "stateFilter",
//	    Selected:  "open",
//	    Items: []datastar.SignalChipItem{
//	        {Label: "All"},
//	        {Label: "Open", Value: "open", Count: &openCount},
//	        {Label: "Done", Value: "done"},
//	    },
//	})
type SignalChipGroupProps struct {
	utils.BaseProps

	// Signal is the signal name every chip assigns to (see SignalChipProps).
	Signal string

	// Items render left to right. Provide an empty-Value item so users can
	// clear the filter.
	Items []SignalChipItem

	// Selected is the Value of the chip that renders pressed at first paint.
	// It is compared to each item's Value, so an empty Selected presses an
	// item whose Value is empty.
	Selected string
}

// DefaultSignalChipGroupProps returns sensible defaults for a signal chip group.
func DefaultSignalChipGroupProps() SignalChipGroupProps {
	return SignalChipGroupProps{} //nolint:exhaustruct_v5 // intentionally minimal defaults
}

// signalChipGroupLabel names a chip row whose consumer set neither AriaLabel
// nor a visible heading. Chip rows are filter groups and must be named.
const signalChipGroupLabel = "Filters"

// signalChipStateClass returns the look matching the server-rendered state.
func signalChipStateClass(selected bool) string {
	return utils.Ternary(selected, signalChipSelectedClass, signalChipUnselectedClass)
}

// signalChipPressedValue returns the server-rendered aria-pressed token.
func signalChipPressedValue(selected bool) string {
	return utils.Ternary(selected, "true", "false")
}

// signalChipSignalGrammar is the signal grammar a chip can bind: dot-separated
// JavaScript identifier segments. The name is interpolated into EXPRESSION
// values (never attribute names), so camelCase is fine, but anything outside
// this grammar could terminate the expression or address an unintended signal.
// Dotted names address nested signals.
var signalChipSignalGrammar = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)*$`)

// signalChipIsLive reports whether the signal name can be bound. Anything
// else renders a static chip (the graceful-degradation rule).
func signalChipIsLive(signal string) bool {
	return signalChipSignalGrammar.MatchString(signal)
}

// signalChipRef returns the expression reference to the signal ($name).
// Only valid for a name that passed signalChipIsLive.
func signalChipRef(signal string) string {
	return "$" + signal
}

// signalChipClickExpr returns the click expression that assigns value to the
// signal: $name='value'.
func signalChipClickExpr(signal, value string) string {
	return signalChipRef(signal) + "=" + jsStringLiteral(value)
}

// signalChipIsExpr returns the equality test used by every reactive binding:
// $name==='value'.
func signalChipIsExpr(signal, value string) string {
	return signalChipRef(signal) + "===" + jsStringLiteral(value)
}

// signalChipClassExpr returns the data-class expression that swaps the chip
// between its pressed and unpressed looks as the signal changes. The runtime's
// class plugin expects an OBJECT of class names to booleans: a plain class
// string is iterated character by character (each index becomes a class).
func signalChipClassExpr(signal, value string) string {
	notExpr := signalChipRef(signal) + "!==" + jsStringLiteral(value)

	return "{ " + jsStringLiteral(signalChipSelectedClass) + ": " + signalChipIsExpr(signal, value) +
		", " + jsStringLiteral(signalChipUnselectedClass) + ": " + notExpr + " }"
}

// signalChipPressedExpr returns the data-attr:aria-pressed expression. It
// yields the STRING 'true' or 'false' on purpose: a boolean expression makes
// the runtime toggle the attribute off when false (and set an empty value
// when true), which leaves assistive tech without a pressed state.
func signalChipPressedExpr(signal, value string) string {
	return signalChipIsExpr(signal, value) + " ? 'true' : 'false'"
}

// signalChipSignalsExpr returns the data-signals expression that declares the
// signal's initial value. A dotted name becomes a nested object literal.
func signalChipSignalsExpr(signal, value string) string {
	var b strings.Builder

	segments := strings.Split(signal, ".")

	for _, segment := range segments {
		b.WriteString("{")
		b.WriteString(jsStringLiteral(segment))
		b.WriteString(":")
	}

	b.WriteString(jsStringLiteral(value))

	for range segments {
		b.WriteString("}")
	}

	return b.String()
}

// signalChipItemProps builds the chip props for one item of a group.
func signalChipItemProps(group SignalChipGroupProps, item SignalChipItem) SignalChipProps {
	return SignalChipProps{ //nolint:exhaustruct_v5 // chips carry no BaseProps of their own
		Signal:   group.Signal,
		Value:    item.Value,
		Label:    item.Label,
		Count:    item.Count,
		Selected: item.Value == group.Selected,
	}
}

// jsStringLiteral renders s as a single-quoted JavaScript string literal.
// Backslashes and quotes are escaped before anything else, so no input can
// terminate the literal. Control characters, the Unicode line and paragraph
// separators, and invalid UTF-8 bytes use short escapes; everything else is
// passed through verbatim.
func jsStringLiteral(s string) string {
	var b strings.Builder

	b.Grow(len(s) + 2)
	b.WriteByte('\'')

	for i := 0; i < len(s); {
		c := s[i]
		if c < utf8.RuneSelf {
			i++

			switch c {
			case '\\', '\'':
				b.WriteByte('\\')
				b.WriteByte(c)
			case '\n':
				b.WriteString(`\n`)
			case '\r':
				b.WriteString(`\r`)
			default:
				if c < 0x20 || c == 0x7f {
					writeHexEscape(&b, c)

					continue
				}

				b.WriteByte(c)
			}

			continue
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == utf8.RuneError && size == 1:
			writeHexEscape(&b, c)
		case r == '\u2028':
			b.WriteString(`\u2028`)
		case r == '\u2029':
			b.WriteString(`\u2029`)
		default:
			b.WriteString(s[i : i+size])
		}

		i += size
	}

	b.WriteByte('\'')

	return b.String()
}

// writeHexEscape appends \xNN for one non-printable byte.
func writeHexEscape(b *strings.Builder, c byte) {
	const hexDigits = "0123456789abcdef"

	b.WriteString(`\x`)
	b.WriteByte(hexDigits[c>>4])
	b.WriteByte(hexDigits[c&0x0f])
}
