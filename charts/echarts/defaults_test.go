package echarts

import "testing"

// TestDefaultProps pins the documented defaults: the dark-mode bridge is on
// (consumers get theme-aware charts without opting in) and the SDK script
// falls back to the pinned ECharts version.
func TestDefaultProps(t *testing.T) {
	t.Parallel()

	if !DefaultEChartsProps().DarkModeBridge {
		t.Error("DefaultEChartsProps().DarkModeBridge = false, want true")
	}

	props := DefaultSDKScriptProps()
	if props.Version != DefaultEChartsVersion {
		t.Errorf("DefaultSDKScriptProps().Version = %q, want %q", props.Version, DefaultEChartsVersion)
	}

	if props.Src != "" || props.CDN != "" {
		t.Errorf("DefaultSDKScriptProps() should leave Src/CDN empty, got Src=%q CDN=%q", props.Src, props.CDN)
	}
}
