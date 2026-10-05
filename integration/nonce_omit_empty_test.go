package integration

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/charts/echarts"
	"github.com/larsartmann/templ-components/datastar"
	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/errorpage"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/layout"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// TestNoEmptyNonceAttribute renders every inline-script component with an
// EMPTY nonce and asserts the output never contains `nonce=""`.
//
// A `nonce=""` attribute is worse than useless on a strict-CSP page: the
// browser rejects the script (a nonce attribute never matches the header's
// nonce unless equal), so the component renders fine but its feature is
// silently dead. The omit-empty rule (utils.ScriptAttrs) renders NO nonce
// attribute when none was provided — the script still runs for non-CSP
// consumers and is blocked there either way for strict-CSP ones.
//
// New script-emitting components MUST be added to this table (mirrored by
// TestAllInlineScriptsHaveNonce for the with-nonce direction).
func TestNoEmptyNonceAttribute(t *testing.T) {
	t.Parallel()

	// Each entry renders a script-emitting component with an unset nonce.
	renderings := []struct {
		name string
		html string
	}{
		{"DatastarSDKScript", utils.Render(t, datastar.SDKScript(datastar.SDKScriptProps{}))},
		{"Accordion", utils.Render(t, display.Accordion(display.AccordionProps{
			Items: []display.AccordionItem{{ID: "a1", Title: "A"}},
		}))},
		{"Modal", utils.Render(t, display.Modal(display.ModalProps{
			ID: "m1", Title: "Test",
		}))},
		{"Drawer", utils.Render(t, display.Drawer(display.DrawerProps{
			ID: "dr1", Title: "Test",
		}))},
		{"Dropdown", utils.Render(t, display.Dropdown(display.DropdownProps{
			ID: "dd", Label: "Menu",
			Items: []display.DropdownItem{{Text: "X", Href: "/x"}},
		}))},
		{"ContextMenu", utils.Render(t, display.ContextMenu(display.ContextMenuProps{
			ID:    "cm",
			Items: []display.ContextMenuItem{{Text: "Edit", Href: "/edit"}},
		}))},
		{"Tabs", utils.Render(t, display.Tabs(display.TabsProps{
			Tabs: []display.Tab{{ID: "t1", Label: "Tab1"}}, ClientSide: true,
		}))},
		{"KanbanBoardWired", utils.Render(t, display.KanbanBoard(display.KanbanBoardProps{
			Columns: []display.KanbanColumn{{ID: "c1", Title: "Todo"}},
			Wire:    wire.Post("/api/kanban"),
		}))},
		{"CopyButton", utils.Render(t, display.CopyButton(display.CopyButtonProps{
			Text: "copy me",
		}))},
		{"TableWithRowHref", utils.Render(t, display.Table(display.TableProps{
			Headers: []string{"Name"},
			Rows: []display.TableRow{
				{Cells: []display.TableCell{{Text: "Alice"}}, Href: "/users/1"},
			},
		}))},
		{"Alert", utils.Render(t, feedback.Alert(feedback.AlertProps{
			Type: feedback.FeedbackInfo, Title: "Info",
		}))},
		{"Toast", utils.Render(t, feedback.Toast(feedback.ToastProps{
			Type: feedback.FeedbackSuccess, Message: "OK",
		}))},
		{"ToastContainer", utils.Render(t, feedback.ToastContainer(""))},
		{"GlobalErrorHandling", utils.Render(t, htmx.GlobalErrorHandling(htmx.ErrorHandlingConfig{}))},
		{"ViewTransitions", utils.Render(t, htmx.ViewTransitions(htmx.ViewTransitionsProps{Global: true}))},
		{"PolledRegionEager", utils.Render(t, htmx.PolledRegion(htmx.PolledRegionProps{
			URL: "/stats", Every: "10s", Eager: true,
		}))},
		{"SSEErrorHandling", utils.Render(t, datastar.SSEErrorHandling(datastar.SSEErrorHandlingConfig{}))},
		{"DatastarLiveRegion", utils.Render(t, datastar.LiveRegion(datastar.LiveRegionProps{
			URL: "/api/stream",
		}))},
		{"DirtyGuard", utils.Render(t, forms.DirtyGuard(forms.DirtyGuardProps{}))},
		{"Combobox", utils.Render(t, forms.Combobox(forms.ComboboxProps{
			Name: "color", Label: "Color",
			Options: []forms.ComboboxOption{{Value: "red", Label: "Red"}},
		}))},
		{"TagsInput", utils.Render(t, forms.TagsInput(forms.TagsInputProps{
			Name: "tags", Label: "Tags",
		}))},
		{
			"MobileMenu",
			utils.Render(
				t,
				navigation.MobileMenu([]navigation.NavLinkProps{{Label: "Home", Href: "/"}}, "/", "", "mm", false),
			),
		},
		{"ErrorPage", utils.Render(t, errorpage.ErrorPage(errorpage.ErrorPageProps{
			Title: "Error",
		}))},
		{"NotFound404", utils.Render(t, errorpage.NotFound404(errorpage.NotFound404Props{}))},
		{"ErrorAlert", utils.Render(t, errorpage.ErrorAlert(errorpage.ErrorAlertProps{
			Dismissible: true,
		}))},
		{"EChart", utils.Render(t, echarts.EChart(echarts.EChartsProps{
			BaseProps:      utils.BaseProps{},
			Element:        `<div id="ec1"></div>`,
			Script:         `echarts.init(document.getElementById('ec1'));`,
			DarkModeBridge: true,
		}))},
		{"EChartsSDKScript", utils.Render(t, echarts.SDKScript(echarts.SDKScriptProps{}))},
		{"ThemeScript", utils.Render(t, layout.ThemeScript(""))},
		{"ThemeToggle", utils.Render(t, layout.ThemeToggle("Toggle theme", ""))},
		{"BaseSelfHostedHTMX", utils.Render(t, layout.Base(layout.PageProps{
			Title: "t", Nonce: "",
		}))},
	}

	for _, r := range renderings {
		t.Run(r.name, func(t *testing.T) {
			t.Parallel()

			if strings.Contains(r.html, `nonce=""`) {
				t.Errorf(
					"%s renders nonce=\"\" with an empty Nonce — apply the omit-empty rule (utils.ScriptAttrs or the scriptComponent pattern)\n%s",
					r.name,
					r.html,
				)
			}
		})
	}
}

