// CodeBlock component: a styled code display with an integrated copy
// button — the shared shape behind "copyable command" rows (mr-sync's
// commandLine) and detail-page ID rows (DiscordSync's copyIDRow), extracted
// per the 2026-10-10 four-project analysis (TODO_LIST #389).
package display

import (
	"github.com/larsartmann/templ-components/utils"
)

// CodeBlockVariant selects the rendering shape. Structural variants use an
// if-branch in the template (map lookups are for class data, not DOM shape).
type CodeBlockVariant string

// Code block variant constants.
const (
	// CodeBlockBlock renders a figure with an optional header (language or
	// label) and a scrollable pre/code body. The zero value.
	CodeBlockBlock CodeBlockVariant = "block"
	// CodeBlockCompactID renders an inline monospace row with a copy button —
	// the detail-page "copy this ID" pattern.
	CodeBlockCompactID CodeBlockVariant = "compact-id"
)

// CodeBlockVariantIsValid reports whether v is one of the defined variants.
// Unknown values render the block variant (graceful fallback, never a panic).
func CodeBlockVariantIsValid(v CodeBlockVariant) bool {
	return v == CodeBlockBlock || v == CodeBlockCompactID
}

// CodeBlockProps configures a code block.
type CodeBlockProps struct {
	utils.BaseProps

	// Code is the text to display and copy. Rendered HTML-escaped by templ.
	Code string
	// Language optionally labels the block's language in the header
	// ("bash", "go"). It is a label only — CodeBlock does not syntax-highlight.
	Language string
	// Label optionally overrides the header label. Defaults to Language
	// when Language is set and Label is empty; both empty renders no header
	// label (a copy-only header when copying is enabled).
	Label string
	// NoCopy hides the copy button. The zero value shows it — both extraction
	// demand shapes (copyable command rows, ID rows) always carry copy.
	NoCopy bool
	// Variant selects block (zero value) or compact-id rendering.
	Variant CodeBlockVariant
}

// DefaultCodeBlockProps returns sensible defaults.
func DefaultCodeBlockProps() CodeBlockProps {
	return CodeBlockProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Language: "bash",
	}
}

// codeBlockLabel resolves the header label: explicit Label wins, then
// Language, then none.
func codeBlockLabel(props CodeBlockProps) string {
	if props.Label != "" {
		return props.Label
	}

	return props.Language
}
