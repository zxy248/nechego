package app

import (
	"regexp"
	"testing"
)

func TestMarkdownEscaper(t *testing.T) {
	testcases := []struct {
		arg  string
		want string
	}{
		{"hello", "\\h\\e\\l\\l\\o"},
		{"(){}[!]", "\\(\\)\\{\\}\\[\\!\\]"},
	}
	for _, tc := range testcases {
		got := markdownEscaper.Replace(tc.arg)
		if got != tc.want {
			t.Errorf("got %v, want %v", got, tc.want)
		}
	}
}

func TestProbabilityTemplates(t *testing.T) {
	re := regexp.MustCompile("^.*%s.*%d%%")
	for _, s := range probabilityTemplates {
		if !re.MatchString(s) {
			t.Errorf("want %q to match regexp", s)
		}
	}
}
