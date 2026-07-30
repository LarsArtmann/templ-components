package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestExternalLinkText(t *testing.T) {
	t.Parallel()
	props := DefaultExternalLinkProps()
	props.Href = "https://example.com"
	props.Text = "Visit site"
	output := utils.Render(t, ExternalLink(props))
	utils.AssertContains(t, output, `href="https://example.com"`)
	utils.AssertContains(t, output, `target="_blank"`)
	utils.AssertContains(t, output, `rel="noopener noreferrer"`)
	utils.AssertContains(t, output, "Visit site")
	utils.AssertContains(t, output, "&#8599;")
}

func TestExternalLinkChildren(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ExternalLink(ExternalLinkProps{
		Href: "https://docs.example.com",
		BaseProps: utils.BaseProps{
			AriaLabel: "Documentation",
		},
	}))
	utils.AssertContains(t, output, `href="https://docs.example.com"`)
	utils.AssertContains(t, output, `aria-label="Documentation"`)
}

func TestExternalLinkHideIcon(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ExternalLink(ExternalLinkProps{
		Href:     "https://example.com",
		Text:     "Link",
		ShowIcon: false,
	}))
	utils.AssertNotContains(t, output, "&#8599;")
}

func TestExternalLinkSanitizesHref(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ExternalLink(ExternalLinkProps{
		Href: "javascript:alert('xss')",
		Text: "Malicious",
	}))
	// templ sanitizes dangerous schemes by rewriting to about:invalid
	utils.AssertNotContains(t, output, `href="javascript:`)
}

func TestExternalLinkDefaults(t *testing.T) {
	t.Parallel()
	props := DefaultExternalLinkProps()
	if !props.ShowIcon {
		t.Error("expected ShowIcon=true")
	}
}

func TestExternalLinkClass(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ExternalLink(ExternalLinkProps{
		Href: "https://example.com",
		Text: "Link",
		BaseProps: utils.BaseProps{
			Class: "text-blue-600 dark:text-blue-400",
		},
	}))
	utils.AssertContains(t, output, "text-blue-600")
}
