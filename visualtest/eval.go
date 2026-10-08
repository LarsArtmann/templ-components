package visualtest

import (
	"context"

	"github.com/chromedp/chromedp"
)

// evalInto adapts a typed chromedp action (v0.20's Action[T] returns its
// result) to a Void action that stores the result into dest, so evaluation
// steps keep their position inside a single chromedp.Do chain instead of
// being split into a separate Run call (order matters: evaluate after
// navigate/wait).
func evalInto[T any](action chromedp.Action[T], dest *T) chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, t *chromedp.Target) error {
		v, err := action(ctx, t)

		*dest = v

		return err
	})
}
