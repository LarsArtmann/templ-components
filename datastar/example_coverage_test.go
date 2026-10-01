package datastar_test

import (
	"github.com/larsartmann/templ-components/datastar"
)

func ExampleIndicator() {
	_ = datastar.Indicator(datastar.DefaultIndicatorProps())
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
