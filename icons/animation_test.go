package icons

import (
	"slices"
	"testing"
)

func TestAnimationIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		anim Animation
		want bool
	}{
		{"pulse", AnimPulse, true},
		{"beat", AnimBeat, true},
		{"bounce", AnimBounce, true},
		{"wiggle", AnimWiggle, true},
		{"spin", AnimSpin, true},
		{"jump", AnimJump, true},
		{"nod", AnimNod, true},
		{"shake", AnimShake, true},
		{"blink", AnimBlink, true},
		{"wobble", AnimWobble, true},
		{"draw", AnimDraw, true},
		{"none is not valid", AnimNone, false},
		{"empty is not valid", Animation(""), false},
		{"bogus is not valid", Animation("bogus"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.anim.IsValid(); got != tt.want {
				t.Errorf("Animation(%q).IsValid() = %v, want %v", tt.anim, got, tt.want)
			}
		})
	}
}

func TestDefaultAnimation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		icon Name
		want Animation
	}{
		// Verified against heroicons-animated source (2026-08-11, 316 icons)
		{"Heart defaults to pulse", Heart, AnimPulse},
		{"Star defaults to beat", Star, AnimBeat},
		{"Bell defaults to wiggle", Bell, AnimWiggle},
		{"Settings defaults to spin", Settings, AnimSpin},
		{"Eye defaults to blink", Eye, AnimBlink},
		{"Home defaults to pulse", Home, AnimPulse},
		{"Search defaults to bounce", Search, AnimBounce},
		{"Beaker defaults to wobble", Beaker, AnimWobble},
		{"Bolt defaults to draw", Bolt, AnimDraw},
		{"Refresh defaults to spin", Refresh, AnimSpin},
		{"Moon defaults to wiggle", Moon, AnimWiggle},
		{"Sun defaults to pulse", Sun, AnimPulse},
		{"Lock defaults to wobble", Lock, AnimWobble},
		{"Trash defaults to nod", Trash, AnimNod},
		{"Check defaults to draw", Check, AnimDraw},
		{"X defaults to draw", X, AnimDraw},
		{"ChevronRight defaults to bounce", ChevronRight, AnimBounce},
		{"ArrowRight defaults to bounce", ArrowRight, AnimBounce},
		{"Wrench defaults to wiggle", Wrench, AnimWiggle},
		{"Cube defaults to wobble", Cube, AnimWobble},
		{"Calculator defaults to beat", Calculator, AnimBeat},
		{"Fire defaults to wiggle", Fire, AnimWiggle},
		{"AcademicCap defaults to shake", AcademicCap, AnimShake},

		// Spinner always opts out
		{"Spinner defaults to none", Spinner, AnimNone},

		// Aliases resolve to canonical icon's animation
		{"ArrowPath alias resolves to spin", ArrowPath, AnimSpin},
		{"Bars3 alias resolves to Menu pulse", Bars3, AnimPulse},
		{"HandThumbUp alias resolves to ThumbUp wiggle", HandThumbUp, AnimWiggle},
		{"MapPin alias resolves to Location bounce", MapPin, AnimBounce},
		{"Close shares X draw", Close, AnimDraw},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := DefaultAnimation(tt.icon); got != tt.want {
				t.Errorf("DefaultAnimation(%q) = %q, want %q", tt.icon, got, tt.want)
			}
		})
	}
}

func TestAllAnimations(t *testing.T) {
	t.Parallel()

	anims := AllAnimations()

	if len(anims) != 11 {
		t.Errorf("AllAnimations() returned %d animations, want 11", len(anims))
	}

	for i := 1; i < len(anims); i++ {
		if string(anims[i-1]) >= string(anims[i]) {
			t.Errorf("AllAnimations() not sorted alphabetically at index %d: %q >= %q", i, anims[i-1], anims[i])
		}
	}

	if slices.Contains(anims, AnimNone) {
		t.Error("AllAnimations() should not include AnimNone")
	}
}

func TestDefaultAnimationConsistency(t *testing.T) {
	t.Parallel()

	t.Run("all mapped icons are valid names", func(t *testing.T) {
		t.Parallel()

		for name := range defaultAnimations {
			if !NameIsValid(name) {
				t.Errorf("defaultAnimations contains invalid icon name %q", name)
			}
		}
	})

	t.Run("all mapped animations are valid", func(t *testing.T) {
		t.Parallel()

		for _, anim := range defaultAnimations {
			if !anim.IsValid() {
				t.Errorf("defaultAnimations contains invalid animation %q", anim)
			}
		}
	})
}

// TestCompleteAnimationCoverage verifies that every icon in iconPathData has
// either a direct entry in defaultAnimations or resolves through iconAliases
// to a name that does. No icon should silently fall back to the generic
// AnimPulse unless it is explicitly mapped to pulse.
func TestCompleteAnimationCoverage(t *testing.T) {
	t.Parallel()

	unmapped := make([]Name, 0)

	for name := range iconPathData {
		if _, ok := defaultAnimations[name]; ok {
			continue
		}

		if alias, ok := iconAliases[name]; ok {
			if _, ok := defaultAnimations[alias]; ok {
				continue
			}
		}

		unmapped = append(unmapped, name)
	}

	if len(unmapped) > 0 {
		t.Errorf("icons without animation mapping (would fall back to AnimPulse): %v", unmapped)
	}
}

// TestBlinkIconsHaveMultiplePaths verifies that all icons mapped to AnimBlink
// have 2+ SVG path elements (required for the per-path nth-child CSS targeting).
func TestBlinkIconsHaveMultiplePaths(t *testing.T) {
	t.Parallel()

	for name, anim := range defaultAnimations {
		if anim != AnimBlink {
			continue
		}

		actualPaths := len(iconPaths(name))
		if actualPaths < 2 {
			t.Errorf(
				"icon %q mapped to AnimBlink (requires 2+ paths) but has only %d",
				name, actualPaths,
			)
		}
	}
}
