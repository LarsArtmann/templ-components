package visualtest

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestZZDumpWizardResponses(t *testing.T) {
	srv := packE2EServer(t)

	post := func(body string) string {
		req, _ := http.NewRequest("POST", srv.URL+"/api/pack/wizard", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Datastar-Request", "true")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)
		t.Logf("status=%d len=%d selector=%q mode=%q", resp.StatusCode, len(b), resp.Header.Get("Datastar-Selector"), resp.Header.Get("Datastar-Mode"))
		return string(b)
	}

	t.Log("RESP1:\n", post("step=0&email="))
	t.Log("RESP2:\n", post("step=0&email=ada@example.com"))
	t.Log("RESP3:\n", post("step=1&name="))
}
