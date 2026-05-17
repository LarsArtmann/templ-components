// Package display provides modal dialog components and related functionality.
package display

import (
	"errors"
	"strconv"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ModalSize defines the width of a modal dialog
type ModalSize string

// Modal size constants
const (
	ModalSizeSM   ModalSize = "sm"
	ModalSizeMD   ModalSize = "md"
	ModalSizeLG   ModalSize = "lg"
	ModalSizeXL   ModalSize = "xl"
	ModalSizeFull ModalSize = "full"
)

// ModalProps configures a modal dialog
type ModalProps struct {
	utils.BaseProps
	Title string
	Open  bool
	Size  ModalSize
}

// DefaultModalProps returns sensible defaults
func DefaultModalProps() ModalProps {
	return ModalProps{ //nolint:exhaustruct // intentionally minimal defaults
		Size: ModalSizeMD,
	}
}

func modalSizeClass(size ModalSize) string {
	switch size {
	case ModalSizeSM:
		return "max-w-sm"
	case ModalSizeMD:
		return "max-w-md"
	case ModalSizeLG:
		return "max-w-lg"
	case ModalSizeXL:
		return "max-w-xl"
	case ModalSizeFull:
		return "max-w-4xl"
	default:
		return "max-w-md"
	}
}

func modalCloseHandler(id string) templ.ComponentScript {
	name := "tcCloseModal_" + id
	fn := "function tcCloseModal_" + id + "(id){tcCloseModal(id)}"
	escapedID := strconv.Quote(id)
	call := "tcCloseModal(" + escapedID + ")"
	return templ.ComponentScript{
		Name:       name,
		Function:   fn,
		Call:       call,
		CallInline: "",
	}
}

func validateModalID(id string) error {
	if id == "" {
		return errors.New(
			"Modal requires a non-empty ID for ARIA attributes and JavaScript functionality",
		)
	}
	return nil
}
