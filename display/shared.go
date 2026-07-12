package display

import (
	"context"
	"fmt"
	"html"
	"io"
	"strconv"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// mutedTextClass is the shared Tailwind class string for secondary/muted text
// (Card subtitle, PageHeader subtitle, EmptyState description). Callers combine
// it with a margin utility via utils.Class, e.g. utils.Class(mutedTextClass, "mt-1").
const mutedTextClass = "text-sm text-gray-500 dark:text-gray-400"

// Motion class constants — re-exported from utils for backward compatibility.
// Use utils.TransitionFast, utils.TransitionNormal, etc. directly in new code.
const (
	transitionFast      = utils.TransitionFast
	transitionNormal    = utils.TransitionNormal
	transitionColors    = utils.TransitionColors
	transitionTransform = utils.TransitionTransform
)

// OverlayKind identifies whether an overlay is a Modal or a Drawer.
type OverlayKind string

const (
	// OverlayModal is the kind for Modal overlays.
	OverlayModal OverlayKind = "modal"
	// OverlayDrawer is the kind for Drawer overlays.
	OverlayDrawer OverlayKind = "drawer"
)

// componentName returns the capitalized name used in generated JS function names
// (tcOpenModal, tcCloseDrawer, etc.). Unknown kinds fall back to "Overlay".
func (k OverlayKind) componentName() string {
	switch k {
	case OverlayModal:
		return "Modal"
	case OverlayDrawer:
		return "Drawer"
	default:
		return "Overlay"
	}
}

// overlayShellProps parameterizes the shared <dialog>-based shell used by
// Modal and Drawer. The native <dialog> element provides focus trapping,
// Escape-to-close, focus restore, top-layer rendering, and ::backdrop —
// so the shell is minimal: just the dialog element, its visual classes,
// and a tiny script for auto-open and click delegation.
type overlayShellProps struct {
	id          string
	open        bool
	title       string
	ariaLabel   string
	kind        OverlayKind
	nonce       string
	dialogClass string           // Tailwind classes for the <dialog> element (visual panel)
	side        DrawerSide       // only for Drawer; zero value for Modal
	attrs       templ.Attributes // consumer extra HTML attributes from BaseProps.Attrs
}

// overlayDialogJS generates the minimal JavaScript for a <dialog>-based overlay.
// The native <dialog> API provides focus trapping, Escape-to-close, focus
// restore, and top-layer rendering — so the JS is just:
//   - tcOpen{Name}(id) / tcClose{Name}(id): thin wrappers for backward compat
//   - Per-instance IIFE: auto-open if data-tc-open is set, click delegation
//     for [data-tc-close] buttons and backdrop-click-to-close
//
// Backdrop click detection: when the user clicks the ::backdrop area, the
// click event target IS the dialog element itself (the backdrop is a
// pseudo-element of dialog). So e.target === dialog detects backdrop clicks.
func overlayDialogJS(id, componentName string) string {
	guard := "window.tcOverlay" + componentName + "Attached"
	escapedID := strconv.Quote(id)
	return "if(!" + guard + "){" + guard + "=true;" +
		"window.tcOpen" + componentName + "=function(id){" +
		"var d=document.getElementById(id);" +
		"if(d&&typeof d.showModal==='function')d.showModal();" +
		"};" +
		"window.tcClose" + componentName + "=function(id){" +
		"var d=document.getElementById(id);if(d)d.close();" +
		"};" +
		"}" +
		"(function(id){" +
		"var d=document.getElementById(id);if(!d)return;" +
		"if(d.hasAttribute('data-tc-open')&&!d.open)d.showModal();" +
		"d.addEventListener('click',function(e){" +
		"if(e.target.closest('[data-tc-close]')||e.target===d)d.close();" +
		"});" +
		"})(" + escapedID + ");\n"
}

// overlayScriptComponent renders a <script nonce="..."> tag containing the
// overlay JS. It writes directly to the output buffer to bypass templ's
// script-context sanitization (which would JSON-encode the JS string).
// The nonce is HTML-attribute-escaped to prevent attribute-boundary breakage
// from a caller-supplied value containing quotes or angle brackets.
// scriptComponent renders a CSP-safe <script nonce="..."> tag wrapping the
// given JS string. Shared by all singleton-script components to avoid duplicating
// the nonce-escaping and error-wrapping pattern. The errLabel is used in the
// wrapped error message for debugging.
func scriptComponent(nonce, js, errLabel string) templ.Component {
	escapedNonce := html.EscapeString(nonce)
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, "<script nonce=\"%s\">\n%s</script>\n", escapedNonce, js); err != nil {
			return fmt.Errorf("write %s: %w", errLabel, err)
		}
		return nil
	})
}

