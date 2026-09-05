package datastar

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestSDKScriptDefault(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SDKScript(DefaultSDKScriptProps()))

	utils.AssertContains(t, output, `<script type="module"`)
	utils.AssertContains(t, output,
		`src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.3/bundles/datastar.js"`)
}

func TestSDKScriptSelfHosted(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SDKScript(SDKScriptProps{
		Src: "/static/datastar.js",
	}))

	utils.AssertContains(t, output, `src="/static/datastar.js"`)
	utils.AssertNotContains(t, output, "cdn.jsdelivr.net")
}

// Self-hosting must not emit the CDN preconnect: there is no cross-origin
// connection to warm, and a dangling preconnect to jsdelivr would leak the
// page load to a third party even on fully self-hosted deployments.
func TestSDKScriptSelfHostedOmitsPreconnect(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SDKScript(SDKScriptProps{
		Src: "/static/datastar.js",
	}))

	utils.AssertNotContains(t, output, "preconnect")
}

func TestSDKScriptPreconnect(t *testing.T) {
	t.Parallel()

	defaultCDN := utils.Render(t, SDKScript(DefaultSDKScriptProps()))
	utils.AssertContains(t, defaultCDN, `<link rel="preconnect" href="https://cdn.jsdelivr.net" crossorigin>`)

	customCDN := utils.Render(t, SDKScript(SDKScriptProps{
		CDN: "https://unpkg.com",
	}))
	utils.AssertContains(t, customCDN, `<link rel="preconnect" href="https://unpkg.com" crossorigin>`)
}

// CDN + nonce must combine: preconnect link first, then the nonced module
// script. Guarded by sdk_script_cdn_nonce.golden for the full document shape.
func TestSDKScriptCDNWithNonce(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SDKScript(SDKScriptProps{
		BaseProps: utils.BaseProps{Nonce: "test-nonce-123"},
		Version:   DatastarVersion1_0_3,
	}))

	utils.AssertContains(t, output, `<link rel="preconnect" href="https://cdn.jsdelivr.net" crossorigin>`)
	utils.AssertContains(t, output, `<script type="module" nonce="test-nonce-123"`)
	utils.AssertContains(
		t,
		output,
		`src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.3/bundles/datastar.js"`,
	)
}

func TestSDKScriptCustomCDN(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SDKScript(SDKScriptProps{
		CDN:     "https://unpkg.com",
		Version: DatastarVersion1_0_3,
	}))

	utils.AssertContains(t, output, `src="https://unpkg.com/starfederation/datastar@1.0.3/bundles/datastar.js"`)
}

func TestSDKScriptNonce(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SDKScript(SDKScriptProps{
		BaseProps: utils.BaseProps{Nonce: "abc123"},
	}))

	utils.AssertContains(t, output, `nonce="abc123"`)
}

func TestSDKScriptDefaultProps(t *testing.T) {
	t.Parallel()

	props := DefaultSDKScriptProps()

	if props.Version != DatastarVersion1_0_3 {
		t.Errorf("expected Version=%s, got %s", DatastarVersion1_0_3, props.Version)
	}
}
