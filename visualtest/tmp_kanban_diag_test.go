package visualtest

import (
	"testing"

	"github.com/chromedp/chromedp"
)

func TestTmpKanbanDiag(t *testing.T) {
	server := StartDemoServer(t)

	ctx, cancel := newFlowTab(t)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(server.BaseURL()+"/kanban"), chromedp.WaitReady("body")); err != nil {
		t.Fatalf("nav: %v", err)
	}

	var diag string
	err := chromedp.Run(ctx, chromedp.Evaluate(`(()=>{
		const btn=document.querySelector('#kanban-demo-htmx [data-tc-kanban-card="kb-1"] [data-tc-kanban-move=next]');
		if(!btn)return "no button";
		const r=btn.getBoundingClientRect();
		const el=document.elementFromPoint(r.x+r.width/2, r.y+r.height/2);
		const desc=(e)=>e? e.tagName+'.'+(e.getAttribute('class')||'').slice(0,80) : 'null';
		const board=document.querySelector('#kanban-demo-htmx');
		const form=board?board.querySelector('form[data-tc-kanban-form]'):null;
		return JSON.stringify({
			rect:{x:r.x,y:r.y,w:r.width,h:r.height},
			atPoint:desc(el),
			isBtn:el===btn,
			formPresent:!!form,
			formAction:form?form.getAttribute('action')||form.getAttribute('hx-post'):'-',
			hxpostTarget:form?form.getAttribute('hx-target'):'-',
			boardCount:document.querySelectorAll('[data-tc-kanban]').length,
			progBefore:Array.from(document.querySelectorAll('#kanban-demo-htmx [data-tc-kanban-column-body="progress"] > [data-tc-kanban-card]')).map(e=>e.getAttribute('data-tc-kanban-card')).join(',')
		});
	})()`, &diag))
	if err != nil {
		t.Fatalf("eval: %v", err)
	}
	t.Logf("DIAG: %s", diag)

	if err := chromedp.Run(ctx, chromedp.Evaluate(`window.__errs=[];window.addEventListener('error',e=>window.__errs.push(String(e.message)));document.addEventListener('click',e=>{window.__lastClick=e.target.tagName+'|'+(e.target.getAttribute('data-tc-kanban-move')||'');},true);true`, nil)); err != nil {
		t.Fatalf("hook: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.Click(`#kanban-demo-htmx [data-tc-kanban-card="kb-1"] [data-tc-kanban-move=next]`, chromedp.ByQuery)); err != nil {
		t.Logf("click err: %v", err)
	}

	err = chromedp.Run(ctx, chromedp.Evaluate(`(()=>{
		const board=document.querySelector('#kanban-demo-htmx');
		const form=board?board.querySelector('form[data-tc-kanban-form]'):null;
		const fd=form?new FormData(form):null;
		const vals={};
		if(fd)fd.forEach((v,k)=>vals[k]=String(v));
		return JSON.stringify({
			errs:window.__errs||[],
			lastClick:window.__lastClick||'none',
			formVals:vals,
			scripts:document.querySelectorAll('script').length,
			noncedScripts:document.querySelectorAll('script[nonce]').length,
			kbFns:typeof window.tcKbBoard+'/'+typeof window.tcKbSubmit+'/'+typeof window.tcKbPendingFor,
			cspMeta:document.querySelector('meta[http-equiv="Content-Security-Policy"]')?'meta':'none',
			progAfter:Array.from(document.querySelectorAll('#kanban-demo-htmx [data-tc-kanban-column-body="progress"] > [data-tc-kanban-card]')).map(e=>e.getAttribute('data-tc-kanban-card')).join(','),
			backlogAfter:Array.from(document.querySelectorAll('#kanban-demo-htmx [data-tc-kanban-column-body="backlog"] > [data-tc-kanban-card]')).map(e=>e.getAttribute('data-tc-kanban-card')).join(',')
		});
	})()`, &diag))
	if err != nil {
		t.Fatalf("eval2: %v", err)
	}
	t.Logf("AFTER: %s", diag)
}
