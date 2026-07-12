package display

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestImageSrcSet(t *testing.T) {
	t.Parallel()
	html := utils.Render(t, Image(ImageProps{
		Src:    "/photo.jpg",
		Alt:    "A photo",
		SrcSet: "/photo-480w.jpg 480w, /photo-800w.jpg 800w",
		Sizes:  "(max-width: 600px) 480px, 800px",
	}))
	if !strings.Contains(html, `srcset="/photo-480w.jpg 480w, /photo-800w.jpg 800w"`) {
		t.Error("SrcSet should emit srcset attribute")
	}
	if !strings.Contains(html, `sizes="(max-width: 600px) 480px, 800px"`) {
		t.Error("Sizes should emit sizes attribute")
	}
}

func TestImageSrcSetEmptyOmitted(t *testing.T) {
	t.Parallel()
	html := utils.Render(t, Image(ImageProps{
		Src: "/photo.jpg",
		Alt: "A photo",
	}))
	if strings.Contains(html, "srcset") {
		t.Error("Empty SrcSet should not emit srcset attribute")
	}
	if strings.Contains(html, "sizes=") {
		t.Error("Empty Sizes should not emit sizes attribute")
	}
}

func TestTableLazyRows(t *testing.T) {
	t.Parallel()
	html := utils.Render(t, Table(TableProps{
		Headers:  []string{"Name", "Email"},
		Striped:  true,
		LazyRows: true,
		Rows: []TableRow{
			SimpleTableRow("Alice", "alice@example.com"),
			SimpleTableRow("Bob", "bob@example.com"),
		},
	}))
	if !strings.Contains(html, "tc-content-auto") {
		t.Error("LazyRows=true should add tc-content-auto class to rows")
	}
}

func TestTableLazyRowsFalseOmits(t *testing.T) {
	t.Parallel()
	html := utils.Render(t, Table(TableProps{
		Headers: []string{"Name", "Email"},
		Rows: []TableRow{
			SimpleTableRow("Alice", "alice@example.com"),
		},
	}))
	if strings.Contains(html, "tc-content-auto") {
		t.Error("LazyRows=false should not add tc-content-auto class")
	}
}

func TestTableLazyRowsCompactUsesCompactVariant(t *testing.T) {
	t.Parallel()
	html := utils.Render(t, Table(TableProps{
		Headers:     []string{"Name", "Email"},
		LazyRows:    true,
		CellPadding: TableCellPaddingCompact,
		Rows: []TableRow{
			SimpleTableRow("Alice", "alice@example.com"),
		},
	}))
	if strings.Contains(html, "tc-content-auto\"") {
		t.Error("LazyRows+Compact should NOT use tc-content-auto (should use tc-content-auto-compact)")
	}
	if !strings.Contains(html, "tc-content-auto-compact") {
		t.Error("LazyRows+Compact should use tc-content-auto-compact class")
	}
}
