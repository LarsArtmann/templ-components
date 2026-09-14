package build

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// RenderedPage is one fully rendered output document, ready to be written.
type RenderedPage struct {
	// Path is the dist-relative file path, e.g. "index.html".
	Path string
	// HTML is the rendered document.
	HTML string
}

var scriptBlockRe = regexp.MustCompile(`(?s)<script\b([^>]*)>(.*?)</script>`)

// InlineScriptHashes extracts CSP source hashes for every inline (src-less)
// <script> body across the rendered pages. Duplicate bodies are hashed once
// and the result is sorted, so the header is byte-stable across builds —
// per-build nonces live in attributes, never in the hashed body. JSON-LD data
// blocks are included defensively (browsers exempt them from script-src, but
// hashing them costs nothing and keeps older engines covered).
func InlineScriptHashes(pages []RenderedPage) []string {
	seen := map[string]bool{}

	hashes := []string{}

	for _, page := range pages {
		for _, match := range scriptBlockRe.FindAllStringSubmatch(page.HTML, -1) {
			attrs, body := match[1], match[2]

			if strings.Contains(attrs, "src=") {
				continue
			}

			digest := sha256.Sum256([]byte(body))
			hash := "'sha256-" + base64.StdEncoding.EncodeToString(digest[:]) + "'"

			if seen[hash] {
				continue
			}

			seen[hash] = true

			hashes = append(hashes, hash)
		}
	}

	sort.Strings(hashes)

	return hashes
}

// CSPHeader renders the site's Content-Security-Policy value. Everything is
// same-origin ('self'); the inline scripts are pinned by content hash so the
// header can be committed to firebase.json and stays valid across builds.
func CSPHeader(scriptHashes []string) string {
	scriptSrc := "script-src 'self'"
	if len(scriptHashes) > 0 {
		scriptSrc += " " + strings.Join(scriptHashes, " ")
	}

	directives := []string{
		"default-src 'self'",
		scriptSrc,
		"style-src 'self'",
		"img-src 'self' data:",
		"font-src 'self'",
		"connect-src 'self'",
		"manifest-src 'self'",
		"object-src 'none'",
		"base-uri 'self'",
		"form-action 'self'",
		"frame-ancestors 'none'",
	}

	return strings.Join(directives, "; ")
}

// cspHeaderKey is the firebase.json header name pinned to the CSP value.
const cspHeaderKey = "Content-Security-Policy"

var (
	errCSPHeaderMissing = errors.New("no Content-Security-Policy header")
	errCSPHeaderStale   = errors.New("CSP header is stale (pages changed their inline scripts)")
)

// wildcardHeaderSource is the firebase.json headers entry that applies to
// every response ("**").
const wildcardHeaderSource = "**"

// SyncFirebaseCSP writes the CSP header into the firebase.json headers config
// (creating the wildcard entry if needed) and reports whether the file
// changed. The round-trip re-indents the file with sorted keys; the result is
// deterministic, so repeated syncs are byte-stable.
func SyncFirebaseCSP(path, header string) (bool, error) {
	config, err := readFirebaseConfig(path)
	if err != nil {
		return false, err
	}

	changed := setCSPHeader(config, header)
	if !changed {
		return false, nil
	}

	if err := writeFirebaseConfig(path, config); err != nil {
		return false, err
	}

	return true, nil
}

// CheckFirebaseCSP verifies that firebase.json pins exactly the CSP the
// rendered pages need. A mismatch is a build failure: the committed header
// would otherwise not cover (or over-permit) the deployed inline scripts.
func CheckFirebaseCSP(path, header string) error {
	config, err := readFirebaseConfig(path)
	if err != nil {
		return err
	}

	current, ok := findCSPHeader(config)
	if !ok {
		return fmt.Errorf("%s: %w; run: go run ./cmd/site --update-csp", path, errCSPHeaderMissing)
	}

	if current != header {
		return fmt.Errorf("%s: %w; run: go run ./cmd/site --update-csp", path, errCSPHeaderStale)
	}

	return nil
}

func readFirebaseConfig(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	return config, nil
}

func writeFirebaseConfig(path string, config map[string]any) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}

	//nolint:gosec // repository-controlled config path
	if err := os.WriteFile(
		path,
		append(data, '\n'),
		0o644,
	); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}

// findCSPHeader returns the current CSP value from the wildcard headers
// entry, if present.
func findCSPHeader(config map[string]any) (string, bool) {
	for _, entry := range hostingHeaders(config) {
		entryMap, ok := entry.(map[string]any)
		if !ok || entryMap["source"] != wildcardHeaderSource {
			continue
		}

		headers, _ := entryMap["headers"].([]any)
		for _, header := range headers {
			headerMap, ok := header.(map[string]any)
			if !ok {
				continue
			}

			if key, _ := headerMap["key"].(string); key == cspHeaderKey {
				value, _ := headerMap["value"].(string)

				return value, true
			}
		}
	}

	return "", false
}

// setCSPHeader stores the header in the wildcard entry, creating entry or
// header as needed. Returns true when the config content changed.
func setCSPHeader(config map[string]any, header string) bool {
	if current, ok := findCSPHeader(config); ok && current == header {
		return false
	}

	hosting, _ := config["hosting"].(map[string]any)
	if hosting == nil {
		hosting = map[string]any{}
		config["hosting"] = hosting
	}

	headers := hostingHeaders(config)

	for _, entry := range headers {
		entryMap, ok := entry.(map[string]any)
		if !ok || entryMap["source"] != wildcardHeaderSource {
			continue
		}

		existing, _ := entryMap["headers"].([]any)
		for _, headerEntry := range existing {
			headerMap, ok := headerEntry.(map[string]any)
			if !ok {
				continue
			}

			if key, _ := headerMap["key"].(string); key == cspHeaderKey {
				headerMap["value"] = header

				return true
			}
		}

		entryMap["headers"] = append(existing, map[string]any{"key": cspHeaderKey, "value": header})

		return true
	}

	hosting["headers"] = append(headers, map[string]any{
		"source": wildcardHeaderSource,
		"headers": []any{
			map[string]any{"key": cspHeaderKey, "value": header},
		},
	})

	return true
}

// hostingHeaders returns the hosting.headers array (empty when absent).
func hostingHeaders(config map[string]any) []any {
	hosting, ok := config["hosting"].(map[string]any)
	if !ok {
		return nil
	}

	headers, _ := hosting["headers"].([]any)

	return headers
}
