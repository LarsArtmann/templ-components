package wire_test

import (
	"fmt"

	"github.com/larsartmann/templ-components/utils/wire"
)

func ExampleAction_Attributes() {
	attrs := wire.Action{Method: wire.MethodPost, URL: "/api/items", Target: "#list"}.Attributes()
	fmt.Println(attrs["hx-post"])
	fmt.Println(attrs["hx-target"])
	// Output:
	// /api/items
	// #list
}
