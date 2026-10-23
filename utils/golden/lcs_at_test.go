package golden

import "testing"

// TestLCSAtOutOfRange pins the defensive accessor used by the LCS diff walker.
// The walker reads one row/column past the line slices; a malformed or short
// table must return 0 instead of panicking.
func TestLCSAtOutOfRange(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		lcs  [][]int
		row  int
		col  int
		want int
	}{
		{name: "nil table", lcs: nil, row: 0, col: 0, want: 0},
		{name: "negative row", lcs: [][]int{{1}}, row: -1, col: 0, want: 0},
		{name: "row past end", lcs: [][]int{{1}}, row: 1, col: 0, want: 0},
		{name: "negative column", lcs: [][]int{{1}}, row: 0, col: -1, want: 0},
		{name: "column past end", lcs: [][]int{{1, 2}}, row: 0, col: 2, want: 0},
		{name: "ragged row", lcs: [][]int{{1}, {}}, row: 1, col: 0, want: 0},
		{name: "in range", lcs: [][]int{{1, 2}, {3, 4}}, row: 1, col: 1, want: 4},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := lcsAt(testCase.lcs, testCase.row, testCase.col); got != testCase.want {
				t.Fatalf("lcsAt(row=%d, col=%d) = %d, want %d", testCase.row, testCase.col, got, testCase.want)
			}
		})
	}
}
