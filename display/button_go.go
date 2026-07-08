// Button component: variant/size enums, default props, and class lookups.
package display

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ButtonType defines the visual style of a button
type ButtonType string

// Button variant constants
const (
	ButtonPrimary   ButtonType = "primary"
	ButtonSecondary ButtonType = "secondary"
	ButtonDanger    ButtonType = "danger"
	ButtonGhost     ButtonType = "ghost"
	ButtonLink      ButtonType = "link"
)

// ButtonSize defines the size of a button
type ButtonSize string

// Button size constants
const (
	ButtonSizeSM ButtonSize = "sm"
	ButtonSizeMD ButtonSize = "md"
	ButtonSizeLG ButtonSize = "lg"
)

// ButtonHTMLType is the HTML type attribute for a <button> element.
type ButtonHTMLType string

// Button HTML type constants
const (
	ButtonHTMLButton ButtonHTMLType = "button"
	ButtonHTMLSubmit ButtonHTMLType = "submit"
	ButtonHTMLReset  ButtonHTMLType = "reset"
)

//nolint:gochecknoglobals // Package-level lookup table for button HTML types
var buttonHTMLTypeLookup = map[ButtonHTMLType]string{
	ButtonHTMLButton: string(ButtonHTMLButton),
	ButtonHTMLSubmit: string(ButtonHTMLSubmit),
	ButtonHTMLReset:  string(ButtonHTMLReset),
}

// buttonHTMLType returns the validated HTML type attribute. Unknown or empty
// values fall back to "button" — matching the map+fallback convention used by
// all other enums in the library and the HTML spec default.
func buttonHTMLType(t ButtonHTMLType) string {
	return utils.Lookup(buttonHTMLTypeLookup, t, string(ButtonHTMLButton))
}

// ButtonProps configures a button or link styled as a button
type ButtonProps struct {
	utils.BaseProps
	Text     string
	Type     ButtonHTMLType // button, submit, reset (default: button; ignored when Href is set)
	Href     string         // if set, renders as <a> instead of <button>
	Variant  ButtonType     // default: Primary
	Size     ButtonSize     // default: MD
	Disabled bool
	Icon     templ.Component
	External bool // adds target="_blank" and rel="noopener noreferrer" for links
}

// DefaultButtonProps returns sensible defaults
func DefaultButtonProps() ButtonProps {
	return ButtonProps{ //nolint:exhaustruct // intentionally minimal defaults
		Type:    ButtonHTMLButton,
		Variant: ButtonPrimary,
		Size:    ButtonSizeMD,
	}
}

// buttonVariantClasses returns the color classes for each button variant
//
//nolint:gochecknoglobals // Package-level lookup table for button variants
var buttonVariantLookup = map[ButtonType]string{
	ButtonPrimary:   "bg-blue-600 text-white hover:bg-blue-500 dark:bg-blue-500 dark:hover:bg-blue-400 focus-visible:outline-blue-600 dark:focus-visible:outline-blue-500",
	ButtonSecondary: "bg-white text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-gray-800 dark:text-white dark:ring-gray-600 dark:hover:bg-gray-700",
	ButtonDanger:    "bg-red-600 text-white hover:bg-red-500 dark:bg-red-500 dark:hover:bg-red-400 focus-visible:outline-red-600 dark:focus-visible:outline-red-500",
	ButtonGhost:     "bg-transparent text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700",
	ButtonLink:      "bg-transparent text-blue-600 hover:text-blue-500 dark:text-blue-400 dark:hover:text-blue-300 underline-offset-2 hover:underline",
}

func buttonVariantClasses(v ButtonType) string {
	return utils.Lookup(buttonVariantLookup, v, buttonVariantLookup[ButtonPrimary])
}

// ButtonTypeIsValid reports whether v is one of the defined ButtonType constants.
func ButtonTypeIsValid(v ButtonType) bool {
	_, ok := buttonVariantLookup[v]
	return ok
}

//nolint:gochecknoglobals // Package-level lookup table for button sizes
var buttonSizeLookup = map[ButtonSize]string{
	ButtonSizeSM: "px-2.5 py-1.5 text-xs",
	ButtonSizeMD: "px-3 py-2 text-sm",
	ButtonSizeLG: "px-4 py-2.5 text-base",
}

func buttonSizeClasses(s ButtonSize) string {
	return utils.Lookup(buttonSizeLookup, s, buttonSizeLookup[ButtonSizeMD])
}

// buttonBaseClass returns the shared base classes for all buttons
const buttonBaseClass = "inline-flex items-center justify-center rounded-md font-semibold transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
