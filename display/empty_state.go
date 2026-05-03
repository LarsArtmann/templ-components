package display

import (
	"github.com/larsartmann/templ-components/icons"
	"github.com/larsartmann/templ-components/utils"
)

//nolint:gochecknoglobals // Package-level lookup table for empty state icons
var emptyStateIconMap = map[string]icons.Name{
	"folder":   icons.Folder,
	"search":   icons.Search,
	"document": icons.Document,
	"inbox":    icons.Inbox,
	"chart":    icons.Chart,
	"users":    icons.Users,
}

func mapEmptyStateIcon(name string) icons.Name {
	return utils.MapEnum(emptyStateIconMap, icons.Inbox, name)
}
