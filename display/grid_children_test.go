package display

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// hybridChild is a programmatically-constructed child component (the hybrid
// rendering path: no surrounding templ file).
func hybridChild(marker string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<p data-hybrid-child=\""+marker+"\">child "+marker+"</p>")

		return err
	})
}

// TestGridHybridChildrenViaWithChildren pins the documented hybrid-path
// escape hatch: children slots populate when the caller threads
// templ.WithChildren through the context (docs/recipes/hybrid-strings-builder-rendering.md).
func TestGridHybridChildrenViaWithChildren(t *testing.T) {
	t.Parallel()

	var b strings.Builder

	ctx := templ.WithChildren(context.Background(), hybridChild("a"))
	if err := Grid(GridProps{Cols: GridCols2}).Render(ctx, &b); err != nil {
		t.Fatalf("render: %v", err)
	}

	utils.AssertContains(t, b.String(), `data-hybrid-child="a"`)
}

// TestGridHybridChildrenEmptyWithoutWithChildren pins the caveat itself:
// rendered standalone into a strings.Builder WITHOUT templ.WithChildren,
// the { children... } slot renders empty (inheritance from a templ call is
// the only thing that populates it).
func TestGridHybridChildrenEmptyWithoutWithChildren(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	if err := Grid(GridProps{Cols: GridCols2}).Render(context.Background(), &b); err != nil {
		t.Fatalf("render: %v", err)
	}

	if strings.Contains(b.String(), "data-hybrid-child") {
		t.Fatalf("expected no children in plain-context render, got: %s", b.String())
	}
}