// TestEmptyNonceStillRendersScripts pins the other half of the omit-empty
// contract: dropping the nonce ATTRIBUTE must not drop the SCRIPT. These
// emitters render their JS unnonced when no nonce is set (theme scripts are
// the documented exception — they omit the whole tag on empty nonce).
func TestEmptyNonceStillRendersScripts(t *testing.T) {
	t.Parallel()

	mustContainScript := map[string]string{
		"Modal":               utils.Render(t, display.Modal(display.ModalProps{ID: "m1", Title: "T"})),
		"CopyButton":          utils.Render(t, display.CopyButton(display.CopyButtonProps{Text: "x"})),
		"GlobalErrorHandling": utils.Render(t, htmx.GlobalErrorHandling(htmx.ErrorHandlingConfig{})),
		"SSEErrorHandling":    utils.Render(t, datastar.SSEErrorHandling(datastar.SSEErrorHandlingConfig{})),
		"MobileMenu": utils.Render(
			t,
			navigation.MobileMenu([]navigation.NavLinkProps{{Label: "Home", Href: "/"}}, "/", "", "mm", false),
		),
		"Combobox": utils.Render(t, forms.Combobox(forms.ComboboxProps{
			Name: "c", Label: "C", Options: []forms.ComboboxOption{{Value: "a", Label: "A"}},
		})),
		"Toast": utils.Render(t, feedback.Toast(feedback.ToastProps{Type: feedback.FeedbackInfo, Message: "m"})),
	}

	for name, html := range mustContainScript {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if !strings.Contains(html, "<script>") {
				t.Errorf(
					"%s dropped its inline script entirely on empty nonce — omit the nonce attribute, not the script\n%s",
					name,
					html,
				)
			}
		})
	}
}
