package forms

import (
	"bytes"
	"context"
	"testing"
)

func BenchmarkForms(b *testing.B) {
	b.Run("Input render", func(b *testing.B) {
		props := InputProps{
			ID:          "email",
			Name:        "email",
			Type:        InputEmail,
			Label:       "Email address",
			Placeholder: "you@example.com",
		}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Input(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Select render", func(b *testing.B) {
		props := SelectProps{
			ID:    "country",
			Name:  "country",
			Label: "Country",
			Options: []SelectOption{
				{Value: "de", Label: "Germany"},
				{Value: "at", Label: "Austria"},
				{Value: "ch", Label: "Switzerland"},
				{Value: "us", Label: "United States"},
				{Value: "uk", Label: "United Kingdom"},
			},
		}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Select(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Textarea render", func(b *testing.B) {
		props := TextareaProps{
			ID:    "bio",
			Name:  "bio",
			Label: "Bio",
			Rows:  4,
		}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Textarea(props).Render(context.Background(), &buf)
		}
	})

	b.Run("Combobox render", func(b *testing.B) {
		opts := make([]ComboboxOption, 10)
		for i := range opts {
			opts[i] = ComboboxOption{Value: "opt-" + string(rune('a'+i)), Label: "Option " + string(rune('A'+i))}
		}

		props := ComboboxProps{
			ID:      "combo",
			Name:    "choice",
			Label:   "Choose",
			Options: opts,
		}

		b.ResetTimer()

		for b.Loop() {
			var buf bytes.Buffer

			_ = Combobox(props).Render(context.Background(), &buf)
		}
	})
}
