// Package display provides button components.
package display

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ButtonType defines the visual style of a button
type ButtonType string

const (
	ButtonPrimary   ButtonType = "primary"
	ButtonSecondary ButtonType = "secondary"
	ButtonDanger    ButtonType = "danger"
	ButtonGhost     ButtonType = "ghost"
	ButtonLink      ButtonType = "link"
)

// ButtonSize defines the size of a button
type ButtonSize string

const (
	ButtonSizeSM ButtonSize = "sm"
	ButtonSizeMD ButtonSize = "md"
	ButtonSizeLG ButtonSize = "lg"
)

// ButtonProps configures a button or link styled as a button
type ButtonProps struct {
	utils.BaseProps
	Text     string
	Type     string     // button, submit, reset (ignored when Href is set)
	Href     string     // if set, renders as <a> instead of <button>
	Variant  ButtonType // default: Primary
	Size     ButtonSize // default: MD
	Disabled bool
	Icon     templ.Component
	External bool // adds target="_blank" and rel="noopener noreferrer" for links
}

// DefaultButtonProps returns sensible defaults
func DefaultButtonProps() ButtonProps {
	return ButtonProps{
		Type:    "button",
		Variant: ButtonPrimary,
		Size:    ButtonSizeMD,
	}
}

// buttonVariantClasses returns the color classes for each button variant
func buttonVariantClasses(v ButtonType) string {
	switch v {
	case ButtonSecondary:
		return "bg-white text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 dark:bg-slate-800 dark:text-white dark:ring-slate-600 dark:hover:bg-slate-700"
	case ButtonDanger:
		return "bg-red-600 text-white hover:bg-red-500 focus-visible:outline-red-600"
	case ButtonGhost:
		return "bg-transparent text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-slate-700"
	case ButtonLink:
		return "bg-transparent text-blue-600 hover:text-blue-500 underline-offset-2 hover:underline"
	default:
		return "bg-blue-600 text-white hover:bg-blue-500 focus-visible:outline-blue-600"
	}
}

// buttonSizeClasses returns the size classes for each button size
func buttonSizeClasses(s ButtonSize) string {
	switch s {
	case ButtonSizeSM:
		return "px-2.5 py-1.5 text-xs"
	case ButtonSizeLG:
		return "px-4 py-2.5 text-base"
	default:
		return "px-3 py-2 text-sm"
	}
}

// buttonBaseClass returns the shared base classes for all buttons
const buttonBaseClass = "inline-flex items-center justify-center rounded-md font-semibold transition-colors focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
