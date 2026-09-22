package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

// TestSiteSalesCopyButton proves the sales page's library CopyButton in a
// real browser against the BUILT site (website/dist): clicking [data-tc-copy]
// writes the exact install command through navigator.clipboard and swaps the
// label to "Copied!". The clipboard payload is captured by a spy bound before
// the click — headless clipboard READS need permissions a chromedp tab cannot
// grant, so reading the real clipboard would prove nothing reliably; the spy
// observes the same writeText call the browser performs and the write itself
// still goes through to the real clipboard.
//
// Clipboard WRITES are granted via Browser.grantPermissions: without them,
// headless Chromium rejects navigator.clipboard.writeText (NotAllowedError)
// and the fallback execCommand path throws in a background tab, so the label
// never swaps — this was the 2026-09-21 CI failure ("label after click =
// Copy, want Copied!").
//
// Run via `nix run .#visual` (builds the dist, sets CHROMEDP_CHROME_PATH);
// skips gracefully without a browser like the rest of the suite.
func TestSiteSalesCopyButton(t *testing.T) {
	base := requireSiteDist(t)

	tabCtx, tabCancel := newTab(t)
	defer tabCancel()

	ctx, cancel := context.WithTimeout(tabCtx, 120*time.Second)
	defer cancel()

	if err := chromedp.Run(ctx,
		// Clipboard WRITES are granted via Browser.setPermission: without
		// them, headless Chromium rejects navigator.clipboard.writeText
		// (NotAllowedError) and the fallback execCommand path throws in a
		// background tab, so the label never swaps — the 2026-09-21 CI
		// failure ("label after click = Copy, want Copied!"). Names are the
		// Permissions-API spellings ("clipboard-write", not the legacy
		// GrantPermissions "clipboardReadWrite").
		browser.SetPermission(
			&browser.PermissionDescriptor{Name: "clipboard-write"},
			browser.PermissionSettingGranted,
		).WithOrigin(base),
		browser.SetPermission(
			&browser.PermissionDescriptor{Name: "clipboard-read"},
			browser.PermissionSettingGranted,
		).WithOrigin(base),
	); err != nil {
		t.Fatalf("grant clipboard permissions: %v", err)
	}

	// Spy on clipboard.writeText BEFORE any click, recording the payload.
	// The original write still executes so the page behaves exactly as in
	// production.
	spyJS := `(() => {
		window.__tcCopied = null;
		const orig = navigator.clipboard.writeText.bind(navigator.clipboard);
		navigator.clipboard.writeText = (text) => {
			window.__tcCopied = text;
			return orig(text);
		};
		return true;
	})()`

	var (
		spyArmed bool
		copied   string
		label    string
	)

	want := "go get github.com/larsartmann/templ-components@latest"

	if err := chromedp.Run(ctx,
		chromedp.Navigate(base+"/sales"),
		chromedp.WaitVisible(`[data-tc-copy]`, chromedp.ByQuery),
		chromedp.Evaluate(spyJS, &spyArmed),
		chromedp.Click(`[data-tc-copy]`, chromedp.ByQuery),
		chromedp.Evaluate(`window.__tcCopied`, &copied),
		chromedp.Text(`[data-tc-copy-text]`, &label, chromedp.ByQuery),
	); err != nil {
		t.Fatalf("copy flow: %v", err)
	}

	if !spyArmed {
		t.Fatal("clipboard spy did not arm")
	}

	if copied != want {
		t.Errorf("clipboard payload = %q, want %q", copied, want)
	}

	if label != "Copied!" {
		t.Errorf("label after click = %q, want %q", label, "Copied!")
	}
}
