package icons

import "sort"

// Animation represents a hover-triggered animation style for animated icons.
// Inspired by https://www.heroicons-animated.com/ — pure CSS reimplementations
// of the Motion-based React originals (zero JavaScript required).
//
// Each animation plays on :hover and :focus-within, respects prefers-reduced-motion,
// and works with any icon via the AnimatedIcon / AnimatedIconWithAnimation components.
//
// Animation types are generalized presets that cover the major patterns from the
// original library (which has bespoke per-icon animations). Each icon is mapped
// to the preset that best matches its original Motion animation.
type Animation string

const (
	// AnimNone disables animation (explicit opt-out).
	AnimNone Animation = ""

	// AnimPulse — gentle scale burst.
	// Source (Heart): scale [1, 1.08, 1] x2, 0.45s ease-in-out.
	AnimPulse Animation = "pulse"

	// AnimBeat — strong scale with overshoot.
	// Source (Star): scale [1, 0.9, 1.2, 1], 0.6s ease-in-out.
	AnimBeat Animation = "beat"

	// AnimBounce — multi-direction translation bounce.
	// Source (MagnifyingGlass): x/y keyframe bounce, 1s ease-in-out.
	AnimBounce Animation = "bounce"

	// AnimWiggle — rotation shake.
	// Source (Bell): rotate [0, -10, 10, -10, 0], 0.5s ease-in-out.
	AnimWiggle Animation = "wiggle"

	// AnimSpin — one-shot rotation with spring-like easing.
	// Source (ArrowPath): rotate with spring (stiffness 250, damping 25).
	// Source (Cog6Tooth): rotate with spring.
	AnimSpin Animation = "spin"

	// AnimJump — scale + vertical lift.
	// Source (Home): scale [1, 1.1, 1] + y [0, -1, 0], 0.4s ease-out.
	AnimJump Animation = "jump"

	// AnimNod — vertical bob.
	// Source (ChevronDown): translateY [0, 2, 0], 0.5s ease-in-out.
	AnimNod Animation = "nod"

	// AnimShake — horizontal shift + rotation burst.
	// Source (Play): x [0, -1, 2, 0] + rotate [0, -10, 0, 0], 0.5s.
	// Source (LockClosed): rotate [-3, 2, -2, 1, 0] + scale.
	AnimShake Animation = "shake"

	// AnimBlink — per-path eyelid blink (requires 2-path icon).
	// Source (Eye): path1 scaleY [1, 0.1, 1] + opacity, path2 scale [1, 0.3, 1].
	// Only icons with 2+ SVG path elements can use this (currently: Eye, Settings,
	// Location/MapPin, Tag). Using it on a single-path icon is a no-op.
	AnimBlink Animation = "blink"

	// AnimWobble — scale down + rotation oscillation.
	// Source (Beaker): scale 0.9 (spring) + rotate [0, 6, -6, 3, -3, 0], 0.8s.
	AnimWobble Animation = "wobble"

	// AnimDraw — self-drawing stroke animation via stroke-dashoffset.
	// Source (Bolt): pathLength [0, 1], pathOffset [1, 0], opacity [0, 1], 0.6s linear.
	// Renders paths with pathLength="1" so CSS stroke-dashoffset works uniformly.
	AnimDraw Animation = "draw"
)

// validAnimations is the complete set of non-empty animation types.
//
//nolint:gochecknoglobals // Lookup table for animation validation
var validAnimations = map[Animation]bool{
	AnimPulse:  true,
	AnimBeat:   true,
	AnimBounce: true,
	AnimWiggle: true,
	AnimSpin:   true,
	AnimJump:   true,
	AnimNod:    true,
	AnimShake:  true,
	AnimBlink:  true,
	AnimWobble: true,
	AnimDraw:   true,
}

// IsValid reports whether a is a known animation type (excluding AnimNone).
// AnimNone is intentionally not "valid" — it is an explicit opt-out, not a
// named animation to validate against.
func (a Animation) IsValid() bool {
	return validAnimations[a]
}

