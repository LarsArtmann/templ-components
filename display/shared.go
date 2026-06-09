package display

import (
	"strconv"

	"github.com/a-h/templ"
)

func closeHandler(componentName, id string) templ.ComponentScript {
	name := "tcClose" + componentName + "_" + id
	fn := "function tcClose" + componentName + "_" + id + "(id){tcClose" + componentName + "(id)}"
	escapedID := strconv.Quote(id)
	call := "tcClose" + componentName + "(" + escapedID + ")"
	return templ.ComponentScript{
		Name:       name,
		Function:   fn,
		Call:       call,
		CallInline: "",
	}
}

func safeID(id string) string {
	return strconv.Quote(id)
}
