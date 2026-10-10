// FilterBar component: the horizontal auto-submit filter row for list
// pages — the converged shape behind DiscordSync's filterForm (9 list pages)
// and dnsblockd's filter rows, extracted per the 2026-10-10 four-project
// analysis (TODO_LIST #390). Composes FilterDropdown/FilterInput/Checkbox
// children; see docs/recipes/horizontal-filter-bar.md for the pattern this
// component productizes (including the two battle-tested footguns it encodes).
package forms

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// FilterBarProps configures a horizontal auto-submit filter bar. The bar is
// a GET form: children render the filter controls (FilterDropdown selects,
// checkboxes, search inputs, hidden inputs); any select or checkbox change
// re-submits the form without a page reload, and the noscript Apply button
// keeps non-JS browsers working with a full-page load.
type FilterBarProps struct {
	utils.BaseProps

	// Action is the list-page URL: the form's GET action (the no-JS
	// fallback), the Reset link target, and — when Wire is unset — the
	// auto-submit fetch URL. An empty Action renders an unwired form
	// (inert bar; useful for pure-CSS composition).
	Action string
	// Wire, when set with a non-empty URL, owns the auto-submit wiring in
	// both dialects. Defaults applied by the component: the change event
	// plus form encoding, so every filter field in the form serializes.
	// Target/PushURL render for htmx only; Datastar targeting is
	// response-driven — wrap the handler in wire.Handler.
	Wire *wire.Action
	// Target is the results-region selector the response swaps into
	// (htmx hx-target). Empty leaves htmx's default (the form itself).
	Target string
	// Sticky renders the top-anchored shell (sticky top-0, z-20, border,
	// solid background) — the pattern for bars above long result lists.
	Sticky bool
	// ResetText renders a ghost Reset link to the unfiltered Action URL
	// when non-empty ("Reset" via DefaultFilterBarProps).
	ResetText string
	// NoPushURL opts out of pushing the filtered URL (the zero value
	// pushes, so filtered views are shareable — the list-page convention;
	// htmx dialect only).
	NoPushURL bool
}

// DefaultFilterBarProps returns sensible defaults.
func DefaultFilterBarProps() FilterBarProps {
	return FilterBarProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		ResetText: "Reset",
	}
}

// filterBarTrigger is the composite auto-submit trigger: selects and
// checkboxes re-submit; text inputs do NOT (they change on blur — wire
// search inputs through FilterInput's debounced keyup instead). This is
// footgun #2 of docs/recipes/horizontal-filter-bar.md, encoded.
const filterBarTrigger = "change from:find select, change from:find input[type=checkbox]"

// filterBarIsHTMX reports whether the effective wiring speaks the htmx
// dialect (the default transport per ADR-0030).
func filterBarIsHTMX(props FilterBarProps) bool {
	if props.Wire == nil || props.Wire.URL == "" {
		return true
	}

	return props.Wire.Transport != wire.TransportDatastar
}

// filterBarAttributes renders the auto-submit wiring for the form element.
//
//   - Wire set: the wire action with the component's form defaults (change
//     event + form encoding), with the htmx-dialect trigger replaced by the
//     bar's composite clause. The consumer's action is never mutated
//     (value-receiver builders).
//   - Wire unset, Action set: default htmx wiring derived from Action.
//   - Neither: nil (inert form).
func filterBarAttributes(props FilterBarProps) templ.Attributes {
	if props.Wire != nil && props.Wire.URL != "" {
		attrs := props.Wire.WithFormDefaults(wire.EventChange).Attributes()
		if attrs == nil {
			return nil
		}

		if props.Wire.Transport != wire.TransportDatastar {
			attrs["hx-trigger"] = filterBarTrigger
			if !props.NoPushURL {
				attrs["hx-push-url"] = "true"
			}
		}

		return attrs
	}

	if props.Action == "" {
		return nil
	}

	attrs := templ.Attributes{
		"hx-get":     props.Action,
		"hx-trigger": filterBarTrigger,
	}

	if props.Target != "" {
		attrs["hx-target"] = props.Target
	}

	if !props.NoPushURL {
		attrs["hx-push-url"] = "true"
	}

	return attrs
}
