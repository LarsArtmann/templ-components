package errorpage

import (
	"errors"
	"fmt"
	"strings"

	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

// Family classifies an error's behavioral profile for web presentation.
// Mirrors the go-error-family library's 6 families — consumers bridge with trivial string constants.
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
	// Tone: apologetic. Visual: gray.
	FamilyInfrastructure Family = "infrastructure"

	// FamilyOrchestration indicates an internal coordination failure (bug, misconfiguration).
	// Tone: factual. Visual: purple.
	FamilyOrchestration Family = "orchestration"
)

// familyVisualStyle holds the CSS classes and icon for a family's visual treatment.
type familyVisualStyle struct {
	Border       string
	BG           string
	Bar          string
	Text         string
	IconColor    string
	Icon         icons.Name
	AccentBG     string
	AccentText   string
	ActionButton string
	// ActionButtonGhost is the secondary (outline) action style — family-
	// tinted border/text on a neutral background, used by SecondaryWayOut.
	ActionButtonGhost string
}

//nolint:gochecknoglobals // Package-level lookup table for family visual styles
var familyStyleMap = map[Family]familyVisualStyle{
	FamilyRejection: {
		Border:       "border-amber-200 dark:border-amber-800",
		BG:           "bg-amber-50 dark:bg-amber-900/20",
		Bar:          "bg-amber-400 dark:bg-amber-500",
		Text:         "text-amber-800 dark:text-amber-200",
		IconColor:    "text-amber-500 dark:text-amber-400",
		Icon:         icons.ExclamationTriangle,
		AccentBG:     "bg-amber-100 dark:bg-amber-900/30",
		AccentText:   "text-amber-700 dark:text-amber-300",
		ActionButton: "bg-amber-600 hover:bg-amber-500 dark:bg-amber-500 dark:hover:bg-amber-400 focus-visible:ring-amber-500 dark:focus-visible:ring-amber-400 text-white",
		ActionButtonGhost: "border border-amber-300 dark:border-amber-700 bg-white dark:bg-transparent text-amber-700 dark:text-amber-300 hover:bg-amber-50 dark:hover:bg-amber-900/30 focus-visible:ring-amber-500 dark:focus-visible:ring-amber-400",
	},
	FamilyConflict: {
		Border:       "border-orange-200 dark:border-orange-800",
		BG:           "bg-orange-50 dark:bg-orange-900/20",
		Bar:          "bg-orange-400 dark:bg-orange-500",
		Text:         "text-orange-800 dark:text-orange-200",
		IconColor:    "text-orange-500 dark:text-orange-400",
		Icon:         icons.ExclamationCircle,
		AccentBG:     "bg-orange-100 dark:bg-orange-900/30",
		AccentText:   "text-orange-700 dark:text-orange-300",
		ActionButton: "bg-orange-600 hover:bg-orange-500 dark:bg-orange-500 dark:hover:bg-orange-400 focus-visible:ring-orange-500 dark:focus-visible:ring-orange-400 text-white",
		ActionButtonGhost: "border border-orange-300 dark:border-orange-700 bg-white dark:bg-transparent text-orange-700 dark:text-orange-300 hover:bg-orange-50 dark:hover:bg-orange-900/30 focus-visible:ring-orange-500 dark:focus-visible:ring-orange-400",
	},
	FamilyTransient: {
		Border:       "border-blue-200 dark:border-blue-800",
		BG:           "bg-blue-50 dark:bg-blue-900/20",
		Bar:          "bg-blue-500 dark:bg-blue-400",
		Text:         "text-blue-800 dark:text-blue-200",
		IconColor:    "text-blue-500 dark:text-blue-400",
		Icon:         icons.Refresh,
		AccentBG:     "bg-blue-100 dark:bg-blue-900/30",
		AccentText:   "text-blue-700 dark:text-blue-300",
		ActionButton: "bg-blue-600 hover:bg-blue-500 dark:bg-blue-500 dark:hover:bg-blue-400 focus-visible:ring-blue-500 dark:focus-visible:ring-blue-400 text-white",
		ActionButtonGhost: "border border-blue-300 dark:border-blue-700 bg-white dark:bg-transparent text-blue-700 dark:text-blue-300 hover:bg-blue-50 dark:hover:bg-blue-900/30 focus-visible:ring-blue-500 dark:focus-visible:ring-blue-400",
	},
	FamilyCorruption: {
		Border:       "border-red-200 dark:border-red-800",
		BG:           "bg-red-50 dark:bg-red-900/20",
		Bar:          "bg-red-500 dark:bg-red-400",
		Text:         "text-red-800 dark:text-red-200",
		IconColor:    "text-red-500 dark:text-red-400",
		Icon:         icons.ExclamationTriangle,
		AccentBG:     "bg-red-100 dark:bg-red-900/30",
		AccentText:   "text-red-700 dark:text-red-300",
		ActionButton: "bg-red-600 hover:bg-red-500 dark:bg-red-500 dark:hover:bg-red-400 focus-visible:ring-red-500 dark:focus-visible:ring-red-400 text-white",
		ActionButtonGhost: "border border-red-300 dark:border-red-700 bg-white dark:bg-transparent text-red-700 dark:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/30 focus-visible:ring-red-500 dark:focus-visible:ring-red-400",
	},
	FamilyInfrastructure: {
		Border:       "border-gray-200 dark:border-gray-700",
		BG:           "bg-gray-50 dark:bg-gray-800/50",
		Bar:          "bg-gray-400 dark:bg-gray-500",
		Text:         "text-gray-800 dark:text-gray-200",
		IconColor:    "text-gray-400 dark:text-gray-500",
		Icon:         icons.Globe,
		AccentBG:     "bg-gray-100 dark:bg-gray-800",
		AccentText:   "text-gray-700 dark:text-gray-300",
		ActionButton: "bg-gray-600 hover:bg-gray-500 dark:bg-gray-500 dark:hover:bg-gray-400 focus-visible:ring-gray-500 dark:focus-visible:ring-gray-400 text-white",
		ActionButtonGhost: "border border-gray-300 dark:border-gray-600 bg-white dark:bg-transparent text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700/40 focus-visible:ring-gray-500 dark:focus-visible:ring-gray-400",
	},
	FamilyOrchestration: {
		Border:       "border-purple-200 dark:border-purple-800",
		BG:           "bg-purple-50 dark:bg-purple-900/20",
		Bar:          "bg-purple-500 dark:bg-purple-400",
		Text:         "text-purple-800 dark:text-purple-200",
		IconColor:    "text-purple-500 dark:text-purple-400",
		Icon:         icons.ExclamationTriangle,
		AccentBG:     "bg-purple-100 dark:bg-purple-900/30",
		AccentText:   "text-purple-700 dark:text-purple-300",
		ActionButton: "bg-purple-600 hover:bg-purple-500 dark:bg-purple-500 dark:hover:bg-purple-400 focus-visible:ring-purple-500 dark:focus-visible:ring-purple-400 text-white",
		ActionButtonGhost: "border border-purple-300 dark:border-purple-700 bg-white dark:bg-transparent text-purple-700 dark:text-purple-300 hover:bg-purple-50 dark:hover:bg-purple-900/30 focus-visible:ring-purple-500 dark:focus-visible:ring-purple-400",
	},
}

