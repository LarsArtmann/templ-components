package display

import "github.com/larsartmann/templ-components/utils/tone"

// Tone/Hue adapters (utils/tone Phase 0, TODO #425): every color-bearing
// enum in display gains a mapping onto the canonical color language. The
// existing lookup maps keep rendering (zero visual change); equivalence
// tests in tone_adapters_test.go prove the adapters describe the shipped
// strings byte-for-byte. The maps collapse into the tone tables at v2
// (ADR-0039/0044, TODO #431).

// Tone maps StatusTone onto the canonical semantic axis.
func (t StatusTone) Tone() tone.Tone {
	switch t {
	case StatusToneSuccess:
		return tone.ToneSuccess
	case StatusToneWarning:
		return tone.ToneWarning
	case StatusToneDanger:
		return tone.ToneDanger
	case StatusToneInfo:
		return tone.ToneInfo
	default:
		return tone.ToneNeutral
	}
}

// DotHue is the hue StatusDot actually renders for t. Warning deliberately
// diverges from the tone table (amber here, yellow in the tone table's
// majority dialect) — the documented warning split brain, unified by the
// #428 decision.
func (t StatusTone) DotHue() tone.Hue {
	switch t {
	case StatusToneSuccess:
		return tone.HueGreen
	case StatusToneWarning:
		return tone.HueAmber
	case StatusToneDanger:
		return tone.HueRed
	case StatusToneInfo:
		return tone.HueBlue
	default:
		return tone.HueGray
	}
}

// Hue maps KanbanTone onto the categorical hue axis. Kanban lanes are
// categorical, not semantic — there is deliberately no Tone() here.
func (t KanbanTone) Hue() tone.Hue {
	switch t {
	case KanbanToneBlue:
		return tone.HueBlue
	case KanbanToneGreen:
		return tone.HueGreen
	case KanbanToneYellow:
		return tone.HueYellow
	case KanbanToneRed:
		return tone.HueRed
	case KanbanTonePurple:
		return tone.HuePurple
	default:
		return tone.HueGray
	}
}

// Hue maps StatTone onto the categorical hue axis (stat tiles are
// decorative emphasis, not semantics).
func (t StatTone) Hue() tone.Hue {
	switch t {
	case StatToneGreen:
		return tone.HueGreen
	case StatToneYellow:
		return tone.HueYellow
	case StatToneRed:
		return tone.HueRed
	case StatTonePurple:
		return tone.HuePurple
	case StatToneBlue:
		return tone.HueBlue
	default:
		return tone.HueGray
	}
}

// Tone maps TrendDirection onto the canonical semantic axis.
func (t TrendDirection) Tone() tone.Tone {
	switch t {
	case TrendUp:
		return tone.ToneSuccess
	case TrendDown:
		return tone.ToneDanger
	case TrendWarn:
		return tone.ToneWarning
	default:
		return tone.ToneNeutral
	}
}

// Tone maps ScrollbackTone onto the canonical semantic axis.
func (t ScrollbackTone) Tone() tone.Tone {
	switch t {
	case ScrollbackToneInfo:
		return tone.ToneInfo
	case ScrollbackToneSuccess:
		return tone.ToneSuccess
	case ScrollbackToneWarning:
		return tone.ToneWarning
	case ScrollbackToneDanger:
		return tone.ToneDanger
	default:
		return tone.ToneNeutral
	}
}

// Tone maps BadgeType onto the canonical semantic axis. Primary maps to
// Accent (the ADR-0044 accent slot); the FeedbackError-style "error"
// spelling maps to Danger.
func (t BadgeType) Tone() tone.Tone {
	switch t {
	case BadgeSuccess:
		return tone.ToneSuccess
	case BadgeWarning:
		return tone.ToneWarning
	case BadgeError:
		return tone.ToneDanger
	case BadgeInfo:
		return tone.ToneInfo
	case BadgePrimary:
		return tone.ToneAccent
	default:
		return tone.ToneNeutral
	}
}
