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
	// Source (MagnifyingGlass): x/y keyframe bounce.
	// Source (ChevronRight): translateX keyframes.
	AnimBounce Animation = "bounce"

	// AnimWiggle — rotation shake.
	// Source (Bell): rotate [0, -10, 10, -10, 0], 0.5s ease-in-out.
	// Source (Moon): rotate [0, -10, 10, -5, 5, 0] — same oscillation family.
	AnimWiggle Animation = "wiggle"

	// AnimSpin — one-shot rotation with spring-like easing.
	// Source (ArrowPath): rotate with spring (stiffness 250, damping 25).
	// Source (Cog6Tooth): rotate with spring.
	AnimSpin Animation = "spin"

	// AnimJump — scale + vertical lift.
	// No icon defaults to jump; available via AnimatedIconWithAnimation.
	// Closest source pattern (Home): scale [1, 1.1, 1] + y [0, -1, 0].
	AnimJump Animation = "jump"

	// AnimNod — vertical bob.
	// Source (ChevronDown): translateY [0, 2, 0], 0.5s ease-in-out.
	AnimNod Animation = "nod"

	// AnimShake — horizontal shift + rotation burst.
	// Source (AcademicCap): rotate [0, -5, 5, 0] + translateY.
	// Source (BugAnt): rotate [0, -2, 2, -2, 2, 0].
	// Source (Key): rotate [0, 3, -3, 0].
	AnimShake Animation = "shake"

	// AnimBlink — per-path eyelid blink (requires 2+ SVG path elements).
	// Source (Eye): path1 scale [1, 0.3, 1] + opacity, path2 scale [1, 0.1, 1].
	// resolveAnimation falls back to AnimPulse for single-path icons.
	AnimBlink Animation = "blink"

	// AnimWobble — scale down + rotation oscillation.
	// Source (Beaker): scale 0.9 + rotate [0, 6, -6, 3, -3, 0].
	// Source (Cube): rotateY [0, 15, -15, 0].
	// Source (Lock): rotate [-3, 2, -2, 1, 0] + scale.
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
// All 96 icons were verified on 2026-08-11 against heroicons-animated source
// (https://github.com/heroicons-animated/heroicons-animated, 316 animated icons).
// 95 of 96 icons have a direct source equivalent; the last (ArrowRightOnRectangle)
// uses a semantic mapping. Each comment records the heroicons-animated filename
// and the Motion variant properties that informed the classification.
//
// Icons whose source uses per-path scaleY/scaleX (blink) but have only one SVG
// path in our implementation are adapted to the closest non-per-path equivalent.
//
//nolint:gochecknoglobals // Default animation lookup table
var defaultAnimations = map[Name]Animation{
	// --- draw (self-drawing stroke animation, 20 icons) ---
	Bolt:         AnimDraw, // HA=bolt: pathLength+opacity+pathOffset
	Check:        AnimDraw, // HA=check: pathLength+opacity+scale
	CheckCircle:  AnimDraw, // HA=check-circle: pathLength+opacity
	CodeBracket:  AnimDraw, // HA=code-bracket: pathLength+opacity+pathOffset
	DocumentText: AnimDraw, // HA=document-text: pathLength+opacity+pathMorph
	FaceSmile:    AnimDraw, // HA=face-smile: pathLength+opacity+rotate+scale
	Globe:        AnimDraw, // HA=globe-americas: pathLength+opacity
	Hashtag:      AnimDraw, // HA=hashtag: pathLength+opacity+pathOffset
	Inbox:        AnimDraw, // HA=inbox: pathLength+opacity
	Link:         AnimDraw, // HA=link: pathLength+pathOffset+rotate
	ListBullet:   AnimDraw, // HA=list-bullet: pathLength+opacity+translateY
	Minus:        AnimDraw, // HA=minus: pathLength+opacity
	NoSymbol:     AnimDraw, // HA=no-symbol: pathLength+opacity+pathOffset+scale
	PuzzlePiece:  AnimDraw, // HA=puzzle-piece: pathLength+opacity+pathOffset
	Share:        AnimDraw, // HA=share: pathLength+opacity+scale (lines draw + nodes pulse)
	ShieldCheck:  AnimDraw, // HA=shield-check: pathLength+opacity
	UserCircle:   AnimDraw, // HA=user-circle: pathLength+opacity+pathOffset
	UserPlus:     AnimDraw, // HA=user-plus: pathLength+opacity
	X:            AnimDraw, // HA=x-mark: pathLength+opacity
	XCircle:      AnimDraw, // HA=x-circle: pathLength+opacity

	// --- blink (per-path eyelid animation, requires 2+ paths, 1 icon) ---
	Eye: AnimBlink, // HA=eye: scale [1, 0.3, 1] per-path

	// --- pulse (gentle scale burst, 19 icons) ---
	Bookmark:            AnimPulse, // HA=bookmark: source=blink (scaleX/scaleY) but single-path → pulse
	BuildingOffice2:     AnimPulse, // HA=building-office-2: opacity
	Calendar:            AnimPulse, // HA=calendar: scale [1, 1.1, 1] + translateY
	Camera:              AnimPulse, // HA=camera: scale [1, 0.95, 1]
	Chart:               AnimPulse, // HA=chart-bar: source=blink (opacity+scaleY) but single-path → pulse
	ExclamationCircle:   AnimPulse, // HA=exclamation-circle: scale [1, 1.1, 1]
	ExclamationTriangle: AnimPulse, // HA=exclamation-triangle: scale [1, 1.1, 1]
	ExternalLink:        AnimPulse, // HA=arrow-top-right-on-square: scale+translateX+translateY
	Filter:              AnimPulse, // HA=funnel: source=blink (scaleX/scaleY) but single-path → pulse
	Heart:               AnimPulse, // HA=heart: scale [1, 1.08, 1]
	Home:                AnimPulse, // HA=home: scale [1, 1.1, 1] + translateY
	Information:         AnimPulse, // HA=information-circle: scale [1, 1.1, 1]
	Menu:                AnimPulse, // HA=bars-3: pathMorph+scaleX
	PaperAirplane:       AnimPulse, // HA=paper-airplane: scale [1, 0.8, 1]
	Question:            AnimPulse, // HA=question-mark-circle: scale [1, 1.1, 1]
	Signal:              AnimPulse, // HA=signal: opacity+scale
	Squares2x2:          AnimPulse, // HA=squares-2x2: scale [0.6, 1]
	Sun:                 AnimPulse, // HA=sun: per-ray opacity stagger
	Unlock:              AnimPulse, // HA=lock-open: scale [1, 1.05, 1]

	// --- beat (strong scale with overshoot, 4 icons) ---
	Calculator:         AnimBeat, // HA=calculator: scale [1, 1.5, 1]
	EllipsisHorizontal: AnimBeat, // HA=ellipsis-horizontal: scale [1, 1.3, 1]
	EllipsisVertical:   AnimBeat, // HA=ellipsis-vertical: scale [1, 1.3, 1]
	Star:               AnimBeat, // HA=star: scale [1, 0.9, 1.2, 1]

	// --- bounce (multi-direction translation, 12 icons) ---
	ArrowLeft:         AnimBounce, // HA=arrow-left: translateX+pathMorph
	ArrowRight:        AnimBounce, // HA=arrow-right: translateX+pathMorph
	BellSlash:         AnimBounce, // HA=bell-slash: translateX
	ChevronLeft:       AnimBounce, // HA=chevron-left: translateX
	ChevronRight:      AnimBounce, // HA=chevron-right: translateX
	Cloud:             AnimBounce, // HA=cloud: translateX+translateY
	DocumentDuplicate: AnimBounce, // HA=document-duplicate: opacity+translateX+translateY
	EyeOff:            AnimBounce, // HA=eye-slash: translateX
	Location:          AnimBounce, // HA=map-pin: translateY
	RocketLaunch:      AnimBounce, // HA=rocket-launch: pathMorph+translateX+translateY
	Search:            AnimBounce, // HA=magnifying-glass: translateX+translateY
	Users:             AnimBounce, // HA=users: opacity+translateX

	// --- wiggle (rotation oscillation, 11 icons) ---
	Bell:       AnimWiggle, // HA=bell: rotate [0, -10, 10, -10, 0]
	Edit:       AnimWiggle, // HA=pencil: rotate [0, -6, 6, 0]
	Fire:       AnimWiggle, // HA=fire: rotate [0, -3, 3, -2, 2, 0]
	FolderOpen: AnimWiggle, // HA=folder-open: rotate [0, -8, 6, -4, 0]
	Mail:       AnimWiggle, // HA=envelope: rotate [0, -5, 5, -3, 3, 0]
	Moon:       AnimWiggle, // HA=moon: rotate [0, -10, 10, -5, 5, 0]
	Phone:      AnimWiggle, // HA=phone: rotate [0, -10, 10, -10, 10, -5, 5, 0]
	Photo:      AnimWiggle, // HA=photo: rotate [0, -3, 3, -2, 2, 0]
	Tag:        AnimWiggle, // HA=tag: rotate [0, -10, 8, -5, 3, 0]
	ThumbUp:    AnimWiggle, // HA=hand-thumb-up: rotate [0, -10, 5, 0]
	Wrench:     AnimWiggle, // HA=wrench: rotate [0, 12, -14, 4, 0]

	// --- wobble (scale down + rotation oscillation, 4 icons) ---
	Beaker: AnimWobble, // HA=beaker: scale 0.9 + rotate [0, 6, -6, 3, -3, 0]
	Cube:   AnimWobble, // HA=cube: rotateY [0, 15, -15, 0]
	Gift:   AnimWobble, // HA=gift: rotate [0, -5, 5, -3, 3, 0] + scale [1, 1.05, 1]
	Lock:   AnimWobble, // HA=lock-closed: rotate [-3, 2, -2, 1, 0] + scale [1, 1.02, 0.98, 1]

	// --- spin (one-shot rotation, 4 icons) ---
	Clock:    AnimSpin, // HA=clock: rotate (spring)
	Plus:     AnimSpin, // HA=plus: rotate (spring)
	Refresh:  AnimSpin, // HA=arrow-path: rotate (spring)
	Settings: AnimSpin, // HA=cog-6-tooth: rotate (spring)

	// --- shake (horizontal shift + rotation burst, 4 icons) ---
	AcademicCap:           AnimShake, // HA=academic-cap: rotate [0, -5, 5, 0] + translateY
	ArrowRightOnRectangle: AnimShake, // No HA equivalent — semantic (directional arrow)
	BugAnt:                AnimShake, // HA=bug-ant: rotate [0, -2, 2, -2, 2, 0]
	Key:                   AnimShake, // HA=key: rotate [0, 3, -3, 0]

	// --- nod (vertical bob, 17 icons) ---
	ArchiveBox:        AnimNod, // HA=archive-box: translateY
	ArrowDown:         AnimNod, // HA=arrow-down: translateY+pathMorph
	ArrowDownOnSquare: AnimNod, // HA=arrow-down-on-square: translateY
	ArrowUp:           AnimNod, // HA=arrow-up: translateY+pathMorph
	ArrowUpOnSquare:   AnimNod, // HA=arrow-up-on-square: translateY
	ChevronDown:       AnimNod, // HA=chevron-down: translateY
	ChevronUp:         AnimNod, // HA=chevron-up: translateY
	Clipboard:         AnimNod, // HA=clipboard: source=blink (scaleY+translateY) but single-path → nod
	Document:          AnimNod, // HA=document: translateY
	Download:          AnimNod, // HA=arrow-down-tray: translateY
	Folder:            AnimNod, // HA=folder: translateY
	Microphone:        AnimNod, // HA=microphone: translateY
	Printer:           AnimNod, // HA=printer: translateY
	QueueList:         AnimNod, // HA=queue-list: opacity+translateY
	Server:            AnimNod, // HA=server: opacity+translateY
	Trash:             AnimNod, // HA=trash: translateY
	Upload:            AnimNod, // HA=arrow-up-tray: translateY
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
