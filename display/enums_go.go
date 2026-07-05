package display

// BadgeTypeIsValid reports whether v is one of the defined BadgeType constants.
func BadgeTypeIsValid(v BadgeType) bool {
	_, ok := badgeStyleMap[v]
	return ok
}

// BadgeSizeIsValid reports whether v is one of the defined BadgeSize constants.
func BadgeSizeIsValid(v BadgeSize) bool {
	_, ok := badgeSizeLookup[string(v)]
	return ok
}

// CardPaddingIsValid reports whether v is one of the defined CardPadding constants.
func CardPaddingIsValid(v CardPadding) bool {
	_, ok := cardPaddingLookup[string(v)]
	return ok
}

// GridColsIsValid reports whether v is one of the defined GridCols constants.
func GridColsIsValid(v GridCols) bool {
	_, ok := gridColsLookup[v]
	return ok
}

// TrendDirectionIsValid reports whether v is one of the defined TrendDirection constants.
func TrendDirectionIsValid(v TrendDirection) bool {
	return validTrendDirections[v]
}

// AvatarSizeIsValid reports whether v is one of the defined AvatarSize constants.
func AvatarSizeIsValid(v AvatarSize) bool {
	_, ok := avatarSizeLookup[string(v)]
	return ok
}

// TooltipPositionIsValid reports whether v is one of the defined TooltipPosition constants.
func TooltipPositionIsValid(v TooltipPosition) bool {
	_, ok := tooltipPositionMap[v]
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
