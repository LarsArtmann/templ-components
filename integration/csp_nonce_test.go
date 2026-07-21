package integration

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/errorpage"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

// TestAllInlineScriptsHaveNonce renders every component that uses inline scripts
// and asserts that every <script> tag contains the nonce attribute. This prevents
// CSP regressions where a new or refactored component forgets the nonce.
func TestAllInlineScriptsHaveNonce(t *testing.T) {
	t.Parallel()

	const testNonce = "test-csp-nonce-12345"

	// Each entry renders a component that emits at least one inline <script>.
	renderings := []struct {
		name string
		html string
	}{
		{"Accordion", utils.Render(t, display.Accordion(display.AccordionProps{
			BaseProps: utils.BaseProps{Nonce: testNonce},
			Items:     []display.AccordionItem{{ID: "a1", Title: "A"}},
		}))},
		{"Modal", utils.Render(t, display.Modal(display.ModalProps{
			BaseProps: utils.BaseProps{ID: "m1", Nonce: testNonce},
			Title:     "Test",
		}))},
		{"Drawer", utils.Render(t, display.Drawer(display.DrawerProps{
			BaseProps: utils.BaseProps{ID: "dr1", Nonce: testNonce},
			Title:     "Test",
		}))},
		{"Dropdown", utils.Render(t, display.Dropdown(display.DropdownProps{
			BaseProps: utils.BaseProps{ID: "dd", Nonce: testNonce},
			Label:     "Menu",
			Items:     []display.DropdownItem{{Text: "X", Href: "/x"}},
		}))},
		{"ContextMenu", utils.Render(t, display.ContextMenu(display.ContextMenuProps{
			BaseProps: utils.BaseProps{ID: "cm", Nonce: testNonce},
			Items:     []display.ContextMenuItem{{Text: "Edit", Href: "/edit"}},
		}))},
		{"Tabs", utils.Render(t, display.Tabs(display.TabsProps{
			BaseProps:  utils.BaseProps{Nonce: testNonce},
			Tabs:       []display.Tab{{ID: "t1", Label: "Tab1"}},
			ClientSide: true,
		}))},
		{"Alert", utils.Render(t, feedback.Alert(feedback.AlertProps{
			BaseProps: utils.BaseProps{Nonce: testNonce},
			Type:      feedback.FeedbackInfo,
			Title:     "Info",
		}))},
		{"Toast", utils.Render(t, feedback.Toast(feedback.ToastProps{
			BaseProps: utils.BaseProps{Nonce: testNonce},
			Type:      feedback.FeedbackSuccess,
			Message:   "OK",
		}))},
		{"CopyButton", utils.Render(t, display.CopyButton(display.CopyButtonProps{
			BaseProps: utils.BaseProps{Nonce: testNonce},
			Text:      "copy me",
		}))},
		{"GlobalErrorHandling", utils.Render(t, htmx.GlobalErrorHandling(htmx.ErrorHandlingConfig{
			Nonce: testNonce,
		}))},
		{"ThemeScript", utils.Render(t, layout.ThemeScript(testNonce))},
		{"ThemeToggle", utils.Render(t, layout.ThemeToggle("Toggle theme", testNonce))},
		{"MobileMenu", utils.Render(t, navigation.MobileMenu(nil, "/", testNonce, "mm", false))},
		{"ErrorPage", utils.Render(t, errorpage.ErrorPage(errorpage.ErrorPageProps{
			BaseProps: utils.BaseProps{Nonce: testNonce},
			Title:     "Error",
		}))},
		{"NotFound404", utils.Render(t, errorpage.NotFound404(errorpage.NotFound404Props{
			BaseProps: utils.BaseProps{Nonce: testNonce},
		}))},
		{"TableWithRowHref", utils.Render(t, display.Table(display.TableProps{
			BaseProps: utils.BaseProps{Nonce: testNonce},
			Headers:   []string{"Name"},
			Rows: []display.TableRow{
				{Cells: []display.TableCell{{Text: "Alice"}}, Href: "/users/1"},
			},
		}))},
	}

	for _, r := range renderings {
		t.Run(r.name, func(t *testing.T) {
			t.Parallel()

			scriptCount := strings.Count(r.html, "<script")
			if scriptCount == 0 {
				t.Skipf("no inline scripts in output")
			}

			nonceCount := strings.Count(r.html, `nonce="`+testNonce+`"`)
			if nonceCount < scriptCount {
				t.Errorf("%s: found %d <script> tags but only %d have nonce=%q\n%s",
					r.name, scriptCount, nonceCount, testNonce, r.html)
			}
		})
	}
}
