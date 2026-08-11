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
		{"Eye with blink", Eye, AnimBlink, "tc-anim-blink"},
		{"Star with beat", Star, AnimBeat, "tc-anim-beat"},
		{"Search with bounce", Search, AnimBounce, "tc-anim-bounce"},
		{"Home with jump", Home, AnimJump, "tc-anim-jump"},
		{"ChevronDown with nod", ChevronDown, AnimNod, "tc-anim-nod"},
		{"ArrowRight with shake", ArrowRight, AnimShake, "tc-anim-shake"},
		{"Beaker with wobble", Beaker, AnimWobble, "tc-anim-wobble"},
		{"Heart with pulse", Heart, AnimPulse, "tc-anim-pulse"},
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

func TestAnimatedIconWithDrawRendersPathLength(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, AnimatedIconWithAnimation(Bolt, AnimDraw, "h-6 w-6"))

	utils.AssertContains(t, output, "tc-anim-draw")
	utils.AssertContains(t, output, `pathLength="1"`)
	utils.AssertContains(t, output, "<svg")
}

func TestAnimatedIconBoltDefaultsToDraw(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, AnimatedIcon(Bolt, "h-6 w-6"))

	utils.AssertContains(t, output, "tc-anim-draw")
	utils.AssertContains(t, output, `pathLength="1"`)
}

func TestAnimatedIconRefreshDefaultsToSpin(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, AnimatedIcon(Refresh, "h-6 w-6"))

	utils.AssertContains(t, output, "tc-anim-spin")
	utils.AssertContains(t, output, "<svg")
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