// defaultAnimations maps each icon to its heroicons-animated-inspired default
// animation. Aliases (ArrowPath, Bars3, MapPin, HandThumbUp) are resolved to
// their canonical name in DefaultAnimation, so only canonical names appear here.
//
// Mappings are categorized as:
//   - "verified" — the icon exists in heroicons-animated and the animation type
//     faithfully recreates the original Motion variant
//   - "semantic" — the icon does not exist in heroicons-animated (or the original
//     uses a pattern not expressible in pure CSS), so the closest preset is chosen
//     based on the icon's visual meaning
//
//nolint:gochecknoglobals // Default animation lookup table
var defaultAnimations = map[Name]Animation{
	// --- Verified from heroicons-animated source ---
	Heart:    AnimPulse,  // scale [1, 1.08, 1] x2
	Star:     AnimBeat,   // scale [1, 0.9, 1.2, 1]
	Bell:     AnimWiggle, // rotate [0, -10, 10, -10, 0]
	Settings: AnimSpin,   // cog rotation (spring)
	Eye:      AnimBlink,  // per-path scaleY [1, 0.1, 1]
	Home:     AnimJump,   // scale + translateY
	Search:   AnimBounce, // x/y bounce
	Beaker:   AnimWobble, // scale 0.9 + rotate [0, 6, -6, 3, -3, 0]
	Bolt:     AnimDraw,   // pathLength [0, 1] self-draw
	Refresh:  AnimSpin,   // rotation with spring

	// --- Semantic: pulse (gentle confirmation/identity) ---
	Check:              AnimPulse,
	X:                  AnimPulse,
	Plus:               AnimPulse,
	Minus:              AnimPulse,
	Bookmark:           AnimPulse,
	ShieldCheck:        AnimPulse,
	CheckCircle:        AnimPulse,
	Information:        AnimPulse,
	UserCircle:         AnimPulse,
	UserPlus:           AnimPulse,
	XCircle:            AnimPulse,
	EllipsisHorizontal: AnimPulse,
	EllipsisVertical:   AnimPulse,
	DocumentDuplicate:  AnimPulse,
	Hashtag:            AnimPulse,
	Server:             AnimPulse,
	Signal:             AnimPulse,
	Squares2x2:         AnimPulse,
	FaceSmile:          AnimPulse,
	Camera:             AnimPulse,
	BuildingOffice2:    AnimPulse,
	Users:              AnimPulse,

	// --- Semantic: beat (strong attention/emphasis) ---
	ExclamationTriangle: AnimBeat,
	ExclamationCircle:   AnimBeat,
	Fire:                AnimBeat,

	// --- Semantic: wiggle (rotation-based playfulness) ---
	Share:         AnimWiggle,
	PaperAirplane: AnimWiggle,
	Gift:          AnimWiggle,
	Tag:           AnimWiggle,
	BellSlash:     AnimWiggle,
	Microphone:    AnimWiggle,
	BugAnt:        AnimWiggle,
	PuzzlePiece:   AnimWiggle,
	Question:      AnimWiggle,
	Trash:         AnimWiggle,

	// --- Semantic: spin (rotation/turning) ---
	Wrench: AnimSpin,
	Globe:  AnimSpin,
	Sun:    AnimSpin,
	Clock:  AnimSpin,
	Key:    AnimSpin,
	Cube:   AnimSpin,

	// --- Semantic: bounce (translate-based movement) ---
	Download:          AnimBounce,
	Upload:            AnimBounce,
	Mail:              AnimBounce,
	Location:          AnimBounce,
	ThumbUp:           AnimBounce,
	Cloud:             AnimBounce,
	Photo:             AnimBounce,
	ArrowDownOnSquare: AnimBounce,
	ArrowUpOnSquare:   AnimBounce,

	// --- Semantic: nod (vertical bob) ---
	ChevronDown:  AnimNod,
	ChevronUp:    AnimNod,
	ArrowUp:      AnimNod,
	ArrowDown:    AnimNod,
	Folder:       AnimNod,
	Document:     AnimNod,
	Chart:        AnimNod,
	Edit:         AnimNod,
	Menu:         AnimNod,
	Calendar:     AnimNod,
	Moon:         AnimNod,
	Clipboard:    AnimNod,
	CodeBracket:  AnimNod,
	DocumentText: AnimNod,
	Filter:       AnimNod,
	ListBullet:   AnimNod,
	Printer:      AnimNod,
	QueueList:    AnimNod,
	ArchiveBox:   AnimNod,
	Calculator:   AnimNod,
	FolderOpen:   AnimNod,
	AcademicCap:  AnimNod,
	Inbox:        AnimNod,

	// --- Semantic: shake (horizontal shift + rotation) ---
	ChevronRight:          AnimShake,
	ChevronLeft:           AnimShake,
	ArrowRight:            AnimShake,
	ArrowLeft:             AnimShake,
	ExternalLink:          AnimShake,
	Link:                  AnimShake,
	EyeOff:                AnimShake,
	Lock:                  AnimShake,
	Unlock:                AnimShake,
	Phone:                 AnimShake,
	NoSymbol:              AnimShake,
	ArrowRightOnRectangle: AnimShake,

	// --- Semantic: jump (scale + lift) ---
	RocketLaunch: AnimJump,
}

// DefaultAnimation returns the default animation for the given icon name.
// Icons not in the explicit mapping default to AnimPulse.
// Spinner always returns AnimNone (it has its own built-in spin animation).
// Aliases (ArrowPath, Bars3, MapPin, HandThumbUp) are resolved to their
// canonical name before lookup, so they inherit the canonical icon's animation.
func DefaultAnimation(name Name) Animation {
	if name == Spinner {
		return AnimNone
	}

	if canonical, ok := iconAliases[name]; ok {
		name = canonical
	}

	if anim, ok := defaultAnimations[name]; ok {
		return anim
	}

	return AnimPulse
}

// resolveAnimation applies safety guards to an animation request.
// Per-path animations (AnimBlink) that require multiple SVG path elements
// fall back to AnimPulse when the icon doesn't have enough paths,
// preventing silent no-ops where the CSS targets nth-child(N) but only
// one path exists.
func resolveAnimation(name Name, anim Animation) Animation {
	if anim == AnimBlink && len(iconPaths(name)) < 2 {
		return AnimPulse
	}

	return anim
}

// AllAnimations returns all valid animation types, sorted alphabetically.
// Useful for documentation, icon galleries, and demos.
func AllAnimations() []Animation {
	anims := make([]Animation, 0, len(validAnimations))
	for a := range validAnimations {
		anims = append(anims, a)
	}

	sort.Slice(anims, func(i, j int) bool {
		return string(anims[i]) < string(anims[j])
	})

	return anims
}
