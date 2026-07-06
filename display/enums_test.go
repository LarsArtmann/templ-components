package display

import "testing"

func TestIsValidEnums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func() bool
		want bool
	}{
		// BadgeType
		{"BadgeType primary", func() bool { return BadgeTypeIsValid(BadgePrimary) }, true},
		{"BadgeType invalid", func() bool { return BadgeTypeIsValid(BadgeType("bogus")) }, false},
		// BadgeSize
		{"BadgeSize SM", func() bool { return BadgeSizeIsValid(BadgeSizeSM) }, true},
		{"BadgeSize invalid", func() bool { return BadgeSizeIsValid(BadgeSize("bogus")) }, false},
		// CardPadding
		{"CardPadding MD", func() bool { return CardPaddingIsValid(CardPaddingMD) }, true},
		{"CardPadding None", func() bool { return CardPaddingIsValid(CardPaddingNone) }, true},
		{"CardPadding invalid", func() bool { return CardPaddingIsValid(CardPadding("bogus")) }, false},
		// GridCols
		{"GridCols 1", func() bool { return GridColsIsValid(GridCols1) }, true},
		{"GridCols invalid", func() bool { return GridColsIsValid(GridCols("bogus")) }, false},
		// GridGap
		{"GridGap SM", func() bool { return GridGapIsValid(GridGapSM) }, true},
		{"GridGap MD", func() bool { return GridGapIsValid(GridGapMD) }, true},
		{"GridGap invalid", func() bool { return GridGapIsValid(GridGap("bogus")) }, false},
		// TrendDirection
		{"TrendDirection Up", func() bool { return TrendDirectionIsValid(TrendUp) }, true},
		{"TrendDirection invalid", func() bool { return TrendDirectionIsValid(TrendDirection("bogus")) }, false},
		// AvatarSize
		{"AvatarSize SM", func() bool { return AvatarSizeIsValid(AvatarSizeSM) }, true},
		{"AvatarSize invalid", func() bool { return AvatarSizeIsValid(AvatarSize("bogus")) }, false},
		// TooltipPosition
		{"TooltipPosition Top", func() bool { return TooltipPositionIsValid(TooltipPositionTop) }, true},
		{"TooltipPosition invalid", func() bool { return TooltipPositionIsValid(TooltipPosition("bogus")) }, false},
		// AvatarShape
		{"AvatarShape Circle", func() bool { return AvatarShapeIsValid(AvatarShapeCircle) }, true},
		{"AvatarShape invalid", func() bool { return AvatarShapeIsValid(AvatarShape("bogus")) }, false},
		// AvatarStatus
		{"AvatarStatus Online", func() bool { return AvatarStatusIsValid(AvatarStatusOnline) }, true},
		{"AvatarStatus None", func() bool { return AvatarStatusIsValid(AvatarStatusNone) }, true},
		{"AvatarStatus invalid", func() bool { return AvatarStatusIsValid(AvatarStatus("bogus")) }, false},
		// DropdownItemKind
		{"DropdownItemKind Link", func() bool { return DropdownItemKindIsValid(DropdownItemLink) }, true},
		{"DropdownItemKind invalid", func() bool { return DropdownItemKindIsValid(DropdownItemKind("bogus")) }, false},
		// DropdownPosition
		{"DropdownPosition Left", func() bool { return DropdownPositionIsValid(DropdownPositionLeft) }, true},
		{"DropdownPosition invalid", func() bool { return DropdownPositionIsValid(DropdownPosition("bogus")) }, false},
		// TabsVariant
		{"TabsVariant Default", func() bool { return TabsVariantIsValid(TabsDefault) }, true},
		{"TabsVariant invalid", func() bool { return TabsVariantIsValid(TabsVariant("bogus")) }, false},
		// OverlayKind
		{"OverlayKind Modal", func() bool { return OverlayKindIsValid(OverlayModal) }, true},
		{"OverlayKind invalid", func() bool { return OverlayKindIsValid(OverlayKind("bogus")) }, false},
		// ButtonSize
		{"ButtonSize SM", func() bool { return ButtonSizeIsValid(ButtonSizeSM) }, true},
		{"ButtonSize invalid", func() bool { return ButtonSizeIsValid(ButtonSize("bogus")) }, false},
		// ButtonHTMLType
		{"ButtonHTMLType Button", func() bool { return ButtonHTMLTypeIsValid(ButtonHTMLButton) }, true},
		{"ButtonHTMLType invalid", func() bool { return ButtonHTMLTypeIsValid(ButtonHTMLType("bogus")) }, false},
		// ModalSize
		{"ModalSize MD", func() bool { return ModalSizeIsValid(ModalSizeMD) }, true},
		{"ModalSize 2XL", func() bool { return ModalSizeIsValid(ModalSize2XL) }, true},
		{"ModalSize invalid", func() bool { return ModalSizeIsValid(ModalSize("bogus")) }, false},
		// DrawerSize
		{"DrawerSize MD", func() bool { return DrawerSizeIsValid(DrawerSizeMD) }, true},
		{"DrawerSize 2XL", func() bool { return DrawerSizeIsValid(DrawerSize2XL) }, true},
		{"DrawerSize invalid", func() bool { return DrawerSizeIsValid(DrawerSize("bogus")) }, false},
		// DrawerSide
		{"DrawerSide Left", func() bool { return DrawerSideIsValid(DrawerLeft) }, true},
		{"DrawerSide invalid", func() bool { return DrawerSideIsValid(DrawerSide("bogus")) }, false},
		// ButtonType
		{"ButtonType Primary", func() bool { return ButtonTypeIsValid(ButtonPrimary) }, true},
		{"ButtonType invalid", func() bool { return ButtonTypeIsValid(ButtonType("bogus")) }, false},
		// SortDirection
		{"SortDirection None", func() bool { return SortDirectionIsValid(SortNone) }, true},
		{"SortDirection Asc", func() bool { return SortDirectionIsValid(SortAsc) }, true},
		{"SortDirection Desc", func() bool { return SortDirectionIsValid(SortDesc) }, true},
		{"SortDirection invalid", func() bool { return SortDirectionIsValid(SortDirection("bogus")) }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.fn(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
