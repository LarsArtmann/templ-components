package display

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/golden"
)

// Golden sweep for CodeBlock (TODO_LIST #389): block and compact-id shapes,
// header variants, and copy suppression.

func TestGoldenSweepCodeBlock(t *testing.T) {
	t.Parallel()

	golden.AssertSnapshots(t, []golden.Snapshot{
		{Name: "code_block_default", HTML: utils.Render(t, CodeBlock(CodeBlockProps{
			Code:     "go install github.com/larsartmann/templ-components/cmd/tc@latest",
			Language: "bash",
		}))},
		{Name: "code_block_label_override", HTML: utils.Render(t, CodeBlock(CodeBlockProps{
			Code:  "SELECT 1;",
			Label: "Query",
		}))},
		{Name: "code_block_no_copy", HTML: utils.Render(t, CodeBlock(CodeBlockProps{
			Code:     "plain snippet",
			NoCopy:   true,
			Language: "text",
		}))},
		{Name: "code_block_no_header", HTML: utils.Render(t, CodeBlock(CodeBlockProps{
			Code:   "no header at all",
			NoCopy: true,
		}))},
		{Name: "code_block_compact_id", HTML: utils.Render(t, CodeBlock(CodeBlockProps{
			Code:    "9007199254740993",
			Variant: CodeBlockCompactID,
		}))},
		{Name: "code_block_compact_id_no_copy", HTML: utils.Render(t, CodeBlock(CodeBlockProps{
			Code:    "9007199254740993",
			Variant: CodeBlockCompactID,
			NoCopy:  true,
		}))},
		{Name: "code_block_custom_class", HTML: utils.Render(t, CodeBlock(CodeBlockProps{
			Class: "max-w-md",
			Code:  "bounded",
		}))},
	})
}