//nolint:gochecknoglobals // Package-level default fallback
var familyStyleDefault = familyVisualStyle{
	Border:       "border-gray-200 dark:border-gray-700",
	BG:           "bg-gray-50 dark:bg-gray-800/50",
	Bar:          "bg-gray-400 dark:bg-gray-500",
	Text:         "text-gray-800 dark:text-gray-200",
	IconColor:    "text-gray-400 dark:text-gray-500",
	Icon:         icons.Information,
	AccentBG:     "bg-gray-100 dark:bg-gray-800",
	AccentText:   "text-gray-700 dark:text-gray-300",
	ActionButton: "bg-gray-600 hover:bg-gray-500 dark:bg-gray-500 dark:hover:bg-gray-400 focus-visible:ring-gray-500 dark:focus-visible:ring-gray-400 text-white",
	ActionButtonGhost: "border border-gray-300 dark:border-gray-600 bg-white dark:bg-transparent text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700/40 focus-visible:ring-gray-500 dark:focus-visible:ring-gray-400",
}

func lookupFamilyStyle(f Family) familyVisualStyle {
	return utils.Lookup(familyStyleMap, f, familyStyleDefault)
}

// FamilyIcon returns the icon name for a given family.
func FamilyIcon(f Family) icons.Name {
	return lookupFamilyStyle(f).Icon
}