// overlayScriptComponent renders the overlay (modal/drawer) JS in a CSP-safe
// <script nonce> tag.
func overlayScriptComponent(nonce, id, componentName string) templ.Component {
	return scriptComponent(nonce, overlayDialogJS(id, componentName), "overlay script")
}

// tooltipJS returns the singleton JavaScript for tooltip touch support and
// Escape-to-dismiss. Uses a window.tcTooltipAttached guard so it executes
// only once per page regardless of how many Tooltip components are rendered.
func tooltipJS() string {
	return `if(!window.tcTooltipAttached){window.tcTooltipAttached=true;` +
		`function tcPropagateTooltipDesc(){` +
		`document.querySelectorAll('[data-tc-tooltip]').forEach(function(w){` +
		`var t=w.querySelector('a[href],button:not([disabled]),input,select,textarea,[tabindex]:not([tabindex="-1"])');` +
		`if(t&&!t.getAttribute('aria-describedby')){var d=w.getAttribute('aria-describedby');if(d)t.setAttribute('aria-describedby',d);}` +
		`});}` +
		`tcPropagateTooltipDesc();` +
		`document.body.addEventListener('htmx:afterSettle',tcPropagateTooltipDesc);` +
		`document.addEventListener('click',function(e){` +
		`var trigger=e.target.closest('[data-tc-tooltip]');` +
		`if(trigger){e.preventDefault();var t=trigger.querySelector('[role="tooltip"]');if(t)t.classList.toggle('hidden');}` +
		`else{document.querySelectorAll('[data-tc-tooltip] [role="tooltip"]:not(.hidden)').forEach(function(t){t.classList.add('hidden');});}` +
		`});` +
		`document.addEventListener('keydown',function(e){` +
		`if(e.key==='Escape'){document.querySelectorAll('[data-tc-tooltip] [role="tooltip"]:not(.hidden)').forEach(function(t){t.classList.add('hidden');});}` +
		`});}` +
		"\n"
}

// tooltipScriptComponent renders a <script nonce="..."> tag containing the
// tooltip touch/Escape JS. The script is a singleton — only the first
// Tooltip on the page injects executable code; subsequent instances
// skip it via the window.tcTooltipAttached guard.
func tooltipScriptComponent(nonce string) templ.Component {
	return scriptComponent(nonce, tooltipJS(), "tooltip script")
}

// copyButtonJS returns the singleton JavaScript for clipboard copy via event
// delegation. Listens for clicks on [data-tc-copy] buttons, copies the text
// via navigator.clipboard.writeText (with a document.execCommand fallback for
// non-secure HTTP contexts or older browsers), and temporarily swaps the label.
func copyButtonJS() string {
	return `if(!window.tcCopyAttached){window.tcCopyAttached=true;` +
		`function tcSwapLabel(el,label,original){` +
		`if(el)el.textContent=label;` +
		`setTimeout(function(){if(el)el.textContent=original;},2000);}` +
		`function tcFallbackCopy(text,done){` +
		`var ta=document.createElement('textarea');ta.value=text;` +
		`ta.style.cssText='position:fixed;top:-9999px;left:-9999px';` +
		`document.body.appendChild(ta);ta.focus();ta.select();` +
		`try{document.execCommand('copy');done();}catch(e){}` +
		`document.body.removeChild(ta);}` +
		`document.addEventListener('click',function(e){` +
		`var btn=e.target.closest('[data-tc-copy]');if(!btn)return;` +
		`e.preventDefault();` +
		`var text=btn.getAttribute('data-tc-copy');` +
		`var label=btn.getAttribute('data-tc-copy-label')||'Copied!';` +
		`var labelEl=btn.querySelector('[data-tc-copy-text]');` +
		`var original=labelEl?labelEl.textContent:'';` +
		`if(navigator.clipboard&&navigator.clipboard.writeText){` +
		`navigator.clipboard.writeText(text).then(function(){` +
		`tcSwapLabel(labelEl,label,original);}).catch(function(){` +
		`tcFallbackCopy(text,function(){tcSwapLabel(labelEl,label,original);});});` +
		`}else{tcFallbackCopy(text,function(){tcSwapLabel(labelEl,label,original);});}` +
		`});}` +
		"\n"
}

// copyButtonScriptComponent renders a <script nonce="..."> tag containing the
// clipboard copy JS. Singleton — only the first CopyButton on the page injects
// executable code; subsequent instances skip via the window.tcCopyAttached guard.
func copyButtonScriptComponent(nonce string) templ.Component {
	return scriptComponent(nonce, copyButtonJS(), "copy button script")
}

