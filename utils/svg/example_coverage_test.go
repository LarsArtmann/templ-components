package svg_test

import (
	"github.com/larsartmann/templ-components/utils/svg"
)

func ExampleFillIcon() {
	_ = svg.FillIcon("w-5 h-5", "M12 2a10 10 0 100 20 10 10 0 000-20z", false)
	// Output:
}

func ExampleSpinnerSVG() {
	_ = svg.SpinnerSVG()
	// Output:
}