// FamilyIsValid reports whether the Family value is one of the six defined constants.
func FamilyIsValid(f Family) bool {
	_, ok := familyStyleMap[f]

	return ok
}

// familyDefaultTitleMap holds the family-derived page title used when an
// error carries no title of its own. Tone-matched to go-error-family's
// per-family DefaultWhy/DefaultFix copy: short, neutral, non-jargony —
// the title heads the page while the message carries the detail.
//
//nolint:gochecknoglobals // Package-level lookup table, mirrors familyStyleMap
var familyDefaultTitleMap = map[Family]string{
	FamilyRejection:      "Request could not be completed",
	FamilyConflict:       "Conflict detected",
	FamilyTransient:      "Temporary error",
	FamilyCorruption:     "Data integrity issue",
	FamilyInfrastructure: "Service unavailable",
	FamilyOrchestration:  "Internal error",
}

// FamilyDefaultTitle returns the family-derived fallback title for a full
// error page. Unknown families return the infrastructure (most apologetic)
// title so a page never renders headingless.
func FamilyDefaultTitle(f Family) string {
	return utils.Lookup(familyDefaultTitleMap, f, "Service unavailable")
}

// Diagnostic panel surfaces for the shared fix/context helpers. ErrorPage's
// white card uses the neutral inset; ErrorDetail's family-tinted card keeps
// the white inset.
const (
	errorInsetNeutral = "bg-gray-50/60 border-gray-200/80 dark:bg-gray-800/50 dark:border-gray-700/60"
	errorInsetCard    = "bg-white border-gray-200 dark:bg-gray-800 dark:border-gray-700"
)

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
	Code    Code
}

// ErrorPageProps configures a full-page error view.
type ErrorPageProps struct {
	utils.BaseProps

	Family     Family
	StatusCode int
	Code       Code
	Title      string
	Message    string
	Why        string
	Fix        string
	WayOut     string
	WayOutHref string
	// WayOutAction is the typed bundle for the primary action. When Text is
	// set it wins entirely over the loose WayOut/WayOutHref strings (no
	// mixing) — callers may adopt it without touching existing fields.
	WayOutAction WayOutAction
	// SecondaryWayOut renders a secondary (ghost) action next to the primary
	// way out ("View status page", "Contact support"). With a
	// SecondaryWayOutHref it renders as a link; without one it behaves like
	// the primary's no-href variant (go back in history).
	SecondaryWayOut     string
	SecondaryWayOutHref string
	// MaxWidth caps the card width. Empty/unknown renders MaxWidthXL (the
	// pre-field look).
	MaxWidth     ErrorMaxWidth
	Context      []ContextPair
	CauseChain   []CauseItem
	Timestamp    string
	Trace        string
	ShowTimestamp bool
}

// WayOutAction bundles a recovery action's label with its destination. Href
// empty means a history-back button (the browser-native way out).
type WayOutAction struct {
	Text string
	Href string
}

