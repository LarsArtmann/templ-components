package errorpage

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestErrorPageMainLandmark verifies ErrorPage uses <main> (not <div role="region">)
// for WCAG 2.4.1 Bypass Blocks compliance.
func TestErrorPageMainLandmark(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ErrorPage(ErrorPageProps{
		Family:  FamilyRejection,
		Message: "test error",
	}))
	utils.AssertContains(t, output, "<main")
	utils.AssertNotContains(t, output, `role="region"`)
}

// TestNotFound404MainLandmark verifies NotFound404 uses <main>.
func TestNotFound404MainLandmark(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, NotFound404(DefaultNotFound404Props()))
	utils.AssertContains(t, output, "<main")
	utils.AssertNotContains(t, output, `role="region"`)
}

// TestErrorPageHasNoRoleRegion ensures no regression to role="region".
func TestErrorPageHasNoRoleRegion(t *testing.T) {
	t.Parallel()

	for _, family := range []Family{FamilyRejection, FamilyConflict, FamilyTransient, FamilyCorruption, FamilyInfrastructure} {
		output := utils.Render(t, ErrorPage(ErrorPageProps{
			Family:  family,
			Message: "test",
		}))
		if strings.Contains(output, `role="region"`) {
			t.Errorf("family %q: found role=region in output", family)
		}
	}
}

// TestFromErrorFallbackFamily verifies unknown errors return FamilyCorruption (→500).
func TestFromErrorFallbackFamily(t *testing.T) {
	t.Parallel()

	props := FromError(&testError{msg: "unknown error"})
	if props.Family != FamilyCorruption {
		t.Errorf("Family = %q, want %q", props.Family, FamilyCorruption)
	}
}

// TestErrorDetailHasRoleAlert verifies the inline error card uses role="alert".
func TestErrorDetailHasRoleAlert(t *testing.T) {
	t.Parallel()
	output := utils.Render(t, ErrorDetail(ErrorDetailProps{
		Family:  FamilyCorruption,
		Message: "detail test",
	}))
	utils.AssertContains(t, output, `role="alert"`)
}
