package layout

import "github.com/larsartmann/templ-components/utils"

// StackGap is a typed enum for vertical spacing between stacked children.
// Unknown values fall back to StackGapMD. Mirrors the GridGap shape so the
// two composition primitives share a vocabulary.
type StackGap string

const (
	StackGapNone StackGap = "none"
	StackGapSM   StackGap = "sm" // space-y-2 (0.5rem)
	StackGapMD   StackGap = "md" // space-y-4 (1rem) — default
	StackGapLG   StackGap = "lg" // space-y-6 (1.5rem)
	StackGapXL   StackGap = "xl" // space-y-8 (2rem)
	// StackGapDefault is the canonical default (MD).
	StackGapDefault StackGap = StackGapMD
)

// StackGapIsValid reports whether g is a recognized StackGap.
func StackGapIsValid(g StackGap) bool {
	_, ok := stackGapLookup[g]

	return ok
}

//nolint:gochecknoglobals // Package-level lookup table for stack gaps
var stackGapLookup = map[StackGap]string{
	StackGapNone: "",
	StackGapSM:   "space-y-2",
	StackGapMD:   "space-y-4",
	StackGapLG:   "space-y-6",
	StackGapXL:   "space-y-8",
}

// stackGapClass returns the Tailwind space-y-* class for a gap. Empty string
// for None. Falls back to StackGapDefault for unknown values.
func stackGapClass(g StackGap) string {
	return utils.Lookup(stackGapLookup, g, stackGapLookup[StackGapDefault])
}

// StackProps configures a vertical stack of children with consistent spacing.
// Use Stack to replace repeated `space-y-N` strings — it is the 1D vertical
// counterpart to the 2D Grid primitive. Stack is flex-col (NOT grid): both
// axes do not matter, only the vertical rhythm. See ADR-0016.
type StackProps struct {
	utils.BaseProps

	// Gap controls the vertical space between children. Default StackGapMD
	// when zero or unknown.
	Gap StackGap
}

// DefaultStackProps returns sensible defaults: MD gap.
//
//nolint:exhaustruct // constructor intentionally sets only non-zero defaults
func DefaultStackProps() StackProps {
	return StackProps{
		Gap: StackGapDefault,
	}
}
