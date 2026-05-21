// Package display provides modal dialog components and related functionality.
package display

import (
	"fmt"
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

//nolint:gochecknoglobals // Package-level lookup table for modal sizes
var modalSizeLookup = map[string]string{
	string(ModalSizeSM):   "max-w-sm",
	string(ModalSizeMD):   "max-w-md",
	string(ModalSizeLG):   "max-w-lg",
	string(ModalSizeXL):   "max-w-xl",
	string(ModalSizeFull): "max-w-4xl",
}

func modalSizeClass(size ModalSize) string {
	if v, ok := modalSizeLookup[string(size)]; ok {
		return v
	}
	return modalSizeLookup[string(ModalSizeMD)]
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
		return fmt.Errorf("modal: id=%q cannot be empty", id)
	}
	return nil
}

func modalSafeID(id string) string {
	return strconv.Quote(id)
}
