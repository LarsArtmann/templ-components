package errorpage

import (
	"errors"
	"testing"
)

func TestErrorPagePropsValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid props return nil", func(t *testing.T) {
		t.Parallel()

		p := ErrorPageProps{Family: FamilyRejection, Title: "Not found"}

		if err := p.Validate(); err != nil {
			t.Errorf("Validate() unexpected error: %v", err)
		}
	})

	t.Run("invalid family returns error", func(t *testing.T) {
		t.Parallel()

		p := ErrorPageProps{Family: Family("bogus"), Title: "X"}

		if err := p.Validate(); err == nil {
			t.Fatal("Validate() expected error for invalid family, got nil")
		} else if !errors.Is(err, errValidateFamily) {
			t.Errorf("Validate() error = %v, want wrap of errValidateFamily", err)
		}
	})

	t.Run("status code below 400 returns error", func(t *testing.T) {
		t.Parallel()

		p := ErrorPageProps{Family: FamilyTransient, Title: "X", StatusCode: 200}
		err := p.Validate()

		if err == nil || !errors.Is(err, errValidateStatusRange) {
			t.Errorf("Validate() with StatusCode=200 error = %v, want errValidateStatusRange", err)
		}
	})

	t.Run("status code above 599 returns error", func(t *testing.T) {
		t.Parallel()

		p := ErrorPageProps{Family: FamilyTransient, Title: "X", StatusCode: 700}
		err := p.Validate()

		if err == nil || !errors.Is(err, errValidateStatusRange) {
			t.Errorf("Validate() with StatusCode=700 error = %v, want errValidateStatusRange", err)
		}
	})

	t.Run("status code 0 is accepted (unset)", func(t *testing.T) {
		t.Parallel()

		p := ErrorPageProps{Family: FamilyRejection, Title: "X"}

		if err := p.Validate(); err != nil {
			t.Errorf("Validate() with StatusCode=0 unexpected error: %v", err)
		}
	})

	t.Run("blank page (no title/message/cause) returns error", func(t *testing.T) {
		t.Parallel()

		p := ErrorPageProps{Family: FamilyRejection}
		err := p.Validate()

		if err == nil || !errors.Is(err, errValidateBlank) {
			t.Errorf("Validate() blank page error = %v, want errValidateBlank", err)
		}
	})

	t.Run("message alone satisfies blank check", func(t *testing.T) {
		t.Parallel()

		p := ErrorPageProps{Family: FamilyRejection, Message: "Something went wrong"}

		if err := p.Validate(); err != nil {
			t.Errorf("Validate() with Message unexpected error: %v", err)
		}
	})

	t.Run("cause chain alone satisfies blank check", func(t *testing.T) {
		t.Parallel()

		p := ErrorPageProps{
			Family:     FamilyRejection,
			CauseChain: []CauseItem{{Message: "root cause"}},
		}

		if err := p.Validate(); err != nil {
			t.Errorf("Validate() with CauseChain unexpected error: %v", err)
		}
	})
}
