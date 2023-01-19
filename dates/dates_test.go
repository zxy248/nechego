package dates

import (
	"testing"
	"time"
)

func TestDates(t *testing.T) {
	want := 24 * time.Hour
	diff := Tomorrow().Sub(Today())
	if diff != want {
		t.Errorf("tomorrow - today == %v, want %v", diff, want)
	}
}
