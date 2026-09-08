package forms

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// formLayoutContextKey is the context key a Form uses to tell the fields it
// renders which layout they sit in (#166).
type formLayoutContextKey struct{}

// formInlineFields renders the Form's children with the Inline-layout marker
// in the context. Every FormFieldWrapper below then groups its
// label/control/error/help into ONE flex item instead of scattering them as
// separate flex children — without the grouping, each w-full control forces
// its own row and the Inline layout degenerates to a stack (the 2026-09-08
// demo audit's f25 finding; the demo had to hand-patch every field with
// Class: "sm:w-auto sm:min-w-40" to work around it).
//
// The ctx argument is the templ render context carrying the Form's children;
// they are re-rendered under the marker context unchanged otherwise.
func formInlineFields(ctx context.Context) templ.Component {
	children := templ.GetChildren(ctx)
	inlineCtx := context.WithValue(ctx, formLayoutContextKey{}, FormLayoutInline)

	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		return children.Render(inlineCtx, w)
	})
}

// formFieldInline reports whether FormFieldWrapper currently renders inside
// an Inline-layout Form (and should emit its grouping wrapper div).
func formFieldInline(ctx context.Context) bool {
	return ctx.Value(formLayoutContextKey{}) == FormLayoutInline
}
