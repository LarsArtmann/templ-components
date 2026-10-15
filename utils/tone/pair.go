package tone

import "strings"

// Pair is one utility in light and dark form. Both halves are COMPLETE
// class literals — Dark carries its own dark: prefix — because Tailwind's
// scanner only emits CSS for tokens it can see in full (the package-level
// scanner-safety contract). Pair makes the halves structural so tests can
// prove every dark half exists.
type Pair struct {
	// Light is the light-mode class literal, without any variant prefix
	// (e.g. "text-green-600").
	Light string
	// Dark is the dark-mode class literal WITH its dark: prefix
	// (e.g. "dark:text-green-400").
	Dark string
}

// Classes returns the renderable class string: light and dark halves joined
// by a single space. Either half may be empty (single-mode utilities); the
// join never leaves leading/trailing whitespace.
func (p Pair) Classes() string {
	return strings.TrimSpace(p.Light + " " + p.Dark)
}

// Complete reports whether both halves are non-empty. Most library pairs
// are complete (every color utility carries a dark: variant — guarded by
// utils.TestDarkModeCompliance); incomplete pairs are deliberate
// single-mode exceptions and must say so at their construction site.
func (p Pair) Complete() bool {
	return p.Light != "" && p.Dark != ""
}
