package utils

import (
	"regexp"
	"strings"
)

// motionReduceReGates bundles the regexes the motion-reduce rule applies per line.
type motionReduceRules struct {
	transitionRe      *regexp.Regexp
	animateRe         *regexp.Regexp
	motionReduceRe    *regexp.Regexp
	transitionConstRe *regexp.Regexp
}

// defaultMotionReduceRules is the guard's rule set: a line carrying
// transition/animate classes needs either an inline motion-reduce fallback or
// a reference to the shared motion constants (which already include the
// fallback).
//
//nolint:gochecknoglobals // declarative rule set
var defaultMotionReduceRules = motionReduceRules{
	transitionRe:   regexp.MustCompile(`transition[-:]`),
	animateRe:      regexp.MustCompile(`animate[-:]`),
	motionReduceRe: regexp.MustCompile(`motion-reduce:`),
	transitionConstRe: regexp.MustCompile(
		`(utils\.Transition|transitionFast|transitionNormal|transitionColors|transitionTransform)`,
	),
}

// needsMotionReduce reports whether a line carries transition/animate classes
// without an inline motion-reduce fallback or a shared motion-constant
// reference (the constants already include motion-reduce fallbacks).
func (r motionReduceRules) needsMotionReduce(line string) bool {
	if !r.transitionRe.MatchString(line) && !r.animateRe.MatchString(line) {
		return false
	}

	if r.motionReduceRe.MatchString(line) {
		return false
	}

	return !r.transitionConstRe.MatchString(line)
}

// ScanMotionReduce walks dir recursively and reports transition/animation
// classes in .templ files that lack a motion-reduce fallback — the
// accessibility safety net for users with reduced-motion enabled. HTML
// comments are blanked before scanning, so prose never trips the class
// regexes.
//
// It is the engine behind TestMotionReduceCompliance (repo guard) and
// `tc lint --motion`. Findings use dir-relative paths.
func ScanMotionReduce(dir string) ([]Finding, error) {
	var findings []Finding

	err := scanSources(dir, func(name string) bool { return strings.HasSuffix(name, ".templ") },
		func(rel, content string) error {
			raw := strings.Split(content, "\n")
			masked := strings.Split(stripCommentsKeepLines(content), "\n")

			for i := range raw {
				if !defaultMotionReduceRules.needsMotionReduce(masked[i]) ||
					lintFindingSuppressed(raw[i], "motion-reduce") {
					continue
				}

				findings = append(findings, Finding{
					File:    rel,
					Line:    i + 1,
					Rule:    "motion-reduce",
					Message: "motion-reduce gap (transition/animate class without a motion-reduce fallback or shared motion constant)",
					Snippet: strings.TrimSpace(masked[i]),
				})
			}

			return nil
		})
	if err != nil {
		return nil, err
	}

	return findings, nil
}
