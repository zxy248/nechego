package input

import "testing"

func TestConstructHelloRe(t *testing.T) {
	args := []string{"п[рл]ивет", "хай", "зд[ао]ров", "ку"}
	want := "(?i)((^|[^а-я])п[рл]ивет[а-я]*([^а-я]|$))" +
		"|((^|[^а-я])хай[а-я]*([^а-я]|$))" +
		"|((^|[^а-я])зд[ао]ров[а-я]*([^а-я]|$))" +
		"|((^|[^а-я])ку[а-я]*([^а-я]|$))"

	re := constructHelloRe(args...)
	if re != want {
		t.Errorf("re = %q, want %q", re, want)
	}
}