// ErrorMaxWidth selects the ErrorPage card's max width. It is a closed-set
// enum; unknown values fall back to MaxWidthXL (map+fallback convention).
type ErrorMaxWidth string

const (
	ErrorMaxWidthLG  ErrorMaxWidth = "lg"
	ErrorMaxWidthXL  ErrorMaxWidth = "xl"
	ErrorMaxWidth2XL ErrorMaxWidth = "2xl"
	ErrorMaxWidth4XL ErrorMaxWidth = "4xl"
)

// ErrorMaxWidthIsValid reports whether w is a known ErrorMaxWidth.
func ErrorMaxWidthIsValid(w ErrorMaxWidth) bool {
	switch w {
	case ErrorMaxWidthLG, ErrorMaxWidthXL, ErrorMaxWidth2XL, ErrorMaxWidth4XL:
		return true
	default:
		return false
	}
}

//nolint:gochecknoglobals // Package-level lookup table
var errorMaxWidthClassMap = map[ErrorMaxWidth]string{
	ErrorMaxWidthLG:  "max-w-lg",
	ErrorMaxWidthXL:  "max-w-xl",
	ErrorMaxWidth2XL: "max-w-2xl",
	ErrorMaxWidth4XL: "max-w-4xl",
}

// errorMaxWidthClass resolves the Tailwind max-width class, defaulting to XL.
func errorMaxWidthClass(w ErrorMaxWidth) string {
	return utils.Lookup(errorMaxWidthClassMap, w, "max-w-xl")
}

// resolvedWayOut returns the effective primary action: WayOutAction wins
// entirely when its Text is set; otherwise the loose strings apply.
func (p ErrorPageProps) resolvedWayOut() (string, string) {
	if p.WayOutAction.Text != "" {
		return p.WayOutAction.Text, p.WayOutAction.Href
	}

	return p.WayOut, p.WayOutHref
}

// errBlankNonRejection is the Validate error returned when the props would
// render as an empty error card.
var errBlankNonRejection = errValidateBlank

// Validate verifies that the props form a coherent error page. Returns an
// error when:
//   - Family is not one of the six defined constants (FamilyIsValid).
//   - StatusCode is set but outside the HTTP error range [400, 599].
//   - The page has no Title AND no Message AND no CauseChain (would render
//     as an empty error card — likely a caller bug).
//
// Validate is intentionally permissive about optional fields (Why, Fix,
// WayOut, Context) — those are presentation, not correctness.
func (p ErrorPageProps) Validate() error {
	if !FamilyIsValid(p.Family) {
		return fmt.Errorf("%w: %q", errValidateFamily, p.Family)
	}

	if p.StatusCode != 0 && (p.StatusCode < 400 || p.StatusCode > 599) {
		return fmt.Errorf("%w: %d", errValidateStatusRange, p.StatusCode)
	}

	if p.Title == "" && p.Message == "" && len(p.CauseChain) == 0 {
		return errBlankNonRejection
	}

	return nil
}

// Validation errors. Kept as package-level errors so callers can errors.Is
// against them.
var (
	errValidateFamily      = errors.New("errorpage: invalid Family")
	errValidateStatusRange = errors.New("errorpage: StatusCode must be in [400, 599]")
	errValidateBlank       = errors.New("errorpage: at least one of Title, Message, CauseChain must be set")
)

// DefaultErrorPageProps returns sensible defaults.
func DefaultErrorPageProps() ErrorPageProps {
	return ErrorPageProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Family: FamilyTransient,
	}
}

// ErrorDetailVariant selects the card shell treatment for ErrorDetail.
// Unknown or empty values render Tinted (the pre-variant look), matching the
// library's map+fallback convention — zero-value props keep rendering
// exactly as before the variant existed.
type ErrorDetailVariant string

