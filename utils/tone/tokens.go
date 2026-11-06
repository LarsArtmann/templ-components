package tone

// SurfaceToken names a surface role: which LAYER of the interface a
// background belongs to (canvas / panel / raised). Role-indexed, not
// shade-indexed, so a brand specifies its light+dark surface pairs as
// data: the token's CSS variable is defined in templates/custom.css and
// re-defined under `.dark`, making one class self-dark-aware — no
// dark: variant, no value duplication in Go (ADR-0049 §3).
//
// Class emission starts with the shell conversion (master plan S11);
// until then nothing renders these classes. The type exists now so the
// namespace contract is pinned before any component commits to it.
type SurfaceToken string

const (
	// SurfaceToken0 is the page canvas behind everything.
	SurfaceToken0 SurfaceToken = "surface-0"
	// SurfaceToken1 is a panel resting on the canvas (cards, tables, sheets).
	SurfaceToken1 SurfaceToken = "surface-1"
	// SurfaceToken2 is a raised layer on a panel (inputs, popovers, tooltips).
	SurfaceToken2 SurfaceToken = "surface-2"
)

// SurfaceTokenDefault is the canonical default (the canvas).
const SurfaceTokenDefault = SurfaceToken0

// SurfaceTokenIsValid reports whether t is a recognized SurfaceToken.
func SurfaceTokenIsValid(t SurfaceToken) bool {
	switch t {
	case SurfaceToken0, SurfaceToken1, SurfaceToken2:
		return true
	default:
		return false
	}
}

// SurfaceTokens returns every SurfaceToken, in stacking order (canvas up).
func SurfaceTokens() []SurfaceToken {
	return []SurfaceToken{SurfaceToken0, SurfaceToken1, SurfaceToken2}
}

// InkToken names a text-ink role: how strongly text reads against the
// surface it sits on. Same variable-flip contract as SurfaceToken
// (ADR-0049 §3): the class is self-dark-aware, values are brand data.
//
// Class emission starts with the surface/ink sweeps (master plan S11);
// the type pins the namespace contract now.
type InkToken string

const (
	// InkToken1 is primary text (headings, body emphasis).
	InkToken1 InkToken = "ink-1"
	// InkToken2 is secondary text (body copy, descriptions).
	InkToken2 InkToken = "ink-2"
	// InkToken3 is muted text (meta, placeholders, footnotes).
	InkToken3 InkToken = "ink-3"
)

// InkTokenDefault is the canonical default (primary).
const InkTokenDefault = InkToken1

// InkTokenIsValid reports whether t is a recognized InkToken.
func InkTokenIsValid(t InkToken) bool {
	switch t {
	case InkToken1, InkToken2, InkToken3:
		return true
	default:
		return false
	}
}

// InkTokens returns every InkToken, in emphasis order (strong to muted).
func InkTokens() []InkToken {
	return []InkToken{InkToken1, InkToken2, InkToken3}
}
