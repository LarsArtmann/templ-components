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
		{"split", AnimSplit, true},
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
		{"Heart defaults to pulse", Heart, AnimPulse},
		{"Star defaults to beat", Star, AnimBeat},
		{"Bell defaults to wiggle", Bell, AnimWiggle},
		{"Settings defaults to spin", Settings, AnimSpin},
		{"Eye defaults to blink", Eye, AnimBlink},
		{"Trash defaults to wiggle", Trash, AnimWiggle},
		{"Home defaults to jump", Home, AnimJump},
		{"Search defaults to bounce", Search, AnimBounce},
		{"Spinner defaults to none", Spinner, AnimNone},
		{"unmapped icon defaults to pulse", Document, AnimPulse},
		{"Wrench defaults to spin", Wrench, AnimSpin},
		{"ChevronDown defaults to nod", ChevronDown, AnimNod},
		{"ExternalLink defaults to shake", ExternalLink, AnimShake},
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

	if len(anims) != 10 {
		t.Errorf("AllAnimations() returned %d animations, want 10", len(anims))
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
