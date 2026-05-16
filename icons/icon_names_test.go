package icons

import "testing"

var allIconNames = []Name{
	Home,
	Users,
	Folder,
	Document,
	Search,
	Settings,
	Chart,
	Inbox,
	Check,
	X,
	Plus,
	Minus,
	ChevronRight,
	ChevronLeft,
	ChevronDown,
	ChevronUp,
	ArrowRight,
	ArrowLeft,
	Refresh,
	ExternalLink,
	Download,
	Upload,
	Trash,
	Edit,
	Eye,
	EyeOff,
	Lock,
	Unlock,
	Menu,
	Bell,
	Calendar,
	Clock,
	Location,
	Phone,
	Mail,
	Globe,
	Sun,
	Moon,
	Spinner,
	Exclamation,
	ExclamationTriangle,
	CheckCircle,
	ExclamationCircle,
	Information,
	Question,
}

func TestIconNames(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		value Name
	}{
		{name: "Home", value: Home},
		{name: "Users", value: Users},
		{name: "Folder", value: Folder},
		{name: "Document", value: Document},
		{name: "Search", value: Search},
		{name: "Settings", value: Settings},
		{name: "Chart", value: Chart},
		{name: "Inbox", value: Inbox},
		{name: "Check", value: Check},
		{name: "X", value: X},
		{name: "Plus", value: Plus},
		{name: "Minus", value: Minus},
		{name: "ChevronRight", value: ChevronRight},
		{name: "ChevronLeft", value: ChevronLeft},
		{name: "ChevronDown", value: ChevronDown},
		{name: "ChevronUp", value: ChevronUp},
		{name: "ArrowRight", value: ArrowRight},
		{name: "ArrowLeft", value: ArrowLeft},
		{name: "Refresh", value: Refresh},
		{name: "ExternalLink", value: ExternalLink},
		{name: "Download", value: Download},
		{name: "Upload", value: Upload},
		{name: "Trash", value: Trash},
		{name: "Edit", value: Edit},
		{name: "Eye", value: Eye},
		{name: "EyeOff", value: EyeOff},
		{name: "Lock", value: Lock},
		{name: "Unlock", value: Unlock},
		{name: "Menu", value: Menu},
		{name: "Bell", value: Bell},
		{name: "Calendar", value: Calendar},
		{name: "Clock", value: Clock},
		{name: "Location", value: Location},
		{name: "Phone", value: Phone},
		{name: "Mail", value: Mail},
		{name: "Globe", value: Globe},
		{name: "Sun", value: Sun},
		{name: "Moon", value: Moon},
		{name: "Spinner", value: Spinner},
		{name: "Exclamation", value: Exclamation},
		{name: "ExclamationTriangle", value: ExclamationTriangle},
		{name: "CheckCircle", value: CheckCircle},
		{name: "ExclamationCircle", value: ExclamationCircle},
		{name: "Information", value: Information},
		{name: "Question", value: Question},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.value == "" {
				t.Errorf("icon %s should not be empty", tt.name)
			}
		})
	}
}

func TestIconCount(t *testing.T) {
	t.Parallel()
	if len(allIconNames) != 45 {
		t.Errorf("expected 45 icons, got %d", len(allIconNames))
	}
}
