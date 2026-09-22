package visualtest

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func TestDebugClipboard(t *testing.T) {
	base := requireSiteDist(t)

	tabCtx, tabCancel := newTab(t)
	defer tabCancel()

	ctx, cancel := context.WithTimeout(tabCtx, 30*time.Second)
	defer cancel()

	var writeResult, permState string

	if err := chromedp.Run(ctx,
		chromedp.Navigate(base+"/sales"),
		chromedp.WaitVisible(`[data-tc-copy]`, chromedp.ByQuery),
		browser.SetPermission(
			&browser.PermissionDescriptor{Name: "clipboard-write"},
			browser.PermissionSettingGranted,
		).WithOrigin(base),
		emulation.SetFocusEmulationEnabled(true),
		chromedp.Evaluate(`window.__p=null; navigator.permissions.query({name:'clipboard-write'}).then(s=>{window.__p=s.state}).catch(e=>{window.__p='err '+e.name}); 'k'`, nil),
		chromedp.Sleep(300*time.Millisecond),
		chromedp.Evaluate(`window.__p`, &permState),
		chromedp.Evaluate(`window.__w = null; navigator.clipboard.writeText('x').then(() => { window.__w = 'resolved'; }).catch(e => { window.__w = 'rejected: '+e.name+' '+e.message; }); 'kicked'`, nil),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.Evaluate(`window.__w`, &writeResult),
	); err != nil {
		t.Fatalf("debug: %v", err)
	}

	t.Logf("perm state: %s", permState)
	t.Logf("writeText: %s", writeResult)
}
