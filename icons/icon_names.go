// Icon name constants — single source of truth for the icon catalogue.
package icons

import "sort"

// Name represents a named icon from the library
type Name string

// Icon name constants for use with the Icon component
const (
	Home                  Name = "home"
	Users                 Name = "users"
	Folder                Name = "folder"
	Document              Name = "document"
	Search                Name = "search"
	Settings              Name = "settings"
	Chart                 Name = "chart"
	Inbox                 Name = "inbox"
	Check                 Name = "check"
	X                     Name = "x"
	Close                 Name = "x"
	Plus                  Name = "plus"
	Minus                 Name = "minus"
	ChevronRight          Name = "chevron-right"
	ChevronLeft           Name = "chevron-left"
	ChevronDown           Name = "chevron-down"
	ChevronUp             Name = "chevron-up"
	ArrowRight            Name = "arrow-right"
	ArrowLeft             Name = "arrow-left"
	Refresh               Name = "refresh"
	ExternalLink          Name = "external-link"
	Download              Name = "download"
	Upload                Name = "upload"
	Trash                 Name = "trash"
	Edit                  Name = "edit"
	Eye                   Name = "eye"
	EyeOff                Name = "eye-off"
	Lock                  Name = "lock"
	Unlock                Name = "unlock"
	Menu                  Name = "menu"
	Bell                  Name = "bell"
	Calendar              Name = "calendar"
	Clock                 Name = "clock"
	Location              Name = "location"
	Phone                 Name = "phone"
	Mail                  Name = "mail"
	Globe                 Name = "globe"
	Sun                   Name = "sun"
	Moon                  Name = "moon"
	Spinner               Name = "spinner"
	ExclamationTriangle   Name = "exclamation-triangle"
	CheckCircle           Name = "check-circle"
	ExclamationCircle     Name = "exclamation-circle"
	Information           Name = "information"
	Question              Name = "question"
	ArrowUp               Name = "arrow-up"
	ArrowDown             Name = "arrow-down"
	Bookmark              Name = "bookmark"
	Clipboard             Name = "clipboard"
	Cloud                 Name = "cloud"
	CodeBracket           Name = "code-bracket"
	DocumentDuplicate     Name = "document-duplicate"
	DocumentText          Name = "document-text"
	EllipsisHorizontal    Name = "ellipsis-horizontal"
	EllipsisVertical      Name = "ellipsis-vertical"
	Filter                Name = "filter"
	Heart                 Name = "heart"
	Link                  Name = "link"
	ListBullet            Name = "list-bullet"
	MapPin                Name = "map-pin"
	Microphone            Name = "microphone"
	PaperAirplane         Name = "paper-airplane"
	Photo                 Name = "photo"
	Printer               Name = "printer"
	QueueList             Name = "queue-list"
	Share                 Name = "share"
	ShieldCheck           Name = "shield-check"
	Star                  Name = "star"
	Tag                   Name = "tag"
	ThumbUp               Name = "thumb-up"
	UserCircle            Name = "user-circle"
	UserPlus              Name = "user-plus"
	Wrench                Name = "wrench"
	XCircle               Name = "x-circle"
	ArchiveBox            Name = "archive-box"
	ArrowPath             Name = "arrow-path"
	Bars3                 Name = "bars-3"
	Beaker                Name = "beaker"
	Bolt                  Name = "bolt"
	BugAnt                Name = "bug-ant"
	Calculator            Name = "calculator"
	Cube                  Name = "cube"
	FaceSmile             Name = "face-smile"
	Fire                  Name = "fire"
	FolderOpen            Name = "folder-open"
	Gift                  Name = "gift"
	HandThumbUp           Name = "hand-thumb-up"
	Hashtag               Name = "hashtag"
	PuzzlePiece           Name = "puzzle-piece"
	RocketLaunch          Name = "rocket-launch"
	Server                Name = "server"
	Signal                Name = "signal"
	Squares2x2            Name = "squares-2x2"
	AcademicCap           Name = "academic-cap"
	ArrowDownOnSquare     Name = "arrow-down-on-square"
	ArrowUpOnSquare       Name = "arrow-up-on-square"
	BellSlash             Name = "bell-slash"
	Camera                Name = "camera"
	NoSymbol              Name = "no-symbol"
	ArrowRightOnRectangle Name = "arrow-right-on-rectangle"
	BuildingOffice2       Name = "building-office-2"
	Key                   Name = "key"
)

// allIconNames returns all icon names, auto-generated from iconPathData + Spinner.
// This eliminates the need to manually maintain a separate list.
func allIconNames() []Name {
	names := make([]Name, 0, len(iconPathData)+1)
	for name := range iconPathData {
		names = append(names, name)
	}
	names = append(names, Spinner)
	sort.Slice(names, func(i, j int) bool {
		return string(names[i]) < string(names[j])
	})
	return names
}

// AllIconNames returns all available icon names, sorted alphabetically.
// Useful for icon galleries, demos, and documentation.
func AllIconNames() []Name {
	return allIconNames()
}

// NameIsValid reports whether v is one of the defined icon names.
func NameIsValid(v Name) bool {
	if v == Spinner {
		return true
	}
	_, ok := iconPathData[v]
	return ok
}
