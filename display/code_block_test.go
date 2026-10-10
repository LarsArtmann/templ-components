package display

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

func TestCodeBlockRender(t *testing.T) {
	t.Parallel()

	t.Run("block variant with language header and copy", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{Code: "go install ./...", Language: "bash"}))
		utils.AssertContains(t, output, "<figure")
		utils.AssertContains(t, output, "<figcaption")
		utils.AssertContains(t, output, "bash")
		utils.AssertContains(t, output, "<pre")
		utils.AssertContains(t, output, "<code")
		utils.AssertContains(t, output, "go install ./...")
		utils.AssertContains(t, output, "data-tc-copy=\"go install ./...\"")
	})

	t.Run("compact id variant", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{
			Code:    "123456789012345678",
			Variant: CodeBlockCompactID,
		}))
		utils.AssertContains(t, output, "font-mono")
		utils.AssertContains(t, output, "123456789012345678")
		utils.AssertContains(t, output, "data-tc-copy=\"123456789012345678\"")
		utils.AssertNotContains(t, output, "<pre")
	})

	t.Run("no copy renders no button", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{Code: "x", NoCopy: true}))
		utils.AssertNotContains(t, output, "data-tc-copy")
	})

	t.Run("label overrides language", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{
			Code:  "SELECT 1",
			Label: "Query",
		}))
		utils.AssertContains(t, output, "Query")
	})

	t.Run("no header when no label and no copy", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{Code: "plain", NoCopy: true}))
		utils.AssertNotContains(t, output, "<figcaption")
	})

	t.Run("unknown variant falls back to block", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{Code: "x", Variant: CodeBlockVariant("bogus")}))
		utils.AssertContains(t, output, "<figure")
	})

	t.Run("code is html escaped", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{Code: "<script>alert(1)</script>"}))
		utils.AssertNotContains(t, output, "<script>alert")
		utils.AssertContains(t, output, "&lt;script&gt;")
	})

	t.Run("propagates id class aria-label and attrs", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{
			ID:        "snippet",
			Class:     "my-4",
			AriaLabel: "Install command",
			Attrs:     templ.Attributes{"data-foo": "bar"},
			Code:      "x",
		}))
		utils.AssertContains(t, output, `id="snippet"`)
		utils.AssertContains(t, output, "my-4")
		utils.AssertContains(t, output, `aria-label="Install command"`)
		utils.AssertContains(t, output, `data-foo="bar"`)
	})

	t.Run("nonce threads to the copy script", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{
			Nonce: "nonce-abc",
			Code:  "x",
		}))
		utils.AssertContains(t, output, `nonce="nonce-abc"`)
	})

	t.Run("empty nonce renders script without empty attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, CodeBlock(CodeBlockProps{Code: "x"}))
		utils.AssertNotContains(t, output, `nonce=""`)
		utils.AssertContains(t, output, "<script")
	})
}

func TestCodeBlockLabelResolution(t *testing.T) {
	t.Parallel()

	t.Run("label wins over language", func(t *testing.T) {
		t.Parallel()

		if got := codeBlockLabel(CodeBlockProps{Label: "L", Language: "go"}); got != "L" {
			t.Errorf("codeBlockLabel(Label=L, Language=go) = %q, want %q", got, "L")
		}
	})

	t.Run("language used when label empty", func(t *testing.T) {
		t.Parallel()

		if got := codeBlockLabel(CodeBlockProps{Language: "go"}); got != "go" {
			t.Errorf("codeBlockLabel(Language=go) = %q, want %q", got, "go")
		}
	})

	t.Run("both empty", func(t *testing.T) {
		t.Parallel()

		if got := codeBlockLabel(CodeBlockProps{}); got != "" {
			t.Errorf("codeBlockLabel() = %q, want empty", got)
		}
	})
}

func TestCodeBlockVariantIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v    CodeBlockVariant
		want bool
	}{
		{"block", CodeBlockBlock, true},
		{"compact id", CodeBlockCompactID, true},
		{"empty", CodeBlockVariant(""), false},
		{"bogus", CodeBlockVariant("bogus"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := CodeBlockVariantIsValid(tt.v); got != tt.want {
				t.Errorf("CodeBlockVariantIsValid(%q) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}