const (
	// ErrorDetailTinted (default) keeps the family-tinted background and
	// border — the compact alert-like card.
	ErrorDetailTinted ErrorDetailVariant = "tinted"
	// ErrorDetailNeutral renders a neutral card shell with a family-colored
	// accent bar on top — visual parity with the redesigned ErrorPage, for
	// inline placement in content that already carries color (dashboards,
	// side panels, tinted sections).
	ErrorDetailNeutral ErrorDetailVariant = "neutral"
)

// ErrorDetailVariantIsValid reports whether v is a known ErrorDetailVariant.
func ErrorDetailVariantIsValid(v ErrorDetailVariant) bool {
	switch v {
	case ErrorDetailTinted, ErrorDetailNeutral:
		return true
	default:
		return false
	}
}

// errorDetailShellClasses returns the border/background classes for the
// ErrorDetail card shell under the given variant.
func errorDetailShellClasses(v ErrorDetailVariant, style familyVisualStyle) string {
	if v == ErrorDetailNeutral {
		return "bg-white border-gray-200 dark:bg-gray-800 dark:border-gray-700"
	}

	return style.Border + " " + style.BG
}

// ErrorDetailProps configures an inline error detail card.
type ErrorDetailProps struct {
	utils.BaseProps

	Family     Family
	Code       Code
	Title      string
	Message    string
	Fix        string
	Context    []ContextPair
	CauseChain []CauseItem
	Timestamp  string
	Trace      string
	// Variant selects the card shell (Tinted default; Neutral adds a family
	// accent bar over a neutral shell). Empty/unknown renders Tinted.
	Variant ErrorDetailVariant
}

// DefaultErrorDetailProps returns sensible defaults.
func DefaultErrorDetailProps() ErrorDetailProps {
	return ErrorDetailProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Family:  FamilyTransient,
		Variant: ErrorDetailTinted,
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
	return ErrorAlertProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Family: FamilyTransient,
	}
}

// FamilyStatusCode returns the HTTP status code for a family.
// Useful for HTTP handlers that need to set the correct response status.
func FamilyStatusCode(f Family) int {
	return utils.Lookup(familyStatusCodeMap, f, 500)
}

//nolint:gochecknoglobals // Package-level lookup table
var familyStatusCodeMap = map[Family]int{
	FamilyRejection:      400,
	FamilyConflict:       409,
	FamilyTransient:      503,
	FamilyCorruption:     500,
	FamilyInfrastructure: 503,
	FamilyOrchestration:  500,
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
// Handles both single-error Unwrap() chains and errors.Join siblings
// (Unwrap() []error, Go 1.20+). Stops after maxDepth levels to prevent
// infinite chains.
func ExtractCauseChain(err error, maxDepth int) []CauseItem {
	if err == nil || maxDepth <= 0 {
		return nil
	}

	var chain []CauseItem

	current := err
	for range maxDepth {
		unwrapped := errors.Unwrap(current)
		if unwrapped == nil {
			chain = appendJoinSiblings(chain, current, maxDepth)

			break
		}

		chain = append(chain, causeItemFromError(unwrapped))
		current = unwrapped
	}

	return chain
}

// appendJoinSiblings adds errors.Join siblings (Unwrap() []error) to the chain,
// respecting maxDepth. No-op if current does not implement Unwrap() []error.
func appendJoinSiblings(chain []CauseItem, current error, maxDepth int) []CauseItem {
	joiner, ok := current.(interface{ Unwrap() []error })
	if !ok {
		return chain
	}

	for _, sibling := range joiner.Unwrap() {
		if len(chain) >= maxDepth {
			break
		}

		chain = append(chain, causeItemFromError(sibling))
	}

	return chain
}

// causeItemFromError builds a CauseItem from an error, extracting ErrorCode if available.
func causeItemFromError(err error) CauseItem {
	item := CauseItem{Message: err.Error()} //nolint:exhaustruct_v5 // Code set conditionally below
	if c, ok := err.(interface{ ErrorCode() string }); ok {
		item.Code = Code(c.ErrorCode())
	}

	return item
}
