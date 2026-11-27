package datastar_test

import (
	"github.com/larsartmann/templ-components/datastar"
)

func ExampleIndicator() {
	_ = datastar.Indicator(datastar.DefaultIndicatorProps())
	// Output:
}

func ExamplePolledRegion() {
	_ = datastar.PolledRegion(datastar.PolledRegionProps{
		ID:    "stats",
		URL:   "/partials/stats",
		Every: "10s",
	})
	// Output:
}

func ExampleLoadingButton() {
	_ = datastar.LoadingButton(datastar.LoadingButtonProps{
		Signal:      "saving",
		DefaultText: "Save",
		LoadingText: "Saving…",
	})
	// Output:
}

func ExampleLiveRegion() {
	_ = datastar.LiveRegion(datastar.DefaultLiveRegionProps())
	// Output:
}

func ExampleSDKScript() {
	_ = datastar.SDKScript(datastar.DefaultSDKScriptProps())
	// Output:
}

func ExampleSSEErrorHandling() {
	_ = datastar.SSEErrorHandling(datastar.DefaultSSEErrorHandlingConfig())
	// Output:
}

func ExampleSignalChip() {
	_ = datastar.SignalChip(datastar.SignalChipProps{
		Signal:   "stateFilter",
		Value:    "open",
		Label:    "Open",
		Selected: true,
	})
	// Output:
}

func ExampleSignalChipGroup() {
	_ = datastar.SignalChipGroup(datastar.SignalChipGroupProps{
		AriaLabel: "Filter by state",
		Signal:    "stateFilter",
		Selected:  "open",
		Items: []datastar.SignalChipItem{
			{Label: "All"},
			{Label: "Open", Value: "open"},
			{Label: "Done", Value: "done"},
		},
	})
	// Output:
}
