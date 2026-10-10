// FilterChips component: zero-JS link-based filter chips — filter state
// lives in the URL, chips are plain links, no JavaScript anywhere (SEO- and
// crawler-friendly). Extracted from nsfw-classifier's historyChip
// (history_page.templ:105, 2026-10-10 analysis, TODO_LIST #391).
package forms

import (
	"net/url"

	"github.com/larsartmann/templ-components/utils"
)

// FilterChip is one filter link. The consumer owns the href — filter state
// is query-param truth, and only the consumer knows the page's full param
// set (cursors, lanes, sibling filters). Build common toggles with
// FilterToggleHref.
type FilterChip struct {
	// Label is the visible chip text.
	Label string
	// Href is the link target carrying the chip's filter state.
	Href string
	// Active marks the currently applied filter. The active chip renders
	// aria-current="true" and the filled style; keeping the link (like the
	// extraction source) is honest about the current state.
	Active bool
}

// FilterChipsProps configures a row of zero-JS filter chips.
type FilterChipsProps struct {
	utils.BaseProps

	// Chips renders left to right. Provide an "All" chip (empty-value href)
	// so users can clear the filter without JavaScript.
	Chips []FilterChip
}

// DefaultFilterChipsProps returns sensible defaults.
func DefaultFilterChipsProps() FilterChipsProps {
	return FilterChipsProps{} //nolint:exhaustruct // no non-zero defaults
}

// FilterToggleHref builds a chip href that sets key=value on basePath,
// preserving every other query param; an empty value REMOVES the key (the
// "All" chip). Values are URL-escaped via url.Values — injection-safe by
// construction.
//
//	href := forms.FilterToggleHref("/history", r.URL.Query(), "source", "upload")
func FilterToggleHref(basePath string, params url.Values, key, value string) string {
	values := make(url.Values, len(params))
	for k, v := range params {
		values[k] = v
	}

	if value == "" {
		values.Del(key)
	} else {
		values.Set(key, value)
	}

	encoded := values.Encode()
	if encoded == "" {
		return basePath
	}

	return basePath + "?" + encoded
}

// filterChipClass resolves the chip style: filled when active, outlined
// otherwise.
func filterChipClass(active bool) string {
	if active {
		return "border-blue-600 bg-blue-600 text-white hover:bg-blue-700 dark:border-blue-500 dark:bg-blue-500 dark:text-white dark:hover:bg-blue-400"
	}

	return "border-gray-300 bg-white text-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
}

// filterChipClasses joins the shared chip shape with the state style.
func filterChipClasses(active bool, extra string) string {
	return utils.Class(
		"rounded-full border px-3 py-1 text-sm transition-colors motion-reduce:transition-none motion-reduce:duration-0",
		filterChipClass(active),
		extra,
	)
}

// filterChipsLabel is the accessible name for a chip row whose consumer set
// neither AriaLabel nor a visible heading — chip rows are filter groups and
// must be named.
const filterChipsLabel = "Filters"
