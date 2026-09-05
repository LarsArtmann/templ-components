package datastar

import "github.com/larsartmann/templ-components/utils"

// SDKScriptProps configures injection of the Datastar runtime <script> tag.
// Include this once per page in your layout, alongside or instead of the HTMX
// script. Datastar loads as an ES module (deferred by default).
//
//	@datastar.SDKScript(datastar.DefaultSDKScriptProps())
type SDKScriptProps struct {
	utils.BaseProps

	// Version pins the Datastar runtime version. Defaults to DatastarVersion1_0_3.
	// Set to "" to use the default.
	Version DatastarVersion

	// CDN overrides the CDN base URL. Defaults to the jsDelivr GitHub CDN.
	// Consumers can point this at their own origin or pnpm mirror.
	CDN string

	// Src overrides the full script URL — use for self-hosting.
	// When set, Version and CDN are ignored.
	Src string
}

// DefaultSDKScriptProps returns sensible defaults for the Datastar runtime script.
func DefaultSDKScriptProps() SDKScriptProps {
	return SDKScriptProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Version: defaultDatastarVersion,
	}
}

// scriptSrc resolves the final src URL for the Datastar runtime script.
// When Src is set, it takes precedence (self-hosting). Otherwise the CDN
// URL is built from Version and CDN.
func scriptSrc(props SDKScriptProps) string {
	if props.Src != "" {
		return props.Src
	}

	return datastarScriptURL(props.Version, props.CDN)
}
