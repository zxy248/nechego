package app

import (
	"regexp"
	"testing"
)

func TestProbabilityTemplates(t *testing.T) {
	re := regexp.MustCompile("^.*%s.*%d%%")
	for _, s := range probabilityTemplates {
		if !re.MatchString(s) {
			t.Errorf("want %q to match regexp", s)
		}
	}
}
