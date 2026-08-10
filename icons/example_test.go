package icons_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/icons"
)

func ExampleIcon() {
	var buf bytes.Buffer

	_ = icons.Icon(icons.Check, "h-5 w-5 text-green-500").Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleAnimatedIcon() {
	var buf bytes.Buffer

	// Renders a Heart icon with its default hover animation (pulse).
	// Requires .tc-anim-* CSS classes from templates/custom.css.
	_ = icons.AnimatedIcon(icons.Heart, "h-6 w-6 text-red-500").Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleAnimatedIconWithAnimation() {
	var buf bytes.Buffer

	// Apply any animation to any icon.
	_ = icons.AnimatedIconWithAnimation(icons.Bell, icons.AnimWiggle, "h-6 w-6").Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleAnimation_IsValid() {
	fmt.Println(icons.AnimPulse.IsValid())
	fmt.Println(icons.AnimNone.IsValid())

	// Output: true
	// false
}

func ExampleDefaultAnimation() {
	fmt.Println(icons.DefaultAnimation(icons.Heart))
	fmt.Println(icons.DefaultAnimation(icons.Settings))
	fmt.Println(icons.DefaultAnimation(icons.Spinner))

	// Output: pulse
	// spin
	//
}
