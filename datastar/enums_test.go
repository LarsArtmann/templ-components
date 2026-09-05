package datastar

import (
	"testing"
)

func TestDatastarVersionIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value DatastarVersion
		want  bool
	}{
		{DatastarVersion1_0_3, true},
		{"0.9.0", false},
		{"", false},
		{"bogus", false},
	}
	for _, tt := range tests {
		got := DatastarVersionIsValid(tt.value)
		if got != tt.want {
			t.Errorf("DatastarVersionIsValid(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestLivePolitenessIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value LivePoliteness
		want  bool
	}{
		{LivePolite, true},
		{LiveAssertive, true},
		{LiveOff, true},
		{"bogus", false},
		{"", false},
	}
	for _, tt := range tests {
		got := LivePolitenessIsValid(tt.value)
		if got != tt.want {
			t.Errorf("LivePolitenessIsValid(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestRetryModeIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value RetryMode
		want  bool
	}{
		{RetryAuto, true},
		{RetryAlways, true},
		{RetryError, true},
		{RetryNever, true},
		{"bogus", false},
		{"", false},
	}
	for _, tt := range tests {
		got := RetryModeIsValid(tt.value)
		if got != tt.want {
			t.Errorf("RetryModeIsValid(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}
