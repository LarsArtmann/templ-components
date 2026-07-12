// Radio and RadioGroup component types and helpers.
package forms

import (
	"github.com/a-h/templ"
	"github.com/larsartmann/templ-components/utils"
)

// RadioProps configures a single radio button input
type RadioProps struct {
	utils.BaseProps
	Name     string
	Value    string
	Label    string
	Checked  bool
	Required bool
	Disabled bool
}

// DefaultRadioProps returns sensible defaults
func DefaultRadioProps() RadioProps {
	return RadioProps{} //nolint:exhaustruct // intentionally minimal defaults
}

// RadioOption represents a single option in a radio group
type RadioOption struct {
	Value    string
	Label    string
	Checked  bool
	Disabled bool
}

// RadioGroupProps configures a group of radio buttons
type RadioGroupProps struct {
	utils.BaseProps
	Name     string
	Label    string // renders as <legend>
	Options  []RadioOption
	Inline   bool   // horizontal layout
	Error    string // group-level error
	HelpText string
	Required bool
}

// DefaultRadioGroupProps returns sensible defaults
func DefaultRadioGroupProps() RadioGroupProps {
	return RadioGroupProps{} //nolint:exhaustruct // intentionally minimal defaults
}

// radioItemProps builds a RadioProps for an option within a group.
// Auto-generates the ID from the group ID and option value.
func radioItemProps(groupID, name string, opt RadioOption, required bool, ariaAttrs templ.Attributes) RadioProps {
	p := RadioProps{ //nolint:exhaustruct // BaseProps.ID set conditionally below
		Name:     name,
		Value:    opt.Value,
		Label:    opt.Label,
		Checked:  opt.Checked,
		Required: required,
		Disabled: opt.Disabled,
	}
	if groupID != "" {
		p.ID = groupID + "-" + SanitizeID(opt.Value)
	}
	if len(ariaAttrs) > 0 {
		p.Attrs = ariaAttrs
	}
	return p
}
