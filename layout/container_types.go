package layout

import "github.com/larsartmann/templ-components/utils"

// ContainerWidth is a typed enum for the maximum width of a Container. Unknown
// values fall back to ContainerWidthLG (graceful degradation — never panic).
type ContainerWidth string

const (
	// ContainerWidthSM narrows the container to max-w-3xl (48rem) — for
	// sign-in forms, focused single-purpose pages, and short articles.
	ContainerWidthSM ContainerWidth = "sm"
	// ContainerWidthMD is max-w-5xl (64rem) — for two-column layouts and
	// medium-density pages.
	ContainerWidthMD ContainerWidth = "md"
	// ContainerWidthLG is max-w-7xl (80rem) — the application default for
	// dashboards, list pages, and admin panels.
	ContainerWidthLG ContainerWidth = "lg"
	// ContainerWidthXL is max-w-[90rem] — for wide dashboards and data-heavy
	// layouts that need more horizontal real estate than LG.
	ContainerWidthXL ContainerWidth = "xl"
	// ContainerWidthFull is max-w-full — full-bleed edge-to-edge content.
	// Pair with explicit px-* padding on the consumer side.
	ContainerWidthFull ContainerWidth = "full"
	// ContainerWidthProse is max-w-prose (65ch) — for long-form articles and
	// documentation where line length must stay readable.
	ContainerWidthProse ContainerWidth = "prose"
	// ContainerWidthDefault is the canonical default (LG).
	ContainerWidthDefault ContainerWidth = ContainerWidthLG
)

// ContainerWidthIsValid reports whether w is a recognized ContainerWidth.
func ContainerWidthIsValid(w ContainerWidth) bool {
	_, ok := containerWidthLookup[w]

	return ok
}

//nolint:gochecknoglobals // Package-level lookup table for container widths
var containerWidthLookup = map[ContainerWidth]string{
	ContainerWidthSM:    "max-w-3xl",
	ContainerWidthMD:    "max-w-5xl",
	ContainerWidthLG:    "max-w-7xl",
	ContainerWidthXL:    "max-w-[90rem]",
	ContainerWidthFull:  "max-w-full",
	ContainerWidthProse: "max-w-prose",
}

// containerWidthClass returns the Tailwind max-width class for a width, with
// fallback to ContainerWidthDefault for unknown values (map+fallback pattern).
func containerWidthClass(w ContainerWidth) string {
	return utils.Lookup(containerWidthLookup, w, containerWidthLookup[ContainerWidthDefault])
}

// ContainerProps configures a content Container — a centered max-width wrapper
// with responsive horizontal padding. This is the single source of truth for
// page content width; use it to replace the repeated
// `max-w-6xl mx-auto px-4 sm:px-6 lg:px-8` snippet.
//
// All new page content should be wrapped in a Container. AppShell already
// wraps its main content slot in a Container by default — set AppShellProps
// .Container=false to opt out (e.g. when the page itself manages width).
type ContainerProps struct {
	utils.BaseProps

	// Width controls the max-width. Default is ContainerWidthLG.
	// Use ContainerWidthProse for long-form articles.
	Width ContainerWidth
	// Pad enables the responsive horizontal padding
	// (px-4 sm:px-6 lg:px-8). Default is true. Set to false when the
	// consumer manages padding (e.g. inside a Card or Grid).
	Pad bool
}

// DefaultContainerProps returns sensible defaults: LG width with responsive
// padding enabled.
//
//nolint:exhaustruct // constructor intentionally sets only non-zero defaults
func DefaultContainerProps() ContainerProps {
	return ContainerProps{
		Width: ContainerWidthDefault,
		Pad:   true,
	}
}
