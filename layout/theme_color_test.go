package layout

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestValidHexColor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "valid 3-digit hex", input: "#abc", want: true},
		{name: "valid 3-digit hex with digits", input: "#123", want: true},
		{name: "valid 4-digit hex (rgba short)", input: "#abcd", want: true},
		{name: "valid 6-digit hex", input: "#aabbcc", want: true},
		{name: "valid 6-digit hex default", input: DefaultThemeColor, want: true},
		{name: "valid 6-digit hex dark default", input: DefaultDarkThemeColor, want: true},
		{name: "valid 8-digit hex (rgba)", input: "#aabbccdd", want: true},
		{name: "valid uppercase hex", input: "#AABBCC", want: true},
		{name: "valid mixed case hex", input: "#AaBbCc", want: true},
		{name: "empty string is invalid", input: "", want: false},
		{name: "missing # is invalid", input: "aabbcc", want: false},
		{name: "only # is invalid", input: "#", want: false},
		{name: "## prefix is invalid", input: "##abc", want: false},
		{name: "2-digit hex is invalid", input: "#ab", want: false},
		{name: "5-digit hex is invalid", input: "#abcde", want: false},
		{name: "7-digit hex is invalid", input: "#abcdefg", want: false},
		{name: "non-hex character in 3-digit is invalid", input: "#abg", want: false},
		{name: "non-hex character in 6-digit is invalid", input: "#aabcgz", want: false},
		{name: "named color is invalid", input: "blue", want: false},
		{name: "rgb() function is invalid", input: "rgb(0,0,0)", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			utils.AssertEqual(
				t,
				"validHexColor("+tt.input+")",
				validHexColor(tt.input),
				tt.want,
			)
		})
	}
}

func TestThemeColor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{"#ff00ff", "#ff00ff"},
		{"", DefaultThemeColor},
		{"ff00ff", DefaultThemeColor},
		{"#xyz", DefaultThemeColor},
		{"#zzzzzz", DefaultThemeColor},
	}
	for _, tt := range tests {
		props := PageProps{ThemeColor: tt.input}
		utils.AssertEqual(t, "themeColor", themeColor(props), tt.want)
	}
}

func TestDarkThemeColor(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{"#000000", "#000000"},
		{"", DefaultDarkThemeColor},
		{"000000", DefaultDarkThemeColor},
		{"#xyz", DefaultDarkThemeColor},
		{"#zzzzzz", DefaultDarkThemeColor},
	}
	for _, tt := range tests {
		props := PageProps{DarkThemeColor: tt.input}
		utils.AssertEqual(t, "darkThemeColor", darkThemeColor(props), tt.want)
	}
}
