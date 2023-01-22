package phone

import "testing"

func TestNumberString(t *testing.T) {
	table := []struct {
		num Number
		str string
	}{
		{0, "00-00-00"},
		{10101, "01-01-01"},
		{111111, "11-11-11"},
		{123, "00-01-23"},
		{123456, "12-34-56"},
		{456000, "45-60-00"},
		{990000, "99-00-00"},
		{999999, "99-99-99"},
	}
	for _, x := range table {
		if x.num.String() != x.str {
			t.Errorf("num == %v, want %v", x.num, x.str)
		}
		num, err := MakeNumber(x.str)
		if err != nil {
			t.Fatal(err)
		}
		if num != x.num {
			t.Errorf("num == %v, want %v", num, x.num)
		}
	}
}

func TestRandomNumber(t *testing.T) {
	for i := 0; i < 100000; i++ {
		const want = 8
		num := RandomNumber()
		if len(num.String()) != want {
			t.Errorf("len == %v, want %v", len(num.String()), want)
		}
	}
}
