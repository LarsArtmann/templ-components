package utils

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"
)

// Finding is a single compliance-scanner violation located in a source file.
// The scanners behind the repo's guard tests share this model so the same
// scan logic can run inside `go test` (repo guards) and inside the `tc lint`
// CLI (consumer adoption).
type Finding struct {
	// File is the path relative to the directory the caller passed.
	File string
	// Line is the 1-based line number in the original file.
	Line int
	// Rule is a stable, machine-readable rule identifier (e.g. "rtl-physical").
	Rule string
	// Message is a human-readable description with the fix hint.
	Message string
	// Snippet is the trimmed offending source line.
	Snippet string
}

// scanSources walks dir and hands every file whose name passes include to fn
// as (dir-relative path, content). Reads go through os.Root so the scanner
// cannot be walked out of the target directory by a swapped symlink, and a
// missing or unreadable directory is an error, never a silent "clean".
func scanSources(dir string, include func(name string) bool, fn func(rel, content string) error) error {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("open %s: %w", dir, err)
	}

	defer func() { _ = root.Close() }()

	walkErr := fs.WalkDir(root.FS(), ".", func(rel string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !include(d.Name()) {
			return nil
		}

		data, readErr := root.ReadFile(rel)
		if readErr != nil {
			return fmt.Errorf("read %s: %w", rel, readErr)
		}

		return fn(rel, string(data))
	})
	if walkErr != nil {
		return fmt.Errorf("scan %s: %w", dir, walkErr)
	}

	return nil
}

// lintIgnoreMarker is the per-line suppression marker for tc lint. It must
// appear on the SAME line as the finding, inside any comment syntax the file
// language offers (`<!-- tc-lint:ignore -->` in .templ, `// tc-lint:ignore`
// in .go). Bare `tc-lint:ignore` suppresses every rule on the line;
// `tc-lint:ignore <rule> [<rule>...]` suppresses exactly the named rules.
const lintIgnoreMarker = "tc-lint:ignore"

// lintIgnoreTrimCuts strips leftover comment punctuation from the tokens
// between the marker and the comment closer.
const lintIgnoreTrimCuts = "\"',;()/-*"

// lintRuleRe is the rule-id shape: lowercase letters/digits/hyphens starting
// with a letter. It keeps markup after the comment closer from leaking in as
// bogus rule names.
var lintRuleRe = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

// lintIgnoreFor parses one raw source line's suppression directive. Only the
// text between the marker and the comment closer (-->, */) is considered, so
// later markup never becomes a rule name. No rules in scope means "suppress
// every rule on this line".
func lintIgnoreFor(line string) (bool, map[string]bool) {
	_, rest, found := strings.Cut(line, lintIgnoreMarker)
	if !found {
		return false, nil
	}

	for _, closer := range []string{"-->", "*/"} {
		if end := strings.Index(rest, closer); end >= 0 {
			rest = rest[:end]
		}
	}

	var rules map[string]bool

	for token := range strings.FieldsSeq(rest) {
		rule := strings.Trim(token, lintIgnoreTrimCuts)
		if !lintRuleRe.MatchString(rule) {
			continue
		}

		if rules == nil {
			rules = map[string]bool{}
		}

		rules[rule] = true
	}

	if rules == nil {
		return true, nil // bare marker: suppress every rule on this line
	}

	return false, rules
}

// lintFindingSuppressed reports whether the raw line's ignore directive
// suppresses a finding of the given rule.
func lintFindingSuppressed(rawLine, rule string) bool {
	all, rules := lintIgnoreFor(rawLine)
	if all {
		return true
	}

	return rules[rule]
}

// htmlCommentAnyRe matches HTML comments across line breaks. An unterminated
// comment is left in place, so its text stays scannable — a violation-shaped
// token inside broken markup still gets reported (fails safe).
var htmlCommentAnyRe = regexp.MustCompile(`(?s)<!--.*?-->`)

// stripCommentsKeepLines blanks HTML comments with spaces while PRESERVING
// newlines, so scanners can report true original line numbers and class-shaped
// compliance regexes judge markup, not prose. Both the motion-reduce and RTL
// guards false-positived on comment prose (2026-10-08).
func stripCommentsKeepLines(content string) string {
	return htmlCommentAnyRe.ReplaceAllStringFunc(content, func(comment string) string {
		return strings.Map(func(r rune) rune {
			if r == '\n' {
				return r
			}

			return ' '
		}, comment)
	})
}

// rtlPhysicalRe matches physical directional Tailwind properties that have
// direct logical equivalents (ml-, mr-, pl-, pr-, text-left, text-right,
// border-l-, border-r-). These have ZERO exceptions — every occurrence is a
// violation.
var rtlPhysicalRe = regexp.MustCompile(`\b(ml-|mr-|pl-|pr-|text-left|text-right|border-l-|border-r-)`)

// rtlInsetAliasRe matches the deprecated physical-positioning inset
// utilities start-*/end-* (e.g. start-0, end-full): probe-verified 2026-10-08
// as output-identical LOGICAL aliases of inset-s-*/inset-e-*, but a future
// major can drop the aliases, so the library stays on the canonical forms.
// The preceding-character class keeps longer classes safe (inset-s-0,
// text-start, items-start, justify-end, col-start-2, rounded-s-*).
var rtlInsetAliasRe = regexp.MustCompile(`(?:^|[\s"'({])(start|end)-[a-z0-9.]`)

// ScanRTL walks dir recursively and reports RTL logical-property violations in
// .templ files: physical directional classes (ml-/mr-/pl-/pr-/text-left/
// text-right/border-l-/border-r-) and the deprecated inset aliases
// (start-*/end-*). Logical properties automatically mirror in RTL (dir="rtl")
// without code changes; physical properties do not and cause broken RTL
// layouts. HTML comments are blanked before scanning, so prose never trips
// the class regexes.
//
// It is the engine behind TestRTLLogicalProperties (repo guard) and
// `tc lint --rtl` (consumer adoption). Findings use dir-relative paths and are
// returned in walk order.
func ScanRTL(dir string) ([]Finding, error) {
	var findings []Finding

	err := scanSources(dir, func(name string) bool { return strings.HasSuffix(name, ".templ") },
		func(rel, content string) error {
			raw := strings.Split(content, "\n")
			masked := strings.Split(stripCommentsKeepLines(content), "\n")

			for i := range raw {
				switch {
				case rtlPhysicalRe.MatchString(masked[i]) && !lintFindingSuppressed(raw[i], "rtl-physical"):
					findings = append(findings, Finding{
						File:    rel,
						Line:    i + 1,
						Rule:    "rtl-physical",
						Message: "RTL physical-property violation (use the logical ms-/me-/ps-/pe-/text-start/text-end/border-s-/border-e- form)",
						Snippet: strings.TrimSpace(masked[i]),
					})
				case rtlInsetAliasRe.MatchString(masked[i]) && !lintFindingSuppressed(raw[i], "rtl-inset-alias"):
					findings = append(findings, Finding{
						File:    rel,
						Line:    i + 1,
						Rule:    "rtl-inset-alias",
						Message: "RTL physical-inset violation (use inset-s-*/inset-e-* instead of start-*/end-*)",
						Snippet: strings.TrimSpace(masked[i]),
					})
				}
			}

			return nil
		})
	if err != nil {
		return nil, err
	}

	return findings, nil
}
