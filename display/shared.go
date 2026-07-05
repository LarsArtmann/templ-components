package display

import (
	"context"
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

// focusableSelector is the CSS selector for focusable elements within an
// overlay panel. Used for focus trapping (Tab cycling) and initial focus.
const focusableSelector = `'a[href], button:not([disabled]), textarea, input, select, [tabindex]:not([tabindex="-1"])'`

// mutedTextClass is the shared Tailwind class string for secondary/muted text
// (Card subtitle, PageHeader subtitle, EmptyState description). Callers combine
// it with a margin utility via utils.Class, e.g. utils.Class(mutedTextClass, "mt-1").
const mutedTextClass = "text-sm text-gray-500 dark:text-gray-400"

// overlayPanelConfig defines the CSS classes that animate the panel element
// between open and closed states. The fields are typed slices, converted to
// JS-quoted classList.add/remove arguments via jsClassArgs at use site.
type overlayPanelConfig struct {
	openClasses  []string // e.g. {"scale-100", "opacity-100"}
	closeClasses []string // e.g. {"scale-95", "opacity-0"}
}

// overlayShellProps parameterizes the shared accessibility shell used by
// Modal and Drawer. It centralizes the outer overlay div, backdrop, and
// per-instance script so the two components stay consistent.
type overlayShellProps struct {
	id            string
	open          bool
	title         string
	ariaLabel     string
	closeKind     string
	componentName string
	outerClass    string
	nonce         string
	cfg           overlayPanelConfig
}

// jsClassArgs converts a slice of CSS class names into comma-separated,
// single-quoted JS string arguments suitable for classList.add/remove.
// Example: []string{"scale-100", "opacity-100"} -> "'scale-100', 'opacity-100'"
func jsClassArgs(classes []string) string {
	quoted := make([]string, len(classes))
	for i, c := range classes {
		quoted[i] = "'" + c + "'"
	}
	return strings.Join(quoted, ", ")
}

func modalPanelConfig() overlayPanelConfig {
	return overlayPanelConfig{
		openClasses:  []string{"scale-100", "opacity-100"},
		closeClasses: []string{"scale-95", "opacity-0"},
	}
}

func drawerPanelConfig(side DrawerSide) overlayPanelConfig {
	if side == DrawerLeft {
		return overlayPanelConfig{
			openClasses:  []string{"translate-x-0"},
			closeClasses: []string{"-translate-x-full"},
		}
	}
	return overlayPanelConfig{
		openClasses:  []string{"translate-x-0"},
		closeClasses: []string{"translate-x-full"},
	}
}

// overlayCloseJS generates the tcClose{Name} function body.
func overlayCloseJS(componentName string, cfg overlayPanelConfig) string {
	return "function tcClose" + componentName + "(id) {\n" +
		"\tvar overlay = document.getElementById(id);\n" +
		"\tif (overlay) {\n" +
		"\t\toverlay.classList.remove('opacity-100', 'pointer-events-auto');\n" +
		"\t\toverlay.classList.add('opacity-0', 'pointer-events-none');\n" +
		"\t\tvar panel = document.getElementById(id + '-panel');\n" +
		"\t\tif (panel) {\n" +
		"\t\t\tpanel.classList.remove(" + jsClassArgs(cfg.openClasses) + ");\n" +
		"\t\t\tpanel.classList.add(" + jsClassArgs(cfg.closeClasses) + ");\n" +
		"\t\t}\n" +
		"\t\tvar prevId = overlay.getAttribute('data-tc-prev-focus');\n" +
		"\t\tif (prevId) {\n" +
		"\t\t\tvar prev = document.getElementById(prevId);\n" +
		"\t\t\tif (prev) prev.focus();\n" +
		"\t\t\toverlay.removeAttribute('data-tc-prev-focus');\n" +
		"\t\t}\n" +
		"\t}\n" +
		"}\n"
}

// overlayOpenJS generates the tcOpen{Name} function body.
func overlayOpenJS(componentName string, cfg overlayPanelConfig) string {
	return "function tcOpen" + componentName + "(id) {\n" +
		"\tvar overlay = document.getElementById(id);\n" +
		"\tif (overlay) {\n" +
		"\t\tif (document.activeElement && document.activeElement.id) {\n" +
		"\t\t\toverlay.setAttribute('data-tc-prev-focus', document.activeElement.id);\n" +
		"\t\t}\n" +
		"\t\toverlay.classList.remove('opacity-0', 'pointer-events-none');\n" +
		"\t\toverlay.classList.add('opacity-100', 'pointer-events-auto');\n" +
		"\t\tvar panel = document.getElementById(id + '-panel');\n" +
		"\t\tif (panel) {\n" +
		"\t\t\tpanel.classList.remove(" + jsClassArgs(cfg.closeClasses) + ");\n" +
		"\t\t\tpanel.classList.add(" + jsClassArgs(cfg.openClasses) + ");\n" +
		"\t\t\tvar focusable = panel.querySelectorAll(" + focusableSelector + ");\n" +
		"\t\t\tif (focusable.length) focusable[0].focus();\n" +
		"\t\t}\n" +
		"\t}\n" +
		"}\n"
}

// overlayTrapJS generates the per-instance IIFE for click delegation and focus trap.
func overlayTrapJS(id, componentName string) string {
	closeFn := "tcClose" + componentName
	openFn := "tcOpen" + componentName
	escapedID := strconv.Quote(id)
	return "(function(id) {\n" +
		"\tvar overlay = document.getElementById(id);\n" +
		"\tif (!overlay) return;\n" +
		"\toverlay.addEventListener('click', function(e) {\n" +
		"\t\tif (e.target.closest('[data-tc-close]')) { " + closeFn + "(id); }\n" +
		"\t});\n" +
		"\toverlay.addEventListener('keydown', function(e) {\n" +
		"\t\tif (e.key === 'Escape') { " + closeFn + "(id); return; }\n" +
		"\t\tif (e.key !== 'Tab') return;\n" +
		"\t\tvar panel = document.getElementById(id + '-panel');\n" +
		"\t\tif (!panel) return;\n" +
		"\t\tvar focusable = panel.querySelectorAll(" + focusableSelector + ");\n" +
		"\t\tif (focusable.length === 0) { e.preventDefault(); return; }\n" +
		"\t\tvar first = focusable[0];\n" +
		"\t\tvar last = focusable[focusable.length - 1];\n" +
		"\t\tif (e.shiftKey && document.activeElement === first) {\n" +
		"\t\t\te.preventDefault(); last.focus();\n" +
		"\t\t} else if (!e.shiftKey && document.activeElement === last) {\n" +
		"\t\t\te.preventDefault(); first.focus();\n" +
		"\t\t}\n" +
		"\t});\n" +
		"\tif (overlay.getAttribute('data-tc-open-on-load') === 'true') { " + openFn + "(id); }\n" +
		"})(" + escapedID + ");\n"
}

// overlayJS generates the complete JavaScript for an overlay component
// (Modal or Drawer). It produces:
//   - tcClose{Name}(id): hide overlay, animate panel closed, restore focus
//   - tcOpen{Name}(id): show overlay, animate panel open, save focus, focus first
//   - Per-instance IIFE: [data-tc-close] click delegation + Tab focus trap + Escape
//
// This replaces inline onclick handlers for CSP compliance (no script-src-attr
// needed) and eliminates ~90 lines of duplicated JS between Modal and Drawer.
func overlayJS(id, componentName string, cfg overlayPanelConfig) string {
	return overlayCloseJS(componentName, cfg) +
		overlayOpenJS(componentName, cfg) +
		overlayTrapJS(id, componentName)
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
func overlayScriptComponent(nonce, id, componentName string, cfg overlayPanelConfig) templ.Component {
	return scriptComponent(nonce, overlayJS(id, componentName, cfg), "overlay script")
}

// tooltipJS returns the singleton JavaScript for tooltip touch support and
// Escape-to-dismiss. Uses a window.tcTooltipAttached guard so it executes
// only once per page regardless of how many Tooltip components are rendered.
func tooltipJS() string {
	return `if(!window.tcTooltipAttached){window.tcTooltipAttached=true;` +
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
