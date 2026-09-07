package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProbeBusyMarkup(t *testing.T) {
	server := httptest.NewServer(newMux())
	defer server.Close()
	resp, err := server.Client().Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	for _, line := range strings.Split(string(body), "<") {
		if strings.Contains(line, "wire/busy") && strings.Contains(line, "data-on") {
			t.Logf("FOUND: <%s", line)
		}
	}
}
