package display

import (
	"github.com/larsartmann/templ-components/utils"
)

// DrawerSide defines which side the drawer slides in from
type DrawerSide string

const (
	DrawerLeft  DrawerSide = "left"
	DrawerRight DrawerSide = "right"
)

// DrawerSize defines the width of the drawer panel
type DrawerSize string

// Drawer size constants. DrawerSize2XL is the largest size (max-w-2xl);
// DrawerFull is a deprecated alias kept for backward compatibility.
const (
	DrawerSizeSM  DrawerSize = "sm"
	DrawerSizeMD  DrawerSize = "md"
	DrawerSizeLG  DrawerSize = "lg"
	DrawerSizeXL  DrawerSize = "xl"
	DrawerSize2XL DrawerSize = "full" // largest available width (max-w-2xl)
	DrawerFull    DrawerSize = "full" // Deprecated: use DrawerSize2XL; "full" is a misnomer (capped at max-w-2xl)
)

// DrawerProps configures a drawer (side panel) component
type DrawerProps struct {
	utils.BaseProps
	Title string
	Open  bool
	Side  DrawerSide
	Size  DrawerSize
}

// DefaultDrawerProps returns sensible defaults
func DefaultDrawerProps() DrawerProps {
	return DrawerProps{ //nolint:exhaustruct // intentionally minimal defaults
		Side: DrawerRight,
		Size: DrawerSizeMD,
	}
}

const (
	maxWSM = "max-w-sm"
	maxWMD = "max-w-md"
	maxWLG = "max-w-lg"
	maxWXL = "max-w-xl"
)

//nolint:gochecknoglobals // Package-level lookup table for drawer sizes
var drawerSizeLookup = map[DrawerSize]string{
	DrawerSizeSM: maxWSM,
	DrawerSizeMD: maxWMD,
	DrawerSizeLG: maxWLG,
	DrawerSizeXL: maxWXL,
	DrawerFull:   "max-w-2xl",
}

func drawerSizeClass(size DrawerSize) string {
	return utils.Lookup(drawerSizeLookup, size, drawerSizeLookup[DrawerSizeMD])
}
