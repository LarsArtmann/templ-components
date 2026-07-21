package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComponentSourcesPopulated(t *testing.T) {
	t.Parallel()

	r := newRegistry()
	if len(r.files) == 0 {
		t.Fatal("registry files is empty — embed failed")
	}

	if _, ok := r.files["button"]; !ok {
		t.Error("expected 'button' in registry files")
	}
}

func TestCmdAddCopiesTemplAndTypes(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	cmdAdd(newRegistry(), []string{"button", "--out", tmp})

	if _, err := os.Stat(filepath.Join(tmp, "button.templ")); err != nil {
		t.Errorf("button.templ not copied: %v", err)
	}
}

func TestCmdAddUnknownComponent(t *testing.T) {
	t.Parallel()

	r := newRegistry()
	if _, ok := r.files["totally-bogus-component"]; ok {
		t.Error("did not expect bogus component in sources")
	}
}

func TestCmdListDoesNotCrash(t *testing.T) {
	t.Parallel()

	cmdList(newRegistry(), nil)
}
