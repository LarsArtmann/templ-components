package icons

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func renderCustom(t *testing.T, icon CustomIcon, class string) string {
	t.Helper()

	return utils.Render(t, Render(icon, class))
}

func TestRender_Defaults(t *testing.T) {
	t.Parallel()

	out := renderCustom(t, CustomIcon{Paths: []string{"M3 3h18v18H3z"}}, "h-4 w-4")

	for _, want := range []string{
		`class="h-4 w-4"`,
		`viewBox="0 0 24 24"`,
		`fill="currentColor"`,
		`d="M3 3h18v18H3z"`,
		`aria-hidden="true"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("default render should contain %q", want)
		}
	}
}

func TestRender_CustomViewBoxAndFill(t *testing.T) {
	t.Parallel()

	out := renderCustom(t, CustomIcon{
		ViewBox: "0 0 128 128",
		Paths:   []string{"M10 10h108"},
		Fill:    "#316192",
	}, "h-6 w-6")

	if !strings.Contains(out, `viewBox="0 0 128 128"`) {
		t.Error("custom viewBox should render")
	}

	if !strings.Contains(out, `fill="#316192"`) {
		t.Error("custom fill (brand hex) should render")
	}
}

func TestRender_MultiPathSkipsEmpty(t *testing.T) {
	t.Parallel()

	out := renderCustom(t, CustomIcon{Paths: []string{"M1 1", "", "M2 2"}}, "")

	if !strings.Contains(out, `d="M1 1"`) || !strings.Contains(out, `d="M2 2"`) {
		t.Error("both non-empty paths should render")
	}

	if strings.Count(out, "<path") != 2 {
		t.Errorf("empty path should be skipped, got %d path nodes", strings.Count(out, "<path"))
	}
}

func TestRender_TitleAccessibility(t *testing.T) {
	t.Parallel()

	decorative := renderCustom(t, CustomIcon{Paths: []string{"M1 1"}}, "")
	if strings.Contains(decorative, "<title>") || !strings.Contains(decorative, `aria-hidden="true"`) {
		t.Error("title-less icon must stay aria-hidden with no <title>")
	}

	labelled := renderCustom(t, CustomIcon{Paths: []string{"M1 1"}, Title: "Go logo"}, "")
	if !strings.Contains(labelled, "<title>Go logo</title>") || strings.Contains(labelled, `aria-hidden="true"`) {
		t.Error("titled icon should render <title> and drop aria-hidden")
	}

	if !strings.Contains(labelled, `role="img"`) {
		t.Error("titled icon should carry role=img")
	}
}

func TestRender_SplitIconPathJSConvention(t *testing.T) {
	t.Parallel()

	joined := "M1 1|M2 2"
	out := renderCustom(t, CustomIcon{Paths: strings.Split(joined, "|")}, "")

	if strings.Count(out, "<path") != 2 {
		t.Error("IconPathJS-style | split should produce one path per segment")
	}
}
