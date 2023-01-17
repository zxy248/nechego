package game

import (
	"testing"
	"time"
)

func TestDates(t *testing.T) {
	want := 24 * time.Hour
	diff := tomorrow().Sub(today())
	if diff != want {
		t.Errorf("tomorrow - today == %v, want %v", diff, want)
	}
}
