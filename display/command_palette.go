// CommandPalette component: a CSP-safe <dialog>-based command palette with
// client-side filtering over server-rendered groups, selection via links or
// wire actions.
package display

import (
	"strconv"
	"strings"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
)

// CommandPaletteSize defines the width of the palette panel.
type CommandPaletteSize string

// Command palette size constants.
const (
	CommandPaletteSizeSM CommandPaletteSize = "sm"
	CommandPaletteSizeMD CommandPaletteSize = "md"
	CommandPaletteSizeLG CommandPaletteSize = "lg"
)

// CommandItem is one selectable row of the palette. Exactly one of Href
// (link navigation) or Wire (HTMX/Datastar action) should be set; with both,
// the link wins and Wire is ignored. Items with neither render as inert rows.
type CommandItem struct {
	Label string
	Hint  string // secondary text rendered under or beside the label
	Href  string // link selection
	Wire  *wire.Action
}

// CommandGroup is a labeled section of palette items. Groups whose items all
// fail the filter are hidden entirely while filtering.
type CommandGroup struct {
	Label string
	Items []CommandItem
}

// CommandPaletteProps configures a command palette.
type CommandPaletteProps struct {
	utils.BaseProps

	Placeholder  string
	Groups       []CommandGroup
	TriggerLabel string // renders an opener button when non-empty
	Hotkey       bool   // register Cmd/Ctrl+K to open
	EmptyMessage string
	Size         CommandPaletteSize
}

// DefaultCommandPaletteProps returns sensible defaults.
func DefaultCommandPaletteProps() CommandPaletteProps {
	return CommandPaletteProps{ //nolint:exhaustruct_v5 // intentionally minimal defaults
		Placeholder:  "Type a command or search…",
		EmptyMessage: "No results found.",
		Size:         CommandPaletteSizeMD,
	}
}

//nolint:gochecknoglobals // Package-level lookup table; Tailwind class strings are intentionally inline
var commandPaletteSizeLookup = map[CommandPaletteSize]string{
	CommandPaletteSizeSM: "max-w-sm",
	CommandPaletteSizeMD: "max-w-lg",
	CommandPaletteSizeLG: "max-w-2xl",
}

func commandPaletteSizeClass(size CommandPaletteSize) string {
	return utils.Lookup(commandPaletteSizeLookup, size, commandPaletteSizeLookup[CommandPaletteSizeMD])
}

// CommandPaletteSizeIsValid reports whether s is one of the defined
// CommandPaletteSize constants.
func CommandPaletteSizeIsValid(s CommandPaletteSize) bool {
	_, ok := commandPaletteSizeLookup[s]

	return ok
}

// commandPaletteJS returns the singleton JavaScript for the palette: open
// triggers (button + optional Cmd/Ctrl+K), client-side filtering with group
// collapsing and an empty state, and arrow-key/Enter navigation over the
// visible options. Filtering is plain substring matching on the item text;
// all highlighting toggles data attributes styled by compiled Tailwind
// variants so no class names are assembled at runtime.
func commandPaletteJS(id string) string {
	guard := "window.tcPaletteAttached"
	escapedID := strconv.Quote(id)

	return "if(!" + guard + "){" + guard + "=true;" +
		"document.addEventListener('keydown',function(e){" +
		"if((e.ctrlKey||e.metaKey)&&e.key.toLowerCase()==='k'){e.preventDefault();" +
		"document.querySelectorAll('[data-tc-palette][data-tc-hotkey]').forEach(function(d){if(!d.open)d.showModal();});}}" +
		",true);" +
		"}" +
		"(function(id){" +
		"var d=document.getElementById(id);if(!d)return;" +
		"var input=d.querySelector('[data-tc-palette-input]');" +
		"var empty=d.querySelector('[data-tc-palette-empty]');" +
		"var active=null;" +
		"function options(){return Array.prototype.slice.call(d.querySelectorAll('[data-tc-palette-item]'))" +
		".filter(function(o){return o.getAttribute('data-tc-palette-hidden')!=='true';});}" +
		"function groups(){return d.querySelectorAll('[data-tc-palette-group]');}" +
		"function setActive(opt){if(active)active.removeAttribute('data-tc-palette-active');" +
		"active=opt;if(active){active.setAttribute('data-tc-palette-active','true');" +
		"input.setAttribute('aria-activedescendant',active.id);" +
		"active.scrollIntoView({block:'nearest'});}else{input.removeAttribute('aria-activedescendant');}}" +
		"function applyFilter(){" +
		"var q=(input.value||'').trim().toLowerCase();var any=false;" +
		"groups().forEach(function(g){" +
		"var gAny=false;" +
		"g.querySelectorAll('[data-tc-palette-item]').forEach(function(o){" +
		"var hit=!q||o.textContent.toLowerCase().indexOf(q)!==-1;" +
		"o.setAttribute('data-tc-palette-hidden',hit?'false':'true');if(hit)gAny=true;});" +
		"g.setAttribute('data-tc-palette-hidden',gAny?'false':'true');if(gAny)any=true;});" +
		"if(empty)empty.setAttribute('data-tc-palette-hidden',any?'true':'false');" +
		"setActive(options()[0]||null);}" +
		"function open(){" +
		"if(input){input.value='';}applyFilter();if(d&&!d.open)d.showModal();" +
		"if(input)input.focus();}" +
		"d.addEventListener('click',function(e){" +
		"if(e.target===d){d.close();return;}" +
		"var opt=e.target.closest('[data-tc-palette-item]');" +
		"if(opt){d.close();return;}" +
		"if(e.target.closest('[data-tc-close]'))d.close();});" +
		"if(input){input.addEventListener('input',applyFilter);" +
		"input.addEventListener('keydown',function(e){" +
		"var opts=options();" +
		"if(e.key==='ArrowDown'||e.key==='ArrowUp'){e.preventDefault();" +
		"if(!opts.length)return;var i=opts.indexOf(active);" +
		"setActive(opts[e.key==='ArrowDown'?(i+1)%opts.length:(i-1+opts.length)%opts.length]);}" +
		"else if(e.key==='Enter'){if(active){e.preventDefault();active.click();}}});}" +
		"document.querySelectorAll('[data-tc-palette-open]').forEach(function(b){" +
		"b.addEventListener('click',function(){open();});});" +
		"d.addEventListener('open-restore',function(){open();});" +
		"})(" + escapedID + ");\n"
}

// commandPaletteItemID derives a DOM-safe option id for the item at index i.
func commandPaletteItemID(paletteID string, i int) string {
	return paletteID + "-item-" + strconv.Itoa(i)
}

// commandPaletteItemIsLink reports whether the item renders as a link row.
func commandPaletteItemIsLink(item CommandItem) bool {
	return strings.TrimSpace(item.Href) != ""
}
