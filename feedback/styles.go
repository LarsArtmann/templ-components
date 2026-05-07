// Package feedback provides user feedback components such as alerts, toasts, and loading indicators.
package feedback

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
