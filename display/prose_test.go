package display

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func proseSample() templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<p>Hello <strong>world</strong></p>")

		return err
	})
}

// TestProseSizeIsValid pins the ProseSize closed set (repo rule: every enum
// ships IsValid + a test in the SAME commit).
func TestProseSizeIsValid(t *testing.T) {
	t.Parallel()

	for _, valid := range []ProseSize{ProseSizeSM, ProseSizeMD, ProseSizeLG} {
		if !ProseSizeIsValid(valid) {
			t.Errorf("ProseSizeIsValid(%q) = false, want true", valid)
		}
	}

	for _, invalid := range []ProseSize{"", "xl", "small", "SM"} {
		if ProseSizeIsValid(invalid) {
			t.Errorf("ProseSizeIsValid(%q) = true, want false", invalid)
		}
	}
}

// TestProseSizeClassZeroValueEqualsDefault pins the struct-literal rule: the
// zero value resolves to the same classes as the documented default.
func TestProseSizeClassZeroValueEqualsDefault(t *testing.T) {
	t.Parallel()

	if proseSizeClass("") != proseSizeClass(ProseSizeMD) {
		t.Errorf("zero value class %q != default class %q", proseSizeClass(""), proseSizeClass(ProseSizeMD))
	}

	if proseSizeClass(ProseSizeSM) == "" {
		t.Error("ProseSizeSM must carry size classes")
	}
}

// TestProseRender pins the wrapper markup: the tc-prose hook class, the
// size class merge, BaseProps propagation, and children rendering.
func TestProseRender(t *testing.T) {
	t.Parallel()

	content := proseSample()

	tests := []struct {
		name     string
		props    ProseProps
		contains []string
		absent   []string
	}{
		{
			name:     "default size renders hook class without size class",
			props:    ProseProps{Children: content},
			contains: []string{`class="tc-prose"`, "<p>Hello <strong>world</strong></p>"},
		},
		{
			name:     "small size class merges",
			props:    ProseProps{Children: content, Size: ProseSizeSM},
			contains: []string{"tc-prose", "text-sm", "leading-relaxed"},
		},
		{
			name:     "large size class merges",
			props:    ProseProps{Children: content, Size: ProseSizeLG},
			contains: []string{"tc-prose", "text-lg", "leading-relaxed"},
		},
		{
			name: "id, class, and attrs propagate",
			props: ProseProps{
				Children: content,
				BaseProps: utils.BaseProps{
					ID:    "article-body",
					Class: "max-w-3xl",
					Attrs: templ.Attributes{"data-article": "true"},
				},
			},
			contains: []string{`id="article-body"`, `max-w-3xl`, `data-article="true"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			html := utils.Render(t, Prose(tt.props))
			utils.AssertContainsAll(t, html, tt.contains...)

			for _, unwanted := range tt.absent {
				if strings.Contains(html, unwanted) {
					t.Errorf("rendered HTML must not contain %q:\n%s", unwanted, html)
				}
			}
		})
	}
}
