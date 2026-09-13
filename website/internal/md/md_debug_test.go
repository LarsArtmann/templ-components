package md

import (
	"strings"
	"testing"

	"github.com/alecthomas/chroma/v2/styles"
)

func TestDebugStylesAndTempl(t *testing.T) {
	t.Log("github-light exists:", styles.Get("github-light") != nil, "github-dark:", styles.Get("github-dark") != nil)
	names := styles.Names()
	var gh []string
	for _, n := range names {
		if strings.Contains(n, "github") {
			gh = append(gh, n)
		}
	}
	t.Log("github styles:", gh)

	p, err := Parse("```templ\n templ Page() {\n  @layout.Base()\n}\n```\n")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("HTML:", p.HTML)
}
