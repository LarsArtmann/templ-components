package wire

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// FuzzDecodeForm exercises DecodeForm against adversarial form bodies
// (M20/F091): deeply nested structs, weird tag spellings, huge values, and
// hostile type coercions. The decoder must never panic — only return a value
// or an error. Malformed bodies intentionally reach ParseForm as raw bytes.
type fuzzNested struct {
	Inner struct {
		Deeper struct {
			Leaf string `form:"leaf"`
		} `form:"deeper"`
	} `form:"inner"`
	Values []string `form:"values"`
	Number int      `form:"number"`
}

type fuzzWeirdTags struct {
	Ok      string `form:"ok"`
	NoTag   string
	Empty   string `form:""`
	Spaced  string `form:"spaced tag"`
	Unicode string `form:"größe"`
}

func FuzzDecodeForm(f *testing.F) {
	f.Add("leaf=x", "POST")
	f.Add("inner.deeper.leaf=x", "POST")
	f.Add("values=a&values=b&values=c", "POST")
	f.Add("number=99999999999999999999", "POST")
	f.Add("number=-1", "POST")
	f.Add("number=abc", "POST")
	f.Add("größe=ü", "POST")
	f.Add(strings.Repeat("a=", 1000), "POST")
	f.Add("q="+strings.Repeat("x", 100_000), "POST")
	f.Add("%zz=broken&%ww=encoding", "POST")
	f.Add("a=1;b=2", "POST")
	f.Add("q=fuzz", "GET")

	f.Fuzz(func(t *testing.T, body, method string) {
		var req *http.Request

		if method == http.MethodGet {
			req = newFuzzRequest(t, "http://x.test/?"+sanitizeFuzzQuery(body))
		} else {
			req = newFuzzRequest(t, "http://x.test/")
			req.Method = http.MethodPost
			req.Header = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
			req.Body = io.NopCloser(strings.NewReader(body))
			req.ContentLength = int64(len(body))
		}

		// Both shapes must be panic-free; errors are fine.
		deep, _ := DecodeForm[fuzzNested](req)
		_ = deep

		weird, _ := DecodeForm[fuzzWeirdTags](req)
		_ = weird
	})
}

// sanitizeFuzzQuery keeps the URL parseable: control bytes, spaces, and
// fragment markers are dropped so the request constructor always succeeds.
func sanitizeFuzzQuery(raw string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f || r == ' ' || r == '#' {
			return -1
		}
		return r
	}, raw)
}

func newFuzzRequest(t *testing.T, target string) *http.Request {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, target, http.NoBody)
	if err != nil {
		t.Fatalf("seed URL must be parseable after sanitize: %v", err)
	}

	return req
}
