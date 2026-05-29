package errorpage

import (
	"errors"
	"strings"

	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

// Family classifies an error's behavioral profile for web presentation.
// Mirrors the go-error-family library's 5 families — consumers bridge with trivial string constants.
//
// Each family maps to a distinct visual treatment (color, icon, tone) that communicates
// the error's nature to the user without technical jargon.
type Family string

const (
	// FamilyRejection indicates bad input, unauthorized access, or resource not found.
	// Tone: helpful, instructional. Visual: amber.
	FamilyRejection Family = "rejection"

	// FamilyConflict indicates version mismatch, duplicate creation, or state machine violation.
	// Tone: explanatory. Visual: orange.
	FamilyConflict Family = "conflict"

	// FamilyTransient indicates a temporary infrastructure failure.
	// Tone: reassuring. Visual: blue.
	FamilyTransient Family = "transient"

	// FamilyCorruption indicates the source of truth is damaged.
	// Tone: urgent. Visual: red.
	FamilyCorruption Family = "corruption"

	// FamilyInfrastructure indicates the system cannot serve.
	// Tone: apologetic. Visual: slate/gray.
	FamilyInfrastructure Family = "infrastructure"
)

// familyVisualStyle holds the CSS classes and icon for a family's visual treatment.
type familyVisualStyle struct {
	Border       string
	BG           string
	Text         string
	IconColor    string
	Icon         icons.Name
	AccentBG     string
	AccentText   string
	ActionButton string
}

//nolint:gochecknoglobals // Package-level lookup table for family visual styles
var familyStyleMap = map[Family]familyVisualStyle{
	FamilyRejection: {
		Border: "border-amber-200 dark:border-amber-800", BG: "bg-amber-50 dark:bg-amber-900/20",
		Text: "text-amber-800 dark:text-amber-200", IconColor: "text-amber-500",
		Icon:     icons.ExclamationTriangle,
		AccentBG: "bg-amber-100 dark:bg-amber-900/30", AccentText: "text-amber-700 dark:text-amber-300",
		ActionButton: "bg-amber-600 hover:bg-amber-500 focus-visible:ring-amber-500 text-white",
	},
	FamilyConflict: {
		Border: "border-orange-200 dark:border-orange-800", BG: "bg-orange-50 dark:bg-orange-900/20",
		Text: "text-orange-800 dark:text-orange-200", IconColor: "text-orange-500",
		Icon:     icons.ExclamationCircle,
		AccentBG: "bg-orange-100 dark:bg-orange-900/30", AccentText: "text-orange-700 dark:text-orange-300",
		ActionButton: "bg-orange-600 hover:bg-orange-500 focus-visible:ring-orange-500 text-white",
	},
	FamilyTransient: {
		Border: "border-blue-200 dark:border-blue-800", BG: "bg-blue-50 dark:bg-blue-900/20",
		Text: "text-blue-800 dark:text-blue-200", IconColor: "text-blue-500",
		Icon:     icons.Refresh,
		AccentBG: "bg-blue-100 dark:bg-blue-900/30", AccentText: "text-blue-700 dark:text-blue-300",
		ActionButton: "bg-blue-600 hover:bg-blue-500 focus-visible:ring-blue-500 text-white",
	},
	FamilyCorruption: {
		Border: "border-red-200 dark:border-red-800", BG: "bg-red-50 dark:bg-red-900/20",
		Text: "text-red-800 dark:text-red-200", IconColor: "text-red-500",
		Icon:     icons.ExclamationTriangle,
		AccentBG: "bg-red-100 dark:bg-red-900/30", AccentText: "text-red-700 dark:text-red-300",
		ActionButton: "bg-red-600 hover:bg-red-500 focus-visible:ring-red-500 text-white",
	},
	FamilyInfrastructure: {
		Border: "border-slate-200 dark:border-slate-700", BG: "bg-slate-50 dark:bg-slate-800/50",
		Text: "text-slate-800 dark:text-slate-200", IconColor: "text-slate-400",
		Icon:     icons.Globe,
		AccentBG: "bg-slate-100 dark:bg-slate-800", AccentText: "text-slate-700 dark:text-slate-300",
		ActionButton: "bg-slate-600 hover:bg-slate-500 focus-visible:ring-slate-500 text-white",
	},
}

//nolint:gochecknoglobals // Package-level default fallback
var familyStyleDefault = familyVisualStyle{
	Border: "border-gray-200 dark:border-gray-700", BG: "bg-gray-50 dark:bg-gray-800/50",
	Text: "text-gray-800 dark:text-gray-200", IconColor: "text-gray-400",
	Icon:     icons.Information,
	AccentBG: "bg-gray-100 dark:bg-gray-800", AccentText: "text-gray-700 dark:text-gray-300",
	ActionButton: "bg-gray-600 hover:bg-gray-500 focus-visible:ring-gray-500 text-white",
}

func lookupFamilyStyle(f Family) familyVisualStyle {
	if s, ok := familyStyleMap[f]; ok {
		return s
	}
	return familyStyleDefault
}

// FamilyIcon returns the icon name for a given family.
func FamilyIcon(f Family) icons.Name {
	return lookupFamilyStyle(f).Icon
}

// FamilyIsValid reports whether the Family value is one of the five defined constants.
func FamilyIsValid(f Family) bool {
	_, ok := familyStyleMap[f]
	return ok
}

// ParseFamily parses a family string (case-insensitive) into a Family.
// Returns FamilyTransient for unrecognized values.
func ParseFamily(s string) Family {
	f := Family(strings.ToLower(strings.TrimSpace(s)))
	if FamilyIsValid(f) {
		return f
	}
	return FamilyTransient
}

// ContextPair is a key-value pair from an error's context map.
type ContextPair struct {
	Key   string
	Value string
}

// CauseItem represents one error in a cause chain.
type CauseItem struct {
	Message string
	Code    string
}

// ErrorPageProps configures a full-page error view.
type ErrorPageProps struct {
	utils.BaseProps
	Family        Family
	Code          string
	Title         string
	Message       string
	Why           string
	Fix           string
	WayOut        string
	WayOutHref    string
	Context       []ContextPair
	CauseChain    []CauseItem
	Timestamp     string
	ShowTimestamp bool
}

// DefaultErrorPageProps returns sensible defaults.
func DefaultErrorPageProps() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct // intentionally minimal defaults
		Family: FamilyTransient,
	}
}

