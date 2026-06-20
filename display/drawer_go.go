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

const (
	DrawerSizeSM DrawerSize = "sm"
	DrawerSizeMD DrawerSize = "md"
	DrawerSizeLG DrawerSize = "lg"
	DrawerSizeXL DrawerSize = "xl"
	DrawerFull   DrawerSize = "full"
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
var drawerSizeLookup = map[string]string{
	string(DrawerSizeSM): maxWSM,
	string(DrawerSizeMD): maxWMD,
	string(DrawerSizeLG): maxWLG,
	string(DrawerSizeXL): maxWXL,
	string(DrawerFull):   "max-w-2xl",
}

func drawerSizeClass(size DrawerSize) string {
	return utils.Lookup(drawerSizeLookup, string(size), drawerSizeLookup[string(DrawerSizeMD)])
}
