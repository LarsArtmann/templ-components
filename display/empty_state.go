package display

import "github.com/larsartmann/templ-components/icons"

func mapEmptyStateIcon(name string) icons.Name {
	switch name {
	case "folder":
		return icons.Folder
	case "search":
		return icons.Search
	case "document":
		return icons.Document
	case "inbox":
		return icons.Inbox
	case "chart":
		return icons.Chart
	case "users":
		return icons.Users
	default:
		return icons.Inbox
	}
}
