package layout_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/layout"
)

func ExampleScript() {
	var buf bytes.Buffer
	_ = layout.Script("nonce-abc", "/static/app.js", nil).Render(context.Background(), &buf)
	fmt.Println(buf.String())
	// Output will contain a CSP-safe script tag with nonce
}
