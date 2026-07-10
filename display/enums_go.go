package display

// BadgeTypeIsValid reports whether v is one of the defined BadgeType constants.
func BadgeTypeIsValid(v BadgeType) bool {
	_, ok := badgeStyleMap[v]
	return ok
}

// BadgeSizeIsValid reports whether v is one of the defined BadgeSize constants.
func BadgeSizeIsValid(v BadgeSize) bool {
	_, ok := badgeSizeLookup[v]
	return ok
}

// CardPaddingIsValid reports whether v is one of the defined CardPadding constants.
func CardPaddingIsValid(v CardPadding) bool {
	_, ok := cardPaddingLookup[v]
	return ok
}

// GridColsIsValid reports whether v is one of the defined GridCols constants.
func GridColsIsValid(v GridCols) bool {
	_, ok := gridColsLookup[v]
	return ok
}

// GridGapIsValid reports whether v is one of the defined GridGap constants.
func GridGapIsValid(v GridGap) bool {
	_, ok := gridGapLookup[v]
	return ok
}

// TrendDirectionIsValid reports whether v is one of the defined TrendDirection constants.
func TrendDirectionIsValid(v TrendDirection) bool {
	return validTrendDirections[v]
}

// AvatarSizeIsValid reports whether v is one of the defined AvatarSize constants.
func AvatarSizeIsValid(v AvatarSize) bool {
	_, ok := avatarSizeLookup[v]
	return ok
}

// TooltipPositionIsValid reports whether v is one of the defined TooltipPosition constants.
func TooltipPositionIsValid(v TooltipPosition) bool {
	_, ok := tooltipPositionMap[v]
	return ok
}

// PopoverPositionIsValid reports whether v is one of the defined PopoverPosition constants.
func PopoverPositionIsValid(v PopoverPosition) bool {
	_, ok := popoverPositionMap[v]
	return ok
}

// AvatarShapeIsValid reports whether v is one of the defined AvatarShape constants.
func AvatarShapeIsValid(v AvatarShape) bool {
	switch v {
	case AvatarShapeCircle, AvatarShapeSquare:
		return true
	default:
		return false
	}
}

// AvatarStatusIsValid reports whether v is one of the defined AvatarStatus constants.
func AvatarStatusIsValid(v AvatarStatus) bool {
	switch v {
	case AvatarStatusNone, AvatarStatusOnline, AvatarStatusOffline:
		return true
	default:
		return false
	}
}

// DropdownItemKindIsValid reports whether v is one of the defined DropdownItemKind constants.
func DropdownItemKindIsValid(v DropdownItemKind) bool {
	switch v {
	case DropdownItemLink, DropdownItemButton:
		return true
	default:
		return false
	}
}

// DropdownPositionIsValid reports whether v is one of the defined DropdownPosition constants.
func DropdownPositionIsValid(v DropdownPosition) bool {
	switch v {
	case DropdownPositionLeft, DropdownPositionRight:
		return true
	default:
		return false
	}
}

// TabsVariantIsValid reports whether v is one of the defined TabsVariant constants.
func TabsVariantIsValid(v TabsVariant) bool {
	switch v {
	case TabsDefault, TabsPills:
		return true
	default:
		return false
	}
}

// OverlayKindIsValid reports whether v is one of the defined OverlayKind constants.
func OverlayKindIsValid(v OverlayKind) bool {
	switch v {
	case OverlayModal, OverlayDrawer:
		return true
	default:
		return false
	}
}

// ButtonSizeIsValid reports whether v is one of the defined ButtonSize constants.
func ButtonSizeIsValid(v ButtonSize) bool {
	_, ok := buttonSizeLookup[v]
	return ok
}

// ButtonHTMLTypeIsValid reports whether v is one of the defined ButtonHTMLType constants.
func ButtonHTMLTypeIsValid(v ButtonHTMLType) bool {
	_, ok := buttonHTMLTypeLookup[v]
	return ok
}

// SortDirectionIsValid reports whether v is one of the defined SortDirection constants.
func SortDirectionIsValid(v SortDirection) bool {
	switch v {
	case SortNone, SortAsc, SortDesc:
		return true
	default:
		return false
	}
}
