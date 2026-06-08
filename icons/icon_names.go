// Package icons provides constants for predefined icon names used throughout the component library.
package icons

import "sort"

// Name represents a named icon from the library
type Name string

// Icon name constants for use with the Icon component
const (
	Home                Name = "home"
	Users               Name = "users"
	Folder              Name = "folder"
	Document            Name = "document"
	Search              Name = "search"
	Settings            Name = "settings"
	Chart               Name = "chart"
	Inbox               Name = "inbox"
	Check               Name = "check"
	X                   Name = "x"
	Plus                Name = "plus"
	Minus               Name = "minus"
	ChevronRight        Name = "chevron-right"
	ChevronLeft         Name = "chevron-left"
	ChevronDown         Name = "chevron-down"
	ChevronUp           Name = "chevron-up"
	ArrowRight          Name = "arrow-right"
	ArrowLeft           Name = "arrow-left"
	Refresh             Name = "refresh"
	ExternalLink        Name = "external-link"
	Download            Name = "download"
	Upload              Name = "upload"
	Trash               Name = "trash"
	Edit                Name = "edit"
	Eye                 Name = "eye"
	EyeOff              Name = "eye-off"
	Lock                Name = "lock"
	Unlock              Name = "unlock"
	Menu                Name = "menu"
	Bell                Name = "bell"
	Calendar            Name = "calendar"
	Clock               Name = "clock"
	Location            Name = "location"
	Phone               Name = "phone"
	Mail                Name = "mail"
	Globe               Name = "globe"
	Sun                 Name = "sun"
	Moon                Name = "moon"
	Spinner             Name = "spinner"
	ExclamationTriangle Name = "exclamation-triangle"
	CheckCircle         Name = "check-circle"
	ExclamationCircle   Name = "exclamation-circle"
	Information         Name = "information"
	Question            Name = "question"
	ArrowUp             Name = "arrow-up"
	ArrowDown           Name = "arrow-down"
	Bookmark            Name = "bookmark"
	Clipboard           Name = "clipboard"
	Cloud               Name = "cloud"
	CodeBracket         Name = "code-bracket"
	DocumentDuplicate   Name = "document-duplicate"
	DocumentText        Name = "document-text"
	EllipsisHorizontal  Name = "ellipsis-horizontal"
	EllipsisVertical    Name = "ellipsis-vertical"
	Filter              Name = "filter"
	Heart               Name = "heart"
	Link                Name = "link"
	ListBullet          Name = "list-bullet"
	MapPin              Name = "map-pin"
	Microphone          Name = "microphone"
	PaperAirplane       Name = "paper-airplane"
	Photo               Name = "photo"
	Printer             Name = "printer"
	QueueList           Name = "queue-list"
	Share               Name = "share"
	ShieldCheck         Name = "shield-check"
	Star                Name = "star"
	Tag                 Name = "tag"
	ThumbUp             Name = "thumb-up"
	UserCircle          Name = "user-circle"
	UserPlus            Name = "user-plus"
	Wrench              Name = "wrench"
	XCircle             Name = "x-circle"
)

// allIconNames returns all icon names, auto-generated from iconPathData + Spinner.
// This eliminates the need to manually maintain a separate list.
func allIconNames() []Name {
	names := make([]Name, 0, len(iconPathData)+1)
	for name := range iconPathData {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return string(names[i]) < string(names[j])
	})
	names = append(names, Spinner)
	return names
}
