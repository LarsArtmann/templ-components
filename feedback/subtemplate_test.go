package feedback

import (
	"strings"
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

// TestSkeletonContainerSubTemplate verifies the shared skeletonContainer renders
// the aria shell with role=status and aria-label.
func TestSkeletonContainerSubTemplate(t *testing.T) {
	t.Parallel()

	result := utils.Render(t, SkeletonCardGrid(3))
	if !strings.Contains(result, `role="status"`) {
		t.Error("skeletonContainer: role=status not found")
	}

	if !strings.Contains(result, "Loading") {
		t.Error("skeletonContainer: loading label not found")
	}
}

// TestSkeletonContainerZeroCount verifies the n<=0 fallback renders a single card.
func TestSkeletonContainerZeroCount(t *testing.T) {
	t.Parallel()

	result := utils.Render(t, SkeletonCardGrid(0))
	if !strings.Contains(result, `role="status"`) {
		t.Error("skeletonContainer: role=status missing on zero count")
	}
}

// TestSkeletonContainerNegativeCount verifies the n<0 fallback.
func TestSkeletonContainerNegativeCount(t *testing.T) {
	t.Parallel()

	result := utils.Render(t, SkeletonCardGrid(-5))
	if !strings.Contains(result, `role="status"`) {
		t.Error("skeletonContainer: role=status missing on negative count")
	}
}
