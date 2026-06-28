// Package feedback provides user feedback components such as alerts, toasts, and loading indicators.
package feedback

import (
	"github.com/larsartmann/templ-components/icons"
)

// FeedbackType represents the severity/visual style of a feedback component.
// Shared by Alert and Toast components.
type FeedbackType string

// Feedback severity levels shared by Alert and Toast components.
const (
	FeedbackSuccess FeedbackType = "success"
	FeedbackError   FeedbackType = "error"
	FeedbackWarning FeedbackType = "warning"
	FeedbackInfo    FeedbackType = "info"
)

// feedbackStyleSet holds the CSS classes for a feedback component variant.
type feedbackStyleSet struct {
	Border, BG, Text, Icon string
}

// lookupFeedbackStyle returns the style set for key t from map m,
// or def if not found.
func lookupFeedbackStyle[T ~string](
	m map[T]feedbackStyleSet,
	def feedbackStyleSet,
	t T,
) feedbackStyleSet {
	if s, ok := m[t]; ok {
		return s
	}
	return def
}

// feedbackIconName maps a FeedbackType to its canonical icon name.
// Shared by Alert and Toast components.
func feedbackIconName(m map[FeedbackType]icons.Name, t FeedbackType) icons.Name {
	if n, ok := m[t]; ok {
		return n
	}
	return icons.Information
}

// feedbackStyleMap is the single source of truth for feedback styles.
// Shared by Alert and Toast to guarantee visual consistency for the same severity.
//
//nolint:gochecknoglobals,goconst // Package-level lookup table; Tailwind class strings are intentionally inline
var feedbackStyleMap = map[FeedbackType]feedbackStyleSet{
	FeedbackSuccess: {
		Border: "border-green-200 dark:border-green-800",
		BG:     "bg-green-50 dark:bg-green-900/20",
		Text:   "text-green-800 dark:text-green-200",
		Icon:   "text-green-600 dark:text-green-400",
	},
	FeedbackError: {
		Border: "border-red-200 dark:border-red-800",
		BG:     "bg-red-50 dark:bg-red-900/20",
		Text:   "text-red-800 dark:text-red-200",
		Icon:   "text-red-600 dark:text-red-400",
	},
	FeedbackWarning: {
		Border: "border-yellow-200 dark:border-yellow-800",
		BG:     "bg-yellow-50 dark:bg-yellow-900/20",
		Text:   "text-yellow-800 dark:text-yellow-200",
		Icon:   "text-yellow-600 dark:text-yellow-400",
	},
	FeedbackInfo: {
		Border: "border-blue-200 dark:border-blue-800",
		BG:     "bg-blue-50 dark:bg-blue-900/20",
		Text:   "text-blue-800 dark:text-blue-200",
		Icon:   "text-blue-600 dark:text-blue-400",
	},
}

//nolint:gochecknoglobals // Package-level default fallback
var feedbackStyleDefault = feedbackStyleSet{
	Border: "border-gray-200 dark:border-gray-700",
	BG:     "bg-gray-50 dark:bg-gray-800/50",
	Text:   "text-gray-800 dark:text-gray-200",
	Icon:   "text-gray-600 dark:text-gray-400",
}

// feedbackIconMap is the single source of truth for feedback icons.
// Uses circle-style icons for clear completion/failure semantics.
//
//nolint:gochecknoglobals // Package-level lookup table for shared feedback icons
var feedbackIconMap = map[FeedbackType]icons.Name{
	FeedbackSuccess: icons.CheckCircle,
	FeedbackError:   icons.XCircle,
	FeedbackWarning: icons.ExclamationTriangle,
	FeedbackInfo:    icons.Information,
}

func feedbackStyle(t FeedbackType) feedbackStyleSet {
	return lookupFeedbackStyle(feedbackStyleMap, feedbackStyleDefault, t)
}

func feedbackIcon(t FeedbackType) icons.Name {
	return feedbackIconName(feedbackIconMap, t)
}
