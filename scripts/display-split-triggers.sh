#!/usr/bin/env bash
# scripts/display-split-triggers.sh — the ADR-0040 trigger dashboard for the
# display family (the 53-component, one-package giant; v2 "display split"
# candidate). Prints the mechanical facts the decision reads: component and
# file counts, the family's external dependency surface (T1: a UNIQUE
# external dependency would fire the extraction trigger), and the afferent
# coupling (who imports display — the blast radius a split would re-path).
#
# T2 (documented consumer needs display-only adoption) and T3 (own build
# pipeline/runtime artifact) are human-judgment rows: the script states them
# and where their evidence would live, so the periodic re-check is ONE
# command instead of archaeology. Decision record: docs/adr/
# 0040-component-module-extraction-triggers.md; v2 order: ROADMAP.md.

set -euo pipefail
cd "$(dirname "$0")/.."

components="$(rg -c '^templ [A-Z][A-Za-z0-9]*\(' display/*.templ | awk -F: '{s+=$2} END {print s+0}')"
files="$(find display -maxdepth 1 -name '*.templ' | wc -l | tr -d ' ')"

echo "display components:   $components (ADR-0040 baseline: 53)"
echo "display .templ files: $files"
echo

echo "T1 — external deps of display (outside templ + tailwind-merge + workspace):"
t1="$(cd display && GOWORK=off go list -f '{{range .Imports}}{{.}}
{{end}}' . 2>/dev/null | grep '^github.com/' |
	grep -v '^github.com/a-h/templ\(/runtime\)\?$' |
	grep -v '^github.com/Oudwins/tailwind-merge-go$' |
	grep -v '^github.com/larsartmann/' || true)"
if [ -n "$t1" ]; then
	echo "$t1" | sed 's/^/    /'
	echo "    → T1 CANDIDATE: display pulls a non-workspace dependency (verify uniqueness + weight)"
else
	echo "    none — T1 unmet (display stays on the shared dependency budget)"
fi
echo

echo "Afferent importers of .../display (workspace-wide):"
afferent=""
for mod in . utils icons errorpage charts/echarts htmx datastar visualtest website; do
	afferent+="$(cd "$mod" && GOWORK=off go list -f '{{.ImportPath}} {{join .Imports " "}}' ./... 2>/dev/null || true)
"
done
echo "$afferent" | awk '$1 != "github.com/larsartmann/templ-components/display" && $0 ~ /templ-components\/display/ {print "    " $1}'
echo "    ($(echo "$afferent" | awk '$1 != "github.com/larsartmann/templ-components/display" && $0 ~ /templ-components\/display/' | wc -l | tr -d ' ') packages — the re-path blast radius of a split)"
echo

echo "T2 — documented consumer needing display-only adoption:"
echo "    judge from TODO_LIST #217 demand-gate rows + issues; no mechanical probe."
echo "T3 — own build pipeline / runtime artifact for display:"
echo "    judge from any future WASM/server-runtime work; no mechanical probe."
echo
echo "Verdict template: all triggers unmet → Keep display inside the root module"
echo "(ADR-0040); any trigger fires → extract per the ADR-0034 playbook, v2 order:"
echo "charts → overlays → tables (ROADMAP.md)."
