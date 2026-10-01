package htmx_test

import (
	"github.com/larsartmann/templ-components/htmx"
)

func ExampleLoadingIndicator() {
	_ = htmx.LoadingIndicator(nil)
	// Output:
}

func ExampleInlineLoadingOverlay() {
	_ = htmx.InlineLoadingOverlay("results", nil)
	// Output:
}

func ExampleLoadingButton() {
	_ = htmx.LoadingButton("Save", "Saving...", nil)
	// Output:
}

func ExampleConfirmDelete() {
	_ = htmx.ConfirmDelete(htmx.ConfirmDeleteProps{Delete: "/items/1", Target: "#list"})
	// Output:
}

func ExampleSwapOOB() {
	_ = htmx.SwapOOB(htmx.SwapOOBProps{Selector: "#list"})
	// Output:
}

func ExampleCSRFToken() {
	_ = htmx.CSRFToken("csrf-token")
	// Output:
}

func ExampleGlobalErrorHandling() {
	_ = htmx.GlobalErrorHandling(htmx.DefaultErrorHandlingConfig())
	// Output:
}

func ExamplePolledRegion() {
	_ = htmx.PolledRegion(htmx.DefaultPolledRegionProps())
	// Output:
}

func ExampleViewTransitions() {
	_ = htmx.ViewTransitions(htmx.DefaultViewTransitionsProps())
	// Output:
}
