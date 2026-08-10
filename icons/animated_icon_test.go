package icons

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestAnimatedIconRendersWrapper(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, AnimatedIcon(Heart, "h-6 w-6 text-red-500"))

	utils.AssertContains(t, output, `<span`)
	utils.AssertContains(t, output, "tc-anim")
	utils.AssertContains(t, output, "tc-anim-pulse")
	utils.AssertContains(t, output, "inline-flex")
	utils.AssertContains(t, output, "<svg")
	utils.AssertContains(t, output, "h-6 w-6 text-red-500")
	utils.AssertContains(t, output, "</span>")
}

func TestAnimatedIconWithAnimation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		icon      Name
		anim      Animation
		wantClass string
	}{
		{"Bell with wiggle", Bell, AnimWiggle, "tc-anim-wiggle"},
		{"Settings with spin", Settings, AnimSpin, "tc-anim-spin"},
		{"Trash with wiggle", Trash, AnimWiggle, "tc-anim-wiggle"},
		{"Eye with blink", Eye, AnimBlink, "tc-anim-blink"},
		{"Star with beat", Star, AnimBeat, "tc-anim-beat"},
		{"Search with bounce", Search, AnimBounce, "tc-anim-bounce"},
		{"Home with jump", Home, AnimJump, "tc-anim-jump"},
		{"ChevronDown with nod", ChevronDown, AnimNod, "tc-anim-nod"},
		{"ArrowRight with shake", ArrowRight, AnimShake, "tc-anim-shake"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, AnimatedIconWithAnimation(tt.icon, tt.anim, "h-5 w-5"))

			utils.AssertContains(t, output, "tc-anim")
			utils.AssertContains(t, output, tt.wantClass)
			utils.AssertContains(t, output, "<svg")
			utils.AssertContains(t, output, "h-5 w-5")
		})
	}
}

func TestAnimatedIconWithNoneRendersPlainIcon(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, AnimatedIconWithAnimation(Heart, AnimNone, "h-6 w-6"))

	utils.AssertNotContains(t, output, "tc-anim")
	utils.AssertNotContains(t, output, "<span")
	utils.AssertContains(t, output, "<svg")
	utils.AssertContains(t, output, "h-6 w-6")
}

func TestAnimatedIconSpinnerDefaultsToNone(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, AnimatedIcon(Spinner, "h-6 w-6"))

	utils.AssertNotContains(t, output, "tc-anim")
	utils.AssertNotContains(t, output, "<span")
	utils.AssertContains(t, output, "animate-spin")
}

func TestAllPathIconsCanAnimate(t *testing.T) {
	t.Parallel()

	for _, name := range allIconNames() {
		if name == Spinner {
			continue
		}

		t.Run(string(name), func(t *testing.T) {
			t.Parallel()

			output := utils.Render(t, AnimatedIcon(name, "h-5 w-5"))

			utils.AssertContains(t, output, "tc-anim")
			utils.AssertContains(t, output, "<svg")

			anim := DefaultAnimation(name)
			if anim != AnimNone {
				utils.AssertContains(t, output, "tc-anim-"+string(anim))
			}
		})
	}
}

func TestPerPathAnimationsHaveCorrectPathCount(t *testing.T) {
	t.Parallel()

	// AnimBlink requires 2+ paths; AnimSplit requires 2+ paths.
	// Verify that all icons defaulted to these animations have enough paths.
	perPathAnimations := map[Animation]int{
		AnimBlink: 2,
		AnimSplit: 2,
	}

	for name, anim := range defaultAnimations {
		requiredPaths, isPerPath := perPathAnimations[anim]
		if !isPerPath {
			continue
		}

		actualPaths := len(iconPaths(name))
		if actualPaths < requiredPaths {
			t.Errorf(
				"icon %q defaulted to %s (requires %d paths) but has only %d",
				name, anim, requiredPaths, actualPaths,
			)
		}
	}
}

func TestAnimatedIconMultiPathStructure(t *testing.T) {
	t.Parallel()

	t.Run("Eye has 2 paths for blink animation", func(t *testing.T) {
		t.Parallel()

		paths := iconPaths(Eye)
		if len(paths) != 2 {
			t.Errorf("Eye should have 2 paths for blink animation, got %d", len(paths))
		}
	})

	t.Run("Trash has 1 combined path", func(t *testing.T) {
		t.Parallel()

		paths := iconPaths(Trash)
		if len(paths) != 1 {
			t.Errorf("Trash should have 1 combined path, got %d", len(paths))
		}
	})
}

func TestAnimatedIconProducesValidHTML(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, AnimatedIcon(Heart, "h-6 w-6"))

	if !strings.HasPrefix(output, "<span") {
		t.Errorf("AnimatedIcon output should start with <span, got: %s", output[:min(len(output), 50)])
	}

	if !strings.HasSuffix(strings.TrimSpace(output), "</span>") {
		t.Errorf("AnimatedIcon output should end with </span>")
	}
}
