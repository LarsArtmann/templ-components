// Package display provides modal dialog components and related functionality.
package display

import (
	"strconv"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ModalProps configures a modal dialog
type ModalProps struct {
	utils.BaseProps
	Title string
	Open  bool
	Size  string // "sm", "md", "lg", "xl", "full"
}

// DefaultModalProps returns sensible defaults
func DefaultModalProps() ModalProps {
	return ModalProps{ //nolint:exhaustruct // intentionally minimal defaults
		Size: "md",
	}
}

func modalSizeClass(size string) string {
	switch size {
	case "sm":
		return "max-w-sm"
	case "lg":
		return "max-w-lg"
	case "xl":
		return "max-w-xl"
	case "full":
		return "max-w-4xl"
	default:
		return "max-w-md"
	}
}

func modalCloseHandler(id string) templ.ComponentScript {
	name := "tcCloseModal_" + id
	fn := "function tcCloseModal_" + id + "(id){tcCloseModal(id)}"
	// Use JSON string escaping to safely embed id in JavaScript
	escapedID := strconv.Quote(id)
	call := "tcCloseModal(" + escapedID + ")"
	return templ.ComponentScript{
		Name:       name,
		Function:   fn,
		Call:       call,
		CallInline: "",
	}
}
