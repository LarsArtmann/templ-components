package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestPolledRegionA11y is the accessibility lens for the polled region.
func TestPolledRegionA11y(t *testing.T) {
	t.Parallel()

	t.Run("region announces updates politely by default", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/partials/stats"}))
		utils.AssertContains(t, output, `aria-live="polite"`)
	})

	t.Run("assertive politeness is consumer-selectable", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x", Live: LiveAssertive}))
		utils.AssertContains(t, output, `aria-live="assertive"`)
	})

	t.Run("off disables announcements", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x", Live: LiveOff}))
		utils.AssertContains(t, output, `aria-live="off"`)
	})

	t.Run("no aria-busy is rendered", func(t *testing.T) {
		t.Parallel()

		// Polling is background refresh, not a blocking load: the initial
		// content is server-rendered in the children slot, so there is no
		// "waiting" window to announce and aria-busy would never clear
		// meaningfully (it re-renders in every patch).
		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x"}))
		utils.AssertNotContains(t, output, "aria-busy")
	})

	t.Run("timestamp footer keeps low visual emphasis", func(t *testing.T) {
		t.Parallel()

		output := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x", ShowTimestamp: true}))
		utils.AssertContains(t, output, "text-xs")
		utils.AssertContains(t, output, "<time")
		utils.AssertContains(t, output, "datetime=")
	})

	t.Run("labelled region is a region landmark", func(t *testing.T) {
		t.Parallel()

		// aria-label on a roleless generic div is invalid ARIA (the class the
		// HTML-validation gate fixed for KanbanBoard, Scrollback, Carousel):
		// a labelled polled region must expose role="region".
		labelled := utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x", AriaLabel: "Activity"}))
		utils.AssertContains(t, labelled, `role="region"`)
		utils.AssertContains(t, labelled, `aria-label="Activity"`)
		utils.AssertNotContains(t, utils.Render(t, PolledRegion(PolledRegionProps{URL: "/x"})), `role="region"`)
	})
}
