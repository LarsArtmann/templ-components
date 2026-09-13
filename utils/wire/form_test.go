package wire

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type decodeFormTarget struct {
	Card   string `form:"card"`
	Column string `form:"column"`
	Index  int    `form:"index"`
	Flag   bool   `form:"flag"`
	Count  int64  `form:"count"`
	Skipped string
}

func TestDecodeFormPostBody(t *testing.T) {
	t.Parallel()

	body := url.Values{"card": {"c1"}, "column": {"todo"}, "index": {"2"}, "flag": {"true"}, "count": {"9007199254740993"}}
	request, err := http.NewRequest(http.MethodPost, "/move", strings.NewReader(body.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	got, err := DecodeForm[decodeFormTarget](request)
	if err != nil {
		t.Fatalf("DecodeForm: %v", err)
	}

	want := decodeFormTarget{Card: "c1", Column: "todo", Index: 2, Flag: true, Count: 9007199254740993}
	if got != want {
		t.Errorf("DecodeForm = %+v, want %+v", got, want)
	}
}

func TestDecodeFormQueryParameters(t *testing.T) {
	t.Parallel()

	request, err := http.NewRequest(http.MethodGet, "/search?q=x", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	type target struct {
		Query string `form:"q"`
		Page  int    `form:"page"`
	}

	got, err := DecodeForm[target](request)
	if err != nil {
		t.Fatalf("DecodeForm: %v", err)
	}

	if got.Query != "x" || got.Page != 0 {
		t.Errorf("DecodeForm = %+v, want Query=x Page=0 (absent field stays zero)", got)
	}
}

func TestDecodeFormCheckboxOn(t *testing.T) {
	t.Parallel()

	request, err := http.NewRequest(http.MethodPost, "/move", strings.NewReader("flag=on"))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	got, err := DecodeForm[decodeFormTarget](request)
	if err != nil {
		t.Fatalf("DecodeForm: %v", err)
	}

	if !got.Flag {
		t.Error("DecodeForm flag=on = false, want true (HTML checkbox convention)")
	}
}

func TestDecodeFormMalformedInt(t *testing.T) {
	t.Parallel()

	request, err := http.NewRequest(http.MethodPost, "/move", strings.NewReader("index=abc"))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if _, err := DecodeForm[decodeFormTarget](request); err == nil {
		t.Fatal("DecodeForm with index=abc: want error, got nil")
	} else if !strings.Contains(err.Error(), "index") {
		t.Errorf("error should name the field, got: %v", err)
	}
}

func TestDecodeFormUnsupportedKind(t *testing.T) {
	t.Parallel()

	type target struct {
		Price float64 `form:"price"`
	}

	request, err := http.NewRequest(http.MethodPost, "/buy", strings.NewReader("price=9.99"))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = DecodeForm[target](request)
	if !errors.Is(err, ErrUnsupportedFormField) {
		t.Errorf("DecodeForm float field: want ErrUnsupportedFormField, got %v", err)
	}
}
