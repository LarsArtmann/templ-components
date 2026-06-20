package display

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/a-h/templ"
)

// focusableSelector is the CSS selector for focusable elements within an
// overlay panel. Used for focus trapping (Tab cycling) and initial focus.
const focusableSelector = `'a[href], button:not([disabled]), textarea, input, select, [tabindex]:not([tabindex="-1"])'`

// overlayPanelConfig defines the CSS classes that animate the panel element
// between open and closed states. The fields are comma-separated JS class
// string arguments suitable for classList.add/remove.
type overlayPanelConfig struct {
	openClasses  string // e.g. "'scale-100', 'opacity-100'"
	closeClasses string // e.g. "'scale-95', 'opacity-0'"
}

func modalPanelConfig() overlayPanelConfig {
	return overlayPanelConfig{
		openClasses:  "'scale-100', 'opacity-100'",
		closeClasses: "'scale-95', 'opacity-0'",
	}
}

func drawerPanelConfig(side DrawerSide) overlayPanelConfig {
	if side == DrawerLeft {
		return overlayPanelConfig{
			openClasses:  "'translate-x-0'",
			closeClasses: "'-translate-x-full'",
		}
	}
	return overlayPanelConfig{
		openClasses:  "'translate-x-0'",
		closeClasses: "'translate-x-full'",
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
		"\t\t\tpanel.classList.remove(" + cfg.openClasses + ");\n" +
		"\t\t\tpanel.classList.add(" + cfg.closeClasses + ");\n" +
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
		"\t\t\tpanel.classList.remove(" + cfg.closeClasses + ");\n" +
		"\t\t\tpanel.classList.add(" + cfg.openClasses + ");\n" +
		"\t\t\tvar focusable = panel.querySelectorAll(" + focusableSelector + ");\n" +
		"\t\t\tif (focusable.length) focusable[0].focus();\n" +
		"\t\t}\n" +
		"\t}\n" +
		"}\n"
}

// overlayTrapJS generates the per-instance IIFE for click delegation and focus trap.
func overlayTrapJS(id, componentName string) string {
	closeFn := "tcClose" + componentName
	escapedID := strconv.Quote(id)
	return "(function(id) {\n" +
		"\tvar overlay = document.getElementById(id);\n" +
		"\tif (!overlay) return;\n" +
		"\toverlay.querySelectorAll('[data-tc-close]').forEach(function(el) {\n" +
		"\t\tel.addEventListener('click', function() { " + closeFn + "(id); });\n" +
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
func overlayScriptComponent(nonce, id, componentName string, cfg overlayPanelConfig) templ.Component {
	js := overlayJS(id, componentName, cfg)
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		if _, err := fmt.Fprintf(w, "<script nonce=\"%s\">\n%s</script>\n", nonce, js); err != nil {
			return fmt.Errorf("write overlay script: %w", err)
		}
		return nil
	})
}