// ErrorDetailProps configures an inline error detail card.
type ErrorDetailProps struct {
	utils.BaseProps
	Family     Family
	Code       string
	Title      string
	Message    string
	Fix        string
	Context    []ContextPair
	CauseChain []CauseItem
	Timestamp  string
}

// DefaultErrorDetailProps returns sensible defaults.
func DefaultErrorDetailProps() ErrorDetailProps {
	return ErrorDetailProps{ //nolint:exhaustruct // intentionally minimal defaults
		Family: FamilyTransient,
	}
}

// ErrorAlertProps configures an alert banner derived from an error family.
type ErrorAlertProps struct {
	utils.BaseProps
	Family      Family
	Title       string
	Message     string
	Fix         string
	Dismissible bool
}

// DefaultErrorAlertProps returns sensible defaults.
func DefaultErrorAlertProps() ErrorAlertProps {
	return ErrorAlertProps{ //nolint:exhaustruct // intentionally minimal defaults
		Family: FamilyTransient,
	}
}

// FamilyStatusCode returns the HTTP status code for a family.
// Useful for HTTP handlers that need to set the correct response status.
func FamilyStatusCode(f Family) int {
	if code, ok := familyStatusCodeMap[f]; ok {
		return code
	}
	return 500
}

//nolint:gochecknoglobals // Package-level lookup table
var familyStatusCodeMap = map[Family]int{
	FamilyRejection:      400,
	FamilyConflict:       409,
	FamilyTransient:      503,
	FamilyCorruption:     500,
	FamilyInfrastructure: 503,
}

// ContextMap converts a map[string]string to a []ContextPair slice.
// Useful for bridging go-error-family's ErrorContext() to errorpage props.
func ContextMap(m map[string]string) []ContextPair {
	if len(m) == 0 {
		return nil
	}
	pairs := make([]ContextPair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, ContextPair{Key: k, Value: v})
	}
	return pairs
}

// ExtractCauseChain walks an error's Unwrap() chain and returns CauseItems.
// Useful for bridging go-error-family errors to errorpage props.
// Stops after maxDepth levels to prevent infinite chains.
func ExtractCauseChain(err error, maxDepth int) []CauseItem {
	if err == nil || maxDepth <= 0 {
		return nil
	}
	var chain []CauseItem
	current := err
	for range maxDepth {
		unwrapped := errors.Unwrap(current)
		if unwrapped == nil {
			break
		}
		item := CauseItem{Message: unwrapped.Error()} //nolint:exhaustruct // Code set conditionally below
		if c, ok := unwrapped.(interface{ ErrorCode() string }); ok {
			item.Code = c.ErrorCode()
		}
		chain = append(chain, item)
		current = unwrapped
	}
	return chain
}
