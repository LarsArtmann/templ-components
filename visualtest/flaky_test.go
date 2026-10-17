package visualtest

import (
	"net/http/httptest"
	"testing"
	"time"
)

// TestFlakyFail pins the doomed-interaction branch of the shared flaky
// contract: the fail ID answers 500 and stops the caller; every other ID
// passes through untouched.
func TestFlakyFail(t *testing.T) {
	t.Parallel()

	t.Run("fail ID writes 500 and stops", func(t *testing.T) {
		t.Parallel()
		rec := httptest.NewRecorder()

		if !FlakyFail(rec, FlakyFailID) {
			t.Fatal("FlakyFail(FlakyFailID) must report handled")
		}

		if rec.Code != 500 {
			t.Errorf("fail branch status = %d, want 500", rec.Code)
		}
	})

	t.Run("slow ID passes through", func(t *testing.T) {
		t.Parallel()
		rec := httptest.NewRecorder()

		if FlakyFail(rec, FlakySlowID) {
			t.Fatal("FlakyFail(FlakySlowID) must report unhandled — the stall is the caller's shape")
		}

		if rec.Code != 200 {
			t.Errorf("pass-through status = %d, want untouched 200", rec.Code)
		}
	})

	t.Run("slow delay is observable, not instant", func(t *testing.T) {
		t.Parallel()
		if FlakySlowDelay <= 0 {
			t.Fatalf("FlakySlowDelay = %v, want > 0", FlakySlowDelay)
		}
	})
}
