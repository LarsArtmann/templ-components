#!/usr/bin/env bash
# First-pass AI vision review of visual-regression golden PNGs.
#
# WHY: the visual goldens were captured by automated agents; the recurring
# TODO (#80/#150/#162) says a human must confirm no rendering bug got
# enshrined as a baseline. The vision-review-agent
# (github.com/LarsArtmann/vision-review-agent) can do that first pass: it
# reads PNGs and flags clipped text, misaligned overlays, broken top-layer
# positioning, illegible contrast, and layout wreckage. A human then reviews
# only the flagged findings.
#
# Usage:
#   scripts/vision-review-goldens.sh            # the TODO #80/#150/#162 flagged set
#   scripts/vision-review-goldens.sh --all      # every golden PNG in visualtest/testdata
#   scripts/vision-review-goldens.sh modal/*.png drawer/*.png   # explicit globs
#
# Environment:
#   VISION_BIN   optional command prefix for the CLI. Default:
#                nix run github:LarsArtmann/vision-review-agent --
#   One provider key among OPENAI_API_KEY, ANTHROPIC_API_KEY,
#   GEMINI_API_KEY, OPENROUTER_API_KEY, XAI_API_KEY (see -list-providers).
#   VISION_MODEL optional model override (default: gpt-4o).
#
# Output: docs/reviews/vision-golden-review-<date>.md plus stdout.
# Cost: one vision call per image. The flagged set is ~20 images; --all is
# 125+ — consent via the explicit flag.

set -euo pipefail

REPO_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$REPO_ROOT"
TESTDATA=visualtest/testdata

VISION_BIN=${VISION_BIN:-"nix run github:LarsArtmann/vision-review-agent --"}
VISION_MODEL=${VISION_MODEL:-}

# --- The flagged set (TODO #80: overlays + datastar/eyebrow/scrollback +
# statcard yellow/purple; TODO #150: wire; TODO #162: progressbar;
# TODO #223: kanban action/tone section) ---
FLAGGED=(
	"dropdown/open_light.png" "dropdown/open_dark.png"
	"popover/open_light.png"
	"contextmenu/open_light.png"
	"modal/open_light.png" "modal/open_dark.png"
	"drawer/left_dark.png" "drawer/right_light.png"
	"datastar/indicator_light.png" "datastar/indicator_dark.png"
	"datastar/live_region_light.png" "datastar/live_region_dark.png"
	"eyebrow/light.png" "eyebrow/dark.png"
	"scrollback/light.png" "scrollback/dark.png"
	"statcard/yellow_light.png" "statcard/purple_light.png"
	"button/outline_success_light.png" "button/outline_info_light.png"
	"progressbar/half_light.png"
	"wire/dual_transport_light.png" "wire/dual_transport_dark.png"
	"kanban/section_action_tone_light.png" "kanban/section_action_tone_dark.png"
)

# --- Argument handling: --all or explicit globs replace the flagged set ---
IMAGES=()
for arg in "$@"; do
	if [ "$arg" = "--all" ]; then
		mapfile -t IMAGES < <(find "$TESTDATA" -name '*.png' -not -name '*.fail.png' | sort)
		break
	fi
	IMAGES+=("$arg")
done
if [ ${#IMAGES[@]} -eq 0 ]; then
	for rel in "${FLAGGED[@]}"; do
		[ -f "$TESTDATA/$rel" ] && IMAGES+=("$TESTDATA/$rel")
	done
fi
if [ ${#IMAGES[@]} -eq 0 ]; then
	echo "ERROR: no golden PNGs matched. What: the flagged list and globs resolved to nothing." >&2
	echo "Fix: run from the repo root, or pass globs like 'modal/*.png', or --all." >&2
	exit 1
fi

# --- Fail fast on missing credentials (What/Why/Fix, per repo error policy) ---
if [ -z "${OPENAI_API_KEY:-}${ANTHROPIC_API_KEY:-}${GEMINI_API_KEY:-}${OPENROUTER_API_KEY:-}${XAI_API_KEY:-}" ]; then
	cat >&2 <<EOF
ERROR: no vision provider API key found.
What: this script calls the vision-review-agent CLI, which needs a vision model.
Why: none of OPENAI_API_KEY, ANTHROPIC_API_KEY, GEMINI_API_KEY,
     OPENROUTER_API_KEY, XAI_API_KEY is set.
Fix:  export one of them (e.g. export OPENAI_API_KEY=sk-...) and re-run.
      Cost control: the default run reviews ${#IMAGES[@]} flagged goldens only;
      pass --all (125+) deliberately.
EOF
	exit 1
fi

MODEL_ARGS=()
[ -n "$VISION_MODEL" ] && MODEL_ARGS=(-model "$VISION_MODEL")

# shellcheck disable=SC2206 # VISION_BIN is intentionally a command prefix
VISION_CMD=(${VISION_BIN})

PROMPT='This PNG is a golden baseline screenshot of one server-rendered Go UI component (Tailwind v4 styling; the component name is in the file path). Review it as a strict UI QA engineer and answer in four short bullets: (1) CLIPPING: is any text, icon, badge, or button clipped or cut off? (2) LAYOUT: is anything misaligned, overlapping, or overflowing its container (for overlays: is the panel positioned sensibly relative to its trigger or screen edges, not detached in a corner)? (3) CONTRAST: is any text illegible against its background (light AND dark variants)? (4) VERDICT: one of CLEAN or SUSPECT plus a one-line reason. If the image shows a plain empty/idle state with nothing visible, say so instead of inventing problems.'

REPORT="docs/reviews/vision-golden-review-$(date +%Y-%m-%d).md"
mkdir -p docs/reviews

echo "# Vision golden review — $(date +%Y-%m-%d)" >"$REPORT"
echo "" >>"$REPORT"
echo "Tool: vision-review-agent (first pass; human reviews SUSPECT findings)." >>"$REPORT"
echo "Images: ${#IMAGES[@]}" >>"$REPORT"
echo "" >>"$REPORT"

suspect=0
for img in "${IMAGES[@]}"; do
	echo "==> $img"
	output=$("${VISION_CMD[@]}" "${MODEL_ARGS[@]}" -prompt "$PROMPT" "$img" 2>&1)
	{
		echo "## $img"
		echo ""
		echo '```'
		echo "$output"
		echo '```'
		echo ""
	} >>"$REPORT"
	if grep -q "SUSPECT" <<<"$output"; then
		suspect=$((suspect + 1))
		echo "    ^^ SUSPECT — needs human eyes"
	fi
done

{
	echo "---"
	echo "Total: ${#IMAGES[@]} images, $suspect flagged SUSPECT."
} >>"$REPORT"

echo ""
echo "Report written: $REPORT ($suspect SUSPECT of ${#IMAGES[@]})"
if [ "$suspect" -gt 0 ]; then
	echo "Review the SUSPECT sections — a human confirms before any golden changes."
fi
