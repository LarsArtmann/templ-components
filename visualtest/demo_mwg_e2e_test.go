package visualtest

import (
	"testing"

	"github.com/chromedp/chromedp"
)

// E2E for the Modern Web Guidance compliance wave's two new client-side
// behaviors, proving them in a real browser rather than string assertions:
// the Form validity bridge (aria-invalid mirroring :user-invalid) and the
// Slider track-fill sync (--tc-slider-value via the CSSOM). Both surfaces
// live on the /forms demo page: the Create Account card (inside forms.Form,
// so the bridge script is live) and the Volume slider in the control
// showcase.

// TestFormValidityBridgeMirrorsUserInvalid drives the required email input
// from valid to invalid and back, asserting aria-invalid follows
// :user-invalid at each step — the exact sync the bridge exists for.
func TestFormValidityBridgeMirrorsUserInvalid(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	emailSel := `#forms-create-account input[name="email"]`

	if err := chromedp.Do(
		ctx,
		chromedp.Navigate(server.BaseURL()+"/forms"),
		chromedp.WaitReady(emailSel),
	); err != nil {
		t.Fatalf("visualtest[demo]: load forms page: %v", err)
	}

	// Server-rendered state: no aria-invalid on the email field.
	var aria string
	if err := chromedp.Do(
		ctx,
		evalExprInto[string](`(document.querySelector('`+emailSel+`').getAttribute('aria-invalid') || '')`, &aria),
	); err != nil {
		t.Fatalf("visualtest[demo]: read initial aria-invalid: %v", err)
	}

	if aria != "" {
		t.Fatalf("visualtest[demo]: email field started with aria-invalid=%q; want unset", aria)
	}

	// Type an invalid email, then blur: the bridge must set aria-invalid.
	if err := chromedp.Do(
		ctx,
		chromedp.SendKeys(emailSel, "not-an-email"),
		chromedp.Evaluate[chromedp.Void](`document.querySelector('`+emailSel+`').blur()`),
	); err != nil {
		t.Fatalf("visualtest[demo]: type invalid email + blur: %v", err)
	}

	if err := chromedp.Do(
		ctx,
		pollText(
			`(document.querySelector('`+emailSel+`').matches(':user-invalid') ? 'user-invalid' : 'clean') + '/' + (document.querySelector('`+emailSel+`').getAttribute('aria-invalid') || 'none')`,
			&aria,
			chromedp.WithPollingTimeout(demoFlowTimeout),
		),
	); err != nil {
		t.Fatalf("visualtest[demo]: validity bridge never mirrored :user-invalid: %v", err)
	}

	if aria != "user-invalid/true" {
		t.Fatalf("visualtest[demo]: after invalid blur got %q; want user-invalid/true", aria)
	}

	// Correct the value and blur again: the bridge must REMOVE its own state.
	if err := chromedp.Do(
		ctx,
		chromedp.Evaluate[chromedp.Void](
			`(function (el) { el.value = ''; })(document.querySelector('`+emailSel+`'))`,
		),
		chromedp.SendKeys(emailSel, "user@example.com"),
		chromedp.Evaluate[chromedp.Void](`document.querySelector('`+emailSel+`').blur()`),
	); err != nil {
		t.Fatalf("visualtest[demo]: fix email + blur: %v", err)
	}

	if err := chromedp.Do(
		ctx,
		pollText(
			`document.querySelector('`+emailSel+`').getAttribute('aria-invalid') || 'none'`,
			&aria,
			chromedp.WithPollingTimeout(demoFlowTimeout),
		),
	); err != nil {
		t.Fatalf("visualtest[demo]: validity bridge never cleared aria-invalid: %v", err)
	}

	if aria != "none" {
		t.Fatalf("visualtest[demo]: after valid blur aria-invalid=%q; want removed", aria)
	}
}

// TestFormValidityBridgePreservesServerErrorState proves the ownership
// contract: the username field ships with a SERVER-set aria-invalid (from
// forms.ErrorAttrs). Interacting with it validly must not strip that state —
// the bridge only manages fields it marked itself (data-tc-user-invalid).
func TestFormValidityBridgePreservesServerErrorState(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	userSel := `#forms-create-account input[name="username"]`

	if err := chromedp.Do(
		ctx,
		chromedp.Navigate(server.BaseURL()+"/forms"),
		chromedp.WaitReady(userSel),
		chromedp.SendKeys(userSel, "validuser"),
		chromedp.Evaluate[chromedp.Void](`document.querySelector('`+userSel+`').blur()`),
	); err != nil {
		t.Fatalf("visualtest[demo]: interact with server-invalid username: %v", err)
	}

	// Give the bridge its shot, then assert the server state survived.
	var aria string
	if err := chromedp.Do(
		ctx,
		evalExprInto[string](`document.querySelector('`+userSel+`').getAttribute('aria-invalid') || 'none'`, &aria),
	); err != nil {
		t.Fatalf("visualtest[demo]: read username aria-invalid: %v", err)
	}

	if aria != "true" {
		t.Fatalf(
			"visualtest[demo]: server-set aria-invalid was %q after valid interaction; want preserved \"true\"",
			aria,
		)
	}
}

// TestSliderSyncMirrorsDragIntoCSSOM drives the volume slider and asserts
// the sync singleton copies the live value into the --tc-slider-value
// custom property via the CSSOM (element.style), which strict style-src
// CSPs permit — the property the WebKit track-fill gradient consumes.
func TestSliderSyncMirrorsDragIntoCSSOM(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	// NB: braceless demoSection calls leave section anchors empty with the
	// content as a following sibling (pre-existing demo convention), so the
	// selector must not require the input to be INSIDE #forms-controls.
	sliderSel := `input[name="volume"]`

	if err := chromedp.Do(
		ctx,
		chromedp.Navigate(server.BaseURL()+"/forms"),
		chromedp.WaitReady(sliderSel),
	); err != nil {
		t.Fatalf("visualtest[demo]: load forms page: %v", err)
	}

	var startProp string
	if err := chromedp.Do(
		ctx,
		evalExprInto[string](
			`document.querySelector('`+sliderSel+`').style.getPropertyValue('--tc-slider-value')`,
			&startProp,
		),
	); err != nil {
		t.Fatalf("visualtest[demo]: read initial --tc-slider-value: %v", err)
	}

	if startProp != "" {
		t.Fatalf("visualtest[demo]: slider shipped with --tc-slider-value=%q before any drag; want unset", startProp)
	}

	// Drag (programmatically): set the value PROPERTY and dispatch input —
	// exactly what a real drag does (the attribute stays server-rendered).
	if err := chromedp.Do(
		ctx,
		chromedp.Evaluate[chromedp.Void](
			`(function (el) { el.value = '75'; el.dispatchEvent(new Event('input', {bubbles: true})); })(document.querySelector('`+sliderSel+`'))`,
		),
	); err != nil {
		t.Fatalf("visualtest[demo]: dispatch slider input: %v", err)
	}

	var synced string
	if err := chromedp.Do(
		ctx,
		pollText(
			`document.querySelector('`+sliderSel+`').style.getPropertyValue('--tc-slider-value')`,
			&synced,
			chromedp.WithPollingTimeout(demoFlowTimeout),
		),
	); err != nil {
		t.Fatalf("visualtest[demo]: slider sync never mirrored the value: %v", err)
	}

	if synced != "75" {
		t.Fatalf("visualtest[demo]: --tc-slider-value=%q after drag; want 75", synced)
	}
}
