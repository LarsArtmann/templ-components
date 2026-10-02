package datastar_test

import (
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/utils"
)

func ExampleIndicator() {
	_ = datastar.Indicator(datastar.DefaultIndicatorProps())
	// Output:
}

func ExamplePolledRegion() {
	_ = datastar.PolledRegion(datastar.PolledRegionProps{
		BaseProps: utils.BaseProps{ID: "stats"},
		URL:       "/partials/stats",
		Every:     "10s",
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
