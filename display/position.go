package display

// Position defines where a floating element appears relative to its trigger,
// expressed in logical (writing-direction aware) sides: start and end map to
// the inline start/end of the surrounding text direction (left/right in LTR,
// mirrored in RTL). It is the shared placement axis for Tooltip, Popover, and
// HoverCard.
type Position string

const (
	PositionTop    Position = "top"
	PositionBottom Position = "bottom"
	PositionStart  Position = "start"
	PositionEnd    Position = "end"
)

// Legacy physical spellings from before the logical axis existed. They keep
// their original physical behavior (never mirrored in RTL) and stay valid so
// existing consumers keep compiling; prefer PositionStart/PositionEnd, which
// are pixel-identical to them in LTR and mirror correctly in RTL. Removal is
// planned for v2 (ADR-0022).
const (
	PositionLeft  Position = "left"
	PositionRight Position = "right"
)

// PositionIsValid reports whether v is one of the defined Position constants.
func PositionIsValid(v Position) bool {
	switch v {
	case PositionTop, PositionBottom, PositionStart, PositionEnd, PositionLeft, PositionRight:
		return true
	default:
		return false
	}
}

// Component aliases of [Position]: one type serves all floating components;
// the per-component spellings remain so call sites read naturally and
// pre-unification source keeps compiling.
type (
	// TooltipPosition is an alias of [Position]: where a Tooltip appears
	// relative to its trigger.
	TooltipPosition = Position

	// PopoverPosition is an alias of [Position]: where a Popover panel
	// appears relative to its trigger.
	PopoverPosition = Position

	// HoverCardPosition is an alias of [Position]: where a HoverCard appears
	// relative to its trigger.
	HoverCardPosition = Position
)

// Legacy per-component constants, kept for source compatibility. The
// Left/Right spellings predate the logical axis; prefer Start/End in new
// code (identical placement in LTR, correct mirroring in RTL).
const (
	TooltipPositionTop    = PositionTop
	TooltipPositionBottom = PositionBottom
	TooltipPositionStart  = PositionStart
	TooltipPositionEnd    = PositionEnd
	TooltipPositionLeft   = PositionLeft
	TooltipPositionRight  = PositionRight

	PopoverPositionTop    = PositionTop
	PopoverPositionBottom = PositionBottom
	PopoverPositionStart  = PositionStart
	PopoverPositionEnd    = PositionEnd
	PopoverPositionLeft   = PositionLeft
	PopoverPositionRight  = PositionRight

	HoverCardPositionTop    = PositionTop
	HoverCardPositionBottom = PositionBottom
	HoverCardPositionStart  = PositionStart
	HoverCardPositionEnd    = PositionEnd
)

// TooltipPositionIsValid reports whether v is one of the defined Position constants.
func TooltipPositionIsValid(v TooltipPosition) bool {
	return PositionIsValid(v)
}

// PopoverPositionIsValid reports whether v is one of the defined Position constants.
func PopoverPositionIsValid(v PopoverPosition) bool {
	return PositionIsValid(v)
}

// HoverCardPositionIsValid reports whether v is one of the defined Position constants.
func HoverCardPositionIsValid(v HoverCardPosition) bool {
	return PositionIsValid(v)
}
