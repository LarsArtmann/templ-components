// Modal component: size enum, default props, and class lookup.
package display

import (
	"github.com/larsartmann/templ-components/utils"
)

// ModalSize defines the width of a modal dialog
type ModalSize string

// Modal size constants. ModalSize2XL is the largest size (max-w-4xl);
// ModalSizeFull is a deprecated alias kept for backward compatibility.
const (
	ModalSizeSM   ModalSize = "sm"
	ModalSizeMD   ModalSize = "md"
	ModalSizeLG   ModalSize = "lg"
	ModalSizeXL   ModalSize = "xl"
	ModalSize2XL  ModalSize = "full" // largest available width (max-w-4xl)
	ModalSizeFull ModalSize = "full" // Deprecated: use ModalSize2XL; "full" is a misnomer (capped at max-w-4xl)
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

// ModalSizeIsValid reports whether s is one of the defined ModalSize constants.
func ModalSizeIsValid(s ModalSize) bool {
	_, ok := modalSizeLookup[s]
	return ok
}
