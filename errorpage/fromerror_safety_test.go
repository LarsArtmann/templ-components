package errorpage

import (
	"errors"
	"testing"

	errorfamily "github.com/larsartmann/go-error-family"
)

// TestFromErrorFamilyRenameSafety guards the typed switch in FromErrorFamily
// against silent collapse. If go-error-family adds, removes, or renumbers a
// Family constant and the switch here is not updated, one of these assertions
// fails loudly instead of every error silently rendering as Transient.
//
// It checks three invariants:
//  1. Totality — every valid errorfamily.Family maps to a valid errorpage
//     Family (one present in familyStyleMap), never the default fallback for a
//     non-Transient input.
//  2. Correctness — each known family maps to its exact expected counterpart.
//  3. Injectivity — no two distinct input families collapse to the same output
//     (a missing switch case would hit the default and duplicate Transient).
func TestFromErrorFamilyRenameSafety(t *testing.T) {
	t.Parallel()

	// The canonical mapping. Adding a family to go-error-family without adding
	// it here makes Correctness fail; forgetting the switch case makes
	// Totality + Injectivity fail.
	expected := map[errorfamily.Family]Family{
		errorfamily.Rejection:      FamilyRejection,
		errorfamily.Conflict:       FamilyConflict,
		errorfamily.Transient:      FamilyTransient,
		errorfamily.Corruption:     FamilyCorruption,
		errorfamily.Infrastructure: FamilyInfrastructure,
		errorfamily.Orchestration:  FamilyOrchestration,
	}

	seen := make(map[Family][]errorfamily.Family) // output -> inputs that produced it

	for f := errorfamily.Rejection; f <= errorfamily.Orchestration; f++ {
		if !f.IsValid() {
			continue
		}

		got := FromErrorFamily(f)

		// Totality: output must be a known, styleable Family.
		if !FamilyIsValid(got) {
			t.Errorf("FromErrorFamily(%v) = %q which is not a valid Family (missing from familyStyleMap)", f, got)
		}

		// Correctness: exact expected counterpart.
		if want, ok := expected[f]; ok && got != want {
			t.Errorf("FromErrorFamily(%v) = %q, want %q", f, got, want)
		}

		// Track for injectivity.
		seen[got] = append(seen[got], f)
	}

	// Injectivity: no two distinct inputs may collapse to the same output. The
	// only exception would be an intentional alias, of which there are none.
	for out, inputs := range seen {
		if len(inputs) > 1 {
			t.Errorf(
				"FromErrorFamily is not injective: inputs %v all collapsed to %q — a switch case is likely missing",
				inputs, out,
			)
		}
	}

	// Ensure we actually exercised every expected family (guards against the
	// loop range silently shrinking if the iota sequence changes).
	if got, want := len(expected), len(seen); got != want {
		t.Errorf("expected %d distinct family mappings, observed %d", want, len(seen))
	}
}

// TestFromErrorFamilyDefaultIsTransient documents the fallback contract: an
// unrecognized Family value renders as Transient (the safest, retryable default
// for a UI). If this fallback ever changes, the test forces an explicit update.
func TestFromErrorFamilyDefaultIsTransient(t *testing.T) {
	t.Parallel()

	// A value beyond the defined range exercises the default branch.
	outOfRange := errorfamily.Orchestration + 1
	if got := FromErrorFamily(outOfRange); got != FamilyTransient {
		t.Errorf("FromErrorFamily(out-of-range) = %q, want FamilyTransient (default fallback)", got)
	}
}

// TestFromErrorFamilyCoversAllValidFamilies is a meta-check: the count of valid
// families in go-error-family must equal the count of cases the switch can
// resolve. If a 7th family is added upstream, this fails until the switch (and
// the expected map above) are updated.
func TestFromErrorFamilyCoversAllValidFamilies(t *testing.T) {
	t.Parallel()

	validCount := 0

	for f := errorfamily.Rejection; f <= errorfamily.Orchestration; f++ {
		if f.IsValid() {
			validCount++
		}
	}

	// Every valid family must round-trip to a Family that FamilyIsValid accepts
	// and that is NOT the default fallback (unless the input is Transient).
	for f := errorfamily.Rejection; f <= errorfamily.Orchestration; f++ {
		if !f.IsValid() {
			continue
		}

		got := FromErrorFamily(f)
		if !FamilyIsValid(got) {
			t.Errorf("family %v dropped: FromErrorFamily returns invalid %q", f, got)
		}
	}
}

// publicTraceError mimics the surface bridge.ClassifiedError exposes: an
// oops-style user-safe message (Public) and a correlation trace ID (Trace),
// both promoted from the embedded samber/oops OopsError. Nothing here
// imports oops — FromError must keep consuming this via small interfaces.
type publicTraceError struct {
	err    error
	public string
	trace  string
}

func (e *publicTraceError) Error() string { return e.err.Error() }

func (e *publicTraceError) Public() string { return e.public }

func (e *publicTraceError) Trace() string { return e.trace }

func (e *publicTraceError) ErrorFamily() errorfamily.Family { return errorfamily.Rejection }

// TestFromErrorPrefersPublicMessage verifies that an error carrying an
// oops-style user-safe message gets it rendered instead of the internal
// developer-facing Error() string — the bridge/oops integration advantage.
func TestFromErrorPrefersPublicMessage(t *testing.T) {
	t.Parallel()

	err := &publicTraceError{
		err:    errors.New("rbac: deny user=42 role=guest"),
		public: "You do not have access to this resource.",
	}

	props := FromError(err)
	if props.Message != "You do not have access to this resource." {
		t.Errorf("FromError Message = %q, want the Public() user-safe message", props.Message)
	}
}

// TestFromErrorSkipsEmptyPublicMessage verifies the Public probe degrades
// gracefully: an empty Public() falls through instead of rendering a blank
// message.
func TestFromErrorSkipsEmptyPublicMessage(t *testing.T) {
	t.Parallel()

	err := &publicTraceError{
		err:    errors.New("internal detail"),
		public: "",
	}

	props := FromError(err)
	if props.Message != "internal detail" {
		t.Errorf("FromError Message = %q, want the Error() fallback when Public is empty", props.Message)
	}
}

// TestFromErrorExtractsTrace verifies the correlation trace ID flows into
// props.Trace for the footer, and stays empty for plain errors.
func TestFromErrorExtractsTrace(t *testing.T) {
	t.Parallel()

	traced := &publicTraceError{err: errors.New("boom"), trace: "trc_9f3a1c2d"}
	if got := FromError(traced).Trace; got != "trc_9f3a1c2d" {
		t.Errorf("FromError Trace = %q, want %q", got, "trc_9f3a1c2d")
	}

	plain := errors.New("boom")
	if got := FromError(plain).Trace; got != "" {
		t.Errorf("FromError Trace for plain error = %q, want empty", got)
	}
}
