package echarts

import "github.com/larsartmann/templ-components/utils"

// EChartsVersion pins the Apache ECharts runtime version loaded from CDN.
type EChartsVersion string

const (
	// EChartsVersion5_5 is ECharts v5.5.0 (stable, Baseline 2024).
	EChartsVersion5_5 EChartsVersion = "5.5.0"
)

// DefaultEChartsVersion is the version used when SDKScriptProps.Version is empty.
const DefaultEChartsVersion EChartsVersion = EChartsVersion5_5

// DefaultCDNHost is the default CDN for loading the ECharts runtime.
const DefaultCDNHost = "https://cdn.jsdelivr.net/npm"

// EChartsProps configures the EChart wrapper component. The consumer builds
// their chart with go-echarts, calls RenderSnippet(), and passes the Element and
// Script strings here. This package does NOT import go-echarts — it just wraps
// the output in a CSP-safe container with optional dark mode bridging.
type EChartsProps struct {
	utils.BaseProps

	// Element is the HTML snippet from go-echarts RenderSnippet().Element —
	// typically a <div> container with a unique ID.
	Element string

	// Script is the JS from go-echarts RenderSnippet().Script — typically
	// echarts.init(...).setOption(...). Injected as a CSP-safe inline script.
	Script string

	// Nonce is the CSP nonce for inline scripts. Required for CSP compliance.
	Nonce string

	// DarkModeBridge injects the singleton dark mode bridge script that syncs
	// the ECharts theme with the Tailwind .dark class. Default: true.
	DarkModeBridge bool
}

// DefaultEChartsProps returns sensible defaults for the EChart component.
func DefaultEChartsProps() EChartsProps {
	return EChartsProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		DarkModeBridge: true,
	}
}

// SDKScriptProps configures injection of the ECharts runtime <script> tag.
// Include this once per page in your layout, similar to the HTMX or Datastar
// script.
type SDKScriptProps struct {
	utils.BaseProps

	// Version pins the ECharts runtime version. Defaults to EChartsVersion5_5.
	Version EChartsVersion

	// CDN overrides the CDN base URL. Defaults to DefaultCDNHost.
	// Consumers can point this at their own origin or pnpm mirror.
	CDN string

	// Src overrides the full script URL — use for self-hosting.
	// When set, Version and CDN are ignored.
	Src string

	// Nonce is the CSP nonce for the external script tag.
	Nonce string
}

// DefaultSDKScriptProps returns sensible defaults for the ECharts runtime script.
func DefaultSDKScriptProps() SDKScriptProps {
	return SDKScriptProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Version: DefaultEChartsVersion,
	}
}

// scriptSrc resolves the final src URL for the ECharts runtime script.
// When Src is set, it takes precedence (self-hosting). Otherwise the CDN
// URL is built from Version and CDN.
func scriptSrc(props SDKScriptProps) string {
	if props.Src != "" {
		return props.Src
	}

	return echartsScriptURL(props.Version, props.CDN)
}

// echartsScriptURL builds the CDN URL for the ECharts runtime.
func echartsScriptURL(version EChartsVersion, cdn string) string {
	if version == "" {
		version = DefaultEChartsVersion
	}

	if cdn == "" {
		cdn = DefaultCDNHost
	}

	return cdn + "/echarts@" + string(version) + "/dist/echarts.min.js"
}
