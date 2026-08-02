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
		`src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.2/bundles/datastar.js"`)
}

func TestSDKScriptSelfHosted(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SDKScript(SDKScriptProps{
		Src: "/static/datastar.js",
	}))

	utils.AssertContains(t, output, `src="/static/datastar.js"`)
	utils.AssertNotContains(t, output, "cdn.jsdelivr.net")
}

func TestSDKScriptCustomCDN(t *testing.T) {
	t.Parallel()

	output := utils.Render(t, SDKScript(SDKScriptProps{
		CDN:     "https://unpkg.com",
		Version: DatastarVersion1_0_2,
	}))

	utils.AssertContains(t, output, `src="https://unpkg.com/starfederation/datastar@1.0.2/bundles/datastar.js"`)
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

	if props.Version != DatastarVersion1_0_2 {
		t.Errorf("expected Version=%s, got %s", DatastarVersion1_0_2, props.Version)
	}
}
