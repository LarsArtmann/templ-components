package icons_test

import (
	"github.com/larsartmann/templ-components/icons"
)

func ExampleIconWithStrokeWidth() {
	_ = icons.IconWithStrokeWidth(icons.Search, "w-5 h-5", 2)
	// Output:
}

func ExampleIconRTL() {
	_ = icons.IconRTL(icons.Search, "w-5 h-5")
	// Output:
}

func ExampleAnimatedIconRTL() {
	_ = icons.AnimatedIconRTL(icons.Search, "w-5 h-5")
	// Output:
}

func ExampleAnimatedIconWithAnimationRTL() {
	_ = icons.AnimatedIconWithAnimationRTL(icons.Search, icons.AnimPulse, "w-5 h-5")
	// Output:
}

func ExampleRender() {
	_ = icons.Render(icons.CustomIcon{
		ViewBox: "0 0 24 24",
		Paths:   []string{"M12 2a10 10 0 100 20 10 10 0 000-20z"},
	}, "w-5 h-5")
	// Output:
}
