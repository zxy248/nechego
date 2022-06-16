package bot

import "testing"

func TestRandomNumbers(t *testing.T) {
	n := 10
	r1 := randomNumbers(n)
	r2 := randomNumbers(n)

	if len(r1) != n || len(r2) != n {
		t.Errorf("len(r1) = %v, len(r2) = %v, want %v %v", len(r1), len(r2), n, n)
	}
	if r1 == r2 {
		t.Errorf("r1 (%v) = r2 (%v), expected to be different", r1, r2)
	}
}

func TestMarkdownEscaper(t *testing.T) {
	e := newMarkdownEscaper()

	testcases := []struct {
		arg  string
		want string
	}{
		{"hello", "\\h\\e\\l\\l\\o"},
		{"(){}[!]", "\\(\\)\\{\\}\\[\\!\\]"},
	}
	for _, tc := range testcases {
		got := e.Replace(tc.arg)
		if got != tc.want {
			t.Errorf("got %v, want %v", got, tc.want)
		}
	}
}