// imageFallbackJS returns the singleton JavaScript for image fallback source
// swapping. Uses event capture (true) because the error event does not bubble.
// When an img with data-tc-img-fallback fails to load, the src is swapped to
// the fallback and the attribute is removed to prevent infinite loops.
func imageFallbackJS() string {
	return `if(!window.tcImageFallbackAttached){window.tcImageFallbackAttached=true;` +
		`document.addEventListener('error',function(e){` +
		`var img=e.target;` +
		`if(img&&img.dataset&&img.dataset.tcImgFallback){` +
		`img.onerror=null;` +
		`img.src=img.dataset.tcImgFallback;` +
		`delete img.dataset.tcImgFallback;` +
		`}` +
		`},true);}` +
		"\n"
}

// imageFallbackScriptComponent renders a <script nonce="..."> tag containing
// the image fallback JS. Singleton — only injected when at least one Image
// with FallbackSrc is rendered on the page.
func imageFallbackScriptComponent(nonce string) templ.Component {
	return scriptComponent(nonce, imageFallbackJS(), "image fallback script")
}

// relativeTimeJS returns the singleton JavaScript for live-updating relative
// timestamps. Uses the native Intl.RelativeTimeFormat API for localization.
// Runs every 30 seconds via setInterval and also listens to htmx:afterSettle
// so newly-swapped <time> elements are formatted immediately after HTMX
// content swaps.
func relativeTimeJS() string {
	return `if(!window.tcRelativeTimeAttached){window.tcRelativeTimeAttached=true;` +
		`function tcRelativeTimeFormat(el){` +
		`var diff=(new Date(el.getAttribute('datetime'))-Date.now())/1000;` +
		`var abs=Math.abs(diff),val,unit;` +
		`if(abs<60){val=0;unit='second'}` +
		`else if(abs<3600){val=Math.round(diff/60);unit='minute'}` +
		`else if(abs<86400){val=Math.round(diff/3600);unit='hour'}` +
		`else if(abs<604800){val=Math.round(diff/86400);unit='day'}` +
		`else if(abs<2592000){val=Math.round(diff/604800);unit='week'}` +
		`else{el.textContent=new Date(el.getAttribute('datetime')).toLocaleDateString();return}` +
		`el.textContent=new Intl.RelativeTimeFormat(document.documentElement.lang||'en',{numeric:'auto'}).format(val,unit);` +
		`}` +
		`function tcRelativeTimeTick(){` +
		`document.querySelectorAll('time[data-tc-relative]').forEach(tcRelativeTimeFormat);` +
		`}` +
		`tcRelativeTimeTick();` +
		`setInterval(tcRelativeTimeTick,30000);` +
		`document.body.addEventListener('htmx:afterSettle',tcRelativeTimeTick);` +
		`}` +
		"\n"
}

// relativeTimeScriptComponent renders the relative-time auto-refresh JS.
// Singleton — only the first RelativeTime with AutoRefresh=true on the page
// injects executable code; subsequent instances skip via the guard flag.
func relativeTimeScriptComponent(nonce string) templ.Component {
	return scriptComponent(nonce, relativeTimeJS(), "relative time script")
}

// tableRowHrefJS returns the singleton JavaScript for clickable table rows.
// When any <tr> has data-tc-row-href, clicking it (or pressing Enter/Space when
// focused) navigates to the URL. Clicks on interactive elements inside the row
// (links, buttons, inputs) are NOT hijacked so they work normally.
func tableRowHrefJS() string {
	return `if(!window.tcTableRowHrefAttached){window.tcTableRowHrefAttached=true;` +
		`document.addEventListener('click',function(e){` +
		`var row=e.target.closest('tr[data-tc-row-href]');` +
		`if(!row)return;` +
		`if(e.target.closest('a[href],button:not([disabled]),input,select,textarea,[contenteditable]'))return;` +
		`window.location.href=row.dataset.tcRowHref;` +
		`});` +
		`document.addEventListener('keydown',function(e){` +
		`var row=e.target.closest('tr[data-tc-row-href]');` +
		`if(!row)return;` +
		`if(e.key==='Enter'||e.key===' '){` +
		`if(e.target.closest('a[href],button:not([disabled]),input,select,textarea,[contenteditable]'))return;` +
		`e.preventDefault();` +
		`window.location.href=row.dataset.tcRowHref;` +
		`}` +
		`});` +
		`document.body.addEventListener('htmx:afterSettle',function(){` +
		`document.querySelectorAll('tr[data-tc-row-href]').forEach(function(r){` +
		`if(!r.hasAttribute('tabindex'))r.setAttribute('tabindex','0');` +
		`if(!r.hasAttribute('role'))r.setAttribute('role','link');` +
		`});` +
		`});` +
		`}` +
		"\n"
}

// tableRowHrefScriptComponent renders the clickable-row JS. Only injected when
// at least one TableRow has Href set.
func tableRowHrefScriptComponent(nonce string) templ.Component {
	return scriptComponent(nonce, tableRowHrefJS(), "table row href script")
}
