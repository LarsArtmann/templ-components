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
