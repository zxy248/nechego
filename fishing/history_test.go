package fishing

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func fakeHistory(size int) (*History, []*Fish) {
	hist := &History{}
	fish := []*Fish{}
	for i := 0; i < size; i++ {
		f := &Fish{Length: float64(i), Weight: float64(i)}
		fish = append(fish, f)
		hist.Add(int64(i), f)
	}
	// Wait for add goroutines to finish.
	time.Sleep(10 * time.Millisecond)
	return hist, fish
}

func TestHistoryTop(t *testing.T) {
	testNums := []int{0, 1, 2, 3, 10, 1000, 1000000}
	t.Run("empty top", func(t *testing.T) {
		h := &History{}
		for _, i := range testNums {
			for _, p := range parameters {
				empty := h.Top(p, i)
				if len(empty) != 0 {
					t.Errorf("want empty top")
				}
			}
		}
	})
	const histSize = 1000
	hist, fish := fakeHistory(histSize)
	for _, l := range testNums {
		for _, p := range parameters {
			top := hist.Top(p, l)
			if len(top) != l && l <= len(top) {
				t.Errorf("len(top) = %v, want %v", len(top), l)
			}
			for i, r := range top {
				j := histSize - 1 - i
				want := fish[j]
				if r.TUID != int64(j) {
					t.Errorf("r.TUID = %v, want %v", r.TUID, j)
				}
				if r.Fish != want {
					t.Errorf("r.Fish = %v, want %v", r.Fish, want)
				}
			}
		}
	}
}

func TestHistoryUnmarshal(t *testing.T) {
	testNums := []int{0, 1, 2, 3, 10, 1000, 1000000}

	h1, _ := fakeHistory(1000)
	raw, err := json.Marshal(h1)
	if err != nil {
		t.Fatal(err)
	}

	var h2 *History
	if err := json.Unmarshal(raw, &h2); err != nil {
		t.Fatal(err)
	}

	for _, l := range testNums {
		for _, p := range parameters {
			t1 := h1.Top(p, l)
			t2 := h2.Top(p, l)
			if len(t1) != l && l <= len(t1) {
				t.Errorf("len(t1) = %v, want %v", len(t1), l)
			}
			if len(t1) != len(t2) {
				t.Errorf("len(t1) = %v, want %v", len(t1), len(t2))
			}
			for i := 0; i < len(t1); i++ {
				a, b := t1[i], t2[i]
				if a.TUID != b.TUID {
					t.Errorf("a.TUID = %v, want %v", a.TUID, b.TUID)
				}
				if *a.Fish != *b.Fish {
					t.Errorf("a.Fish = %v, want %v", a.Fish, b.Fish)
				}
			}
		}
	}
}

func TestHistoryAdd(t *testing.T) {
	const small, big = 50, 100
	const wait = 10 * time.Millisecond
	smallFish := func() *Fish { return &Fish{Weight: small, Length: small} }
	bigFish := func() *Fish { return &Fish{Weight: big, Length: big} }
	for _, p := range parameters {
		s := strconv.Itoa(int(p))
		history := &History{}
		recordc := history.Records(p)
		add := func(id int64, f *Fish) {
			go func() { time.Sleep(time.Millisecond); history.Add(id, f) }()
		}
		t.Run("small_fish_record_"+s, func(t *testing.T) {
			const id = 123
			add(id, smallFish())
			select {
			case r := <-recordc:
				if r.TUID != id || *r.Fish != *smallFish() {
					t.Errorf("bad record")
				}
			case <-time.After(wait):
				t.Errorf("want new record")
			}
		})
		t.Run("big_fish_record_"+s, func(t *testing.T) {
			const id = 456
			add(id, bigFish())
			select {
			case r := <-recordc:
				if r.TUID != id || *r.Fish != *bigFish() {
					t.Errorf("bad record")
				}
			case <-time.After(wait):
				t.Errorf("want new record")
			}
		})
		t.Run("small_fish_no_record_"+s, func(t *testing.T) {
			const id = 789
			add(id, smallFish())
			select {
			case <-recordc:
				t.Errorf("don't want new record")
			case <-time.After(wait):
			}
		})
	}
}

func TestParam(t *testing.T) {
	f := RandomFish()
	l := f.Length
	w := f.Weight
	x := f.Price()
	if p := param(f, Length); p != l {
		t.Errorf("p = %v, want %v", p, l)
	}
	if p := param(f, Weight); p != w {
		t.Errorf("p = %v, want %v", p, w)
	}
	if p := param(f, Price); p != x {
		t.Errorf("p = %v, want %v", p, x)
	}
}
