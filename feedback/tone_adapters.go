package feedback

import "github.com/larsartmann/templ-components/utils/tone"

// Tone maps FeedbackType onto the canonical semantic axis (utils/tone
// Phase 0, TODO #425). The "error" spelling maps to Danger — the canonical
// spelling this library converges on at v2 (ADR-0022 removes the alias;
// the value rename rides the #431 batch).
func (t FeedbackType) Tone() tone.Tone {
	switch t {
	case FeedbackSuccess:
		return tone.ToneSuccess
	case FeedbackError:
		return tone.ToneDanger
	case FeedbackWarning:
		return tone.ToneWarning
	case FeedbackInfo:
		return tone.ToneInfo
	default:
		return tone.ToneNeutral
	}
}
