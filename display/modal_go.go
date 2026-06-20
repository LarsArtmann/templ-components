// Package display provides modal dialog components and related functionality.
package display

import (
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
var modalSizeLookup = map[ModalSize]string{
	ModalSizeSM:   "max-w-sm",
	ModalSizeMD:   "max-w-md",
	ModalSizeLG:   "max-w-lg",
	ModalSizeXL:   "max-w-xl",
	ModalSizeFull: "max-w-4xl",
}

func modalSizeClass(size ModalSize) string {
	return utils.Lookup(modalSizeLookup, size, modalSizeLookup[ModalSizeMD])
}
