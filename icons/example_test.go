package icons_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/icons"
)

func ExampleIcon() {
	var buf bytes.Buffer
	_ = icons.Icon(icons.Check, "h-5 w-5 text-green-500").Render(context.Background(), &buf)
	fmt.Println(buf.String())
}
