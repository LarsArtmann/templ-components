// Feedback style lookups: shared style sets, icon mapping, and generic enum lookup.
package feedback

import (
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/tone"
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

// feedbackStyle derives the style set for t from the canonical tone table
// (utils/tone) — the ADR-0049 collapse: the former feedbackStyleMap was
// the tone table's anchor dialect, and TestToneAdapterEquivalence keeps
// those byte-for-byte expectations as literals, so the derived classes
// still ship exactly the strings feedback has always rendered. Semantic
// severities stay palette classes (ADR-0044 decision 3); only the brand
// axis emits token families.
func feedbackStyle(t FeedbackType) feedbackStyleSet {
	s := tone.For(t.Tone())

	return feedbackStyleSet{
		Border: s.Border.Classes(),
		BG:     s.BGTint.Classes(),
		Text:   s.Text.Classes(),
		Icon:   s.Icon.Classes(),
	}
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

// feedbackRole maps the feedback type onto the ARIA role: errors announce
// assertively (role="alert"); every other severity reports politely
// (role="status") so success and info banners are not read as errors by
// screen readers. Mirrors the inlineMessage twins' split (InlineError
// "alert" vs InlineSuccess "status").
func feedbackRole(t FeedbackType) string {
	if t == FeedbackError {
		return "alert"
	}

	return "status"
}

func feedbackIcon(t FeedbackType) icons.Name {
	return utils.Lookup(feedbackIconMap, t, icons.Information)
}

// FeedbackTypeIsValid reports whether the FeedbackType value is one of the
// four defined constants (Success, Error, Warning, Info).
func FeedbackTypeIsValid(t FeedbackType) bool {
	switch t {
	case FeedbackSuccess, FeedbackError, FeedbackWarning, FeedbackInfo:
		return true
	default:
		return false
	}
}
