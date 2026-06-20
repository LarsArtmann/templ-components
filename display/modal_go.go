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
var modalSizeLookup = map[string]string{
	string(ModalSizeSM):   "max-w-sm",
	string(ModalSizeMD):   "max-w-md",
	string(ModalSizeLG):   "max-w-lg",
	string(ModalSizeXL):   "max-w-xl",
	string(ModalSizeFull): "max-w-4xl",
}

func modalSizeClass(size ModalSize) string {
	return utils.Lookup(modalSizeLookup, string(size), modalSizeLookup[string(ModalSizeMD)])
}
