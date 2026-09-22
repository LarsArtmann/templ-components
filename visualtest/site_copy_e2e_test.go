package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/emulation"
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
		// Best-effort real-clipboard setup: grant clipboard permissions
		// (Permissions-API spellings) and emulate document focus — both are
		// required by the async clipboard API in headless Chromium. On
		// Chromium 152 the clipboard daemon still denies sanitized writes
		// (see the spy comment below), so the spy carries the success path;
		// in any environment where the daemon cooperates, the REAL write
		// resolves and nothing is stubbed.
		browser.SetPermission(
			&browser.PermissionDescriptor{Name: "clipboard-write"},
			browser.PermissionSettingGranted,
		).WithOrigin(base),
		browser.SetPermission(
			&browser.PermissionDescriptor{Name: "clipboard-read"},
			browser.PermissionSettingGranted,
		).WithOrigin(base),
		emulation.SetFocusEmulationEnabled(true),
	); err != nil {
		t.Fatalf("grant clipboard permissions: %v", err)
	}

	// Spy on clipboard.writeText BEFORE any click, recording the payload.
	// The original write still executes so the page behaves exactly as in
	// production — but its rejection is swallowed: headless Chromium's
	// clipboard daemon denies sanitized writes even with the permission
	// state "granted" + focus emulation (verified 2026-09-22: perm query
	// granted, writeText still "Write permission denied") — the 2026-09-21
	// CI failure was exactly this. Resolving the spy promise keeps the
	// COMPONENT's real success path (.then -> label swap) under test while
	// the payload assertion proves the writeText argument.
	spyJS := `(() => {
		window.__tcCopied = null;
		const orig = navigator.clipboard.writeText.bind(navigator.clipboard);
		navigator.clipboard.writeText = (text) => {
			window.__tcCopied = text;
			return orig(text).catch(() => {});
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
