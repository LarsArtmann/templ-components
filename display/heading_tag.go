package display

// HeadingTagType selects the heading element (h1–h6) a component renders for
// its title. The zero value and any unknown value fall back to h3 at render
// time. As a defined string type, untyped literals ("h2") still assign, so
// existing consumer call sites keep compiling.
type HeadingTagType string

const (
	HeadingTagH1 HeadingTagType = "h1"
	HeadingTagH2 HeadingTagType = "h2"
	HeadingTagH3 HeadingTagType = "h3"
	HeadingTagH4 HeadingTagType = "h4"
	HeadingTagH5 HeadingTagType = "h5"
	HeadingTagH6 HeadingTagType = "h6"

	// HeadingTagDefault is the heading level used when TitleTag is empty or
	// unknown — h3 keeps the common card/section/empty-state outline sane
	// under a page-level h1.
	HeadingTagDefault = HeadingTagH3
)

// HeadingTagTypeIsValid reports whether t is one of h1–h6.
func HeadingTagTypeIsValid(t HeadingTagType) bool {
	switch t {
	case HeadingTagH1, HeadingTagH2, HeadingTagH3, HeadingTagH4, HeadingTagH5, HeadingTagH6:
		return true
	default:
		return false
	}
}
