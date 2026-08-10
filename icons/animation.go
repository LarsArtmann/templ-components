package icons

import "sort"

// Animation represents a hover-triggered animation style for animated icons.
// Inspired by https://www.heroicons-animated.com/ — pure CSS reimplementations
// of the Motion-based React originals (zero JavaScript required).
//
// Each animation plays on :hover and :focus-within, respects prefers-reduced-motion,
// and works with any icon via the AnimatedIcon / AnimatedIconWithAnimation components.
type Animation string

const (
	// AnimNone disables animation (explicit opt-out).
	AnimNone Animation = ""

	// AnimPulse — gentle scale burst (Heart style).
	// Source: scale [1, 1.08, 1] x2, 0.45s.
	AnimPulse Animation = "pulse"

	// AnimBeat — strong scale with overshoot (Star style).
	// Source: scale [1, 0.9, 1.2, 1], 0.6s.
	AnimBeat Animation = "beat"

	// AnimBounce — multi-direction bounce (Search style).
	// Source: x/y keyframe bounce, 1s.
	AnimBounce Animation = "bounce"

	// AnimWiggle — rotation shake (Bell style).
	// Source: rotate [0, -10, 10, -10, 0], 0.5s.
	AnimWiggle Animation = "wiggle"

	// AnimSpin — one-shot rotation with spring-like easing (Cog/Settings style).
	// Source: rotate 180deg, spring (stiffness 50, damping 10).
	AnimSpin Animation = "spin"

	// AnimJump — scale + vertical lift (Home style).
	// Source: scale [1, 1.1, 1] + y [0, -1, 0], 0.4s.
	AnimJump Animation = "jump"

	// AnimNod — vertical bob (Chevron style).
	// Source: translateY [0, 2, 0], 0.5s.
	AnimNod Animation = "nod"

	// AnimShake — horizontal shift + rotation burst (Play style).
	// Source: x [0, -1, 2, 0] + rotate [0, -10, 0, 0], 0.5s.
	AnimShake Animation = "shake"

	// AnimBlink — per-path eyelid blink (Eye style, requires 2-path icon).
	// Source: path1 scaleY [1, 0.1, 1] + opacity, path2 scale [1, 0.3, 1] + opacity.
	AnimBlink Animation = "blink"

	// AnimSplit — per-path lid/body split (Trash style, requires 2-path icon).
	// Source: path1 translateY -1.5px (lid up), path2 translateY 1px (body down).
	AnimSplit Animation = "split"
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
	AnimSplit:  true,
}

// IsValid reports whether a is a known animation type (excluding AnimNone).
// AnimNone is intentionally not "valid" — it is an explicit opt-out, not a
// named animation to validate against.
func (a Animation) IsValid() bool {
	return validAnimations[a]
}

// defaultAnimations maps each icon to its heroicons-animated-inspired default
// animation. Icons not listed here default to AnimPulse. The mapping is based
// on the animation semantics from the heroicons-animated source project.
//
//nolint:gochecknoglobals // Default animation lookup table
var defaultAnimations = map[Name]Animation{
	// Verified from heroicons-animated source
	Heart:    AnimPulse,
	Star:     AnimBeat,
	Bell:     AnimWiggle,
	Settings: AnimSpin,
	Eye:      AnimBlink,
	Trash:    AnimWiggle,
	Home:     AnimJump,
	Search:   AnimBounce,

	// Semantically assigned based on icon meaning
	EyeOff:              AnimShake,
	Wrench:              AnimSpin,
	Globe:               AnimSpin,
	Sun:                 AnimJump,
	Moon:                AnimJump,
	RocketLaunch:        AnimJump,
	Bolt:                AnimBeat,
	Fire:                AnimBeat,
	Check:               AnimPulse,
	Plus:                AnimPulse,
	Minus:               AnimPulse,
	Bookmark:            AnimPulse,
	Share:               AnimWiggle,
	PaperAirplane:       AnimWiggle,
	Gift:                AnimWiggle,
	ExternalLink:        AnimShake,
	Link:                AnimShake,
	Download:            AnimBounce,
	Upload:              AnimBounce,
	Mail:                AnimBounce,
	ChevronDown:         AnimNod,
	ChevronUp:           AnimNod,
	ArrowDown:           AnimNod,
	ArrowUp:             AnimNod,
	ChevronRight:        AnimShake,
	ChevronLeft:         AnimShake,
	ArrowRight:          AnimShake,
	ArrowLeft:           AnimShake,
	ExclamationTriangle: AnimBeat,
	BellSlash:           AnimWiggle,
	Tag:                 AnimWiggle,
}

// DefaultAnimation returns the default animation for the given icon name.
// Icons not in the explicit mapping default to AnimPulse.
// Spinner always returns AnimNone (it has its own built-in spin animation).
func DefaultAnimation(name Name) Animation {
	if name == Spinner {
		return AnimNone
	}

	if anim, ok := defaultAnimations[name]; ok {
		return anim
	}

	return AnimPulse
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
