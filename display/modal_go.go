package display

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// ModalProps configures a modal dialog
type ModalProps struct {
	utils.BaseProps
	Title string
	Open  bool
	Size  string // "sm", "md", "lg", "xl", "full"
	Nonce string
}

// DefaultModalProps returns sensible defaults
func DefaultModalProps() ModalProps {
	return ModalProps{
		Size: "md",
	}
}

func modalSizeClass(size string) string {
	switch size {
	case "sm":
		return "max-w-sm"
	case "lg":
		return "max-w-lg"
	case "xl":
		return "max-w-xl"
	case "full":
		return "max-w-4xl"
	default:
		return "max-w-md"
	}
}

func modalCloseHandler(id string) templ.ComponentScript {
	return templ.ComponentScript{
		Name:     fmt.Sprintf("tcCloseModal_%s", id),
		Function: fmt.Sprintf(`function tcCloseModal_%s(id){tcCloseModal(id)}`, id),
		Call:     fmt.Sprintf(`tcCloseModal('%s')`, id),
	}
}
