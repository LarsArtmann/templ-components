package utils

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

// failingScanWriter always fails on Write, modelling a closed connection.
type failingScanWriter struct{}

func (failingScanWriter) Write([]byte) (int, error) {
	return 0, errors.New("connection closed")
}

// TestScannersReportMissingDirectory proves every repo guard surfaces an
// unreadable root as an error instead of silently reporting zero findings
// (a silent scanner would turn a whole compliance lint into a no-op).
func TestScannersReportMissingDirectory(t *testing.T) {
	t.Parallel()

	missing := filepath.Join(t.TempDir(), "does-not-exist")

	scanners := map[string]func(string) ([]Finding, error){
		"RTL":        ScanRTL,
		"dark mode":  ScanDarkMode,
		"motion":     ScanMotionReduce,
		"icon paths": ScanIconPaths,
	}

	for name, scan := range scanners {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			findings, err := scan(missing)
			if err == nil {
				t.Fatalf("%s scanner returned nil error for a missing dir (findings %v)", name, findings)
			}
		})
	}
}

// TestScriptComponentPropagatesWriteErrors pins the CSP-safety writer's error
// contract: a failed write is reported with the caller's label so a broken
// response is never mistaken for a rendered script.
func TestScriptComponentPropagatesWriteErrors(t *testing.T) {
	t.Parallel()

	err := ScriptComponent("", "console.log(1)", "test script").
		Render(context.Background(), failingScanWriter{})
	if err == nil {
		t.Fatal("ScriptComponent swallowed a write error")
	}

	if !strings.Contains(err.Error(), "test script") {
		t.Errorf("error %q does not carry the caller label", err)
	}
}
