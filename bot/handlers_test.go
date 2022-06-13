package bot

import "testing"

func TestRandomNumbers(t *testing.T) {
	r1 := randomNumbers(5)
	r2 := randomNumbers(5)

	if len(r1) != 5 || len(r2) != 5 {
		t.Errorf("len(r1) == %v, len(r2) == %v, %v %v expected", len(r1), len(r2), 5, 5)
	}

	if r1 == r2 {
		t.Errorf("r1 (%v) == r2 (%v), expected to be different", r1, r2)
	}
}

func TestMarkdownEscaper(t *testing.T) {
	e := newMarkdownEscaper()

	r := e.Replace("hello")
	if r != "\\h\\e\\l\\l\\o" {
		t.Fatalf(`r == %v, expected "\\h\\e\\l\\l\\o"`, r)
	}

	r = e.Replace("(){}[!]")
	if r != "\\(\\)\\{\\}\\[\\!\\]" {
		t.Fatalf(`r == %v, expected "\\(\\)\\{\\}\\[\\!\\]"`, r)
	}
}
