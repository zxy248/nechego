package clock

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type Clock struct {
	Hours   int
	Minutes int
}

func (c Clock) String() string {
	if c.Hours > 0 {
		return fmt.Sprintf("%d ч %d мин", c.Hours, c.Minutes)
	}
	return fmt.Sprintf("%d мин", c.Minutes)
}

func Now() Clock {
	return FromTime(time.Now())
}

func FromTime(t time.Time) Clock {
	return Clock{t.Hour(), t.Minute()}
}

var clockRegexp = regexp.MustCompile(`(\d?\d):(\d\d)`)

func FromString(s string) (Clock, error) {
	match := clockRegexp.FindStringSubmatch(s)
	if match == nil {
		return Clock{}, fmt.Errorf("clock: cannot parse: %s", s)
	}
	h, m := match[1], match[2]
	hours, err := strconv.Atoi(h)
	if err != nil {
		return Clock{}, fmt.Errorf("clock: cannot parse hours: %s", h)
	}
	minutes, err := strconv.Atoi(m)
	if err != nil {
		return Clock{}, fmt.Errorf("clock: cannot parse minutes: %s", m)
	}
	c := Clock{hours, minutes}
	return c, c.Valid()
}

func (c1 Clock) Sub(c2 Clock) Clock {
	if c1.Minutes < c2.Minutes {
		c1.Hours--
		c1.Minutes += 60
	}
	if c1.Hours < c2.Hours {
		c1.Hours += 24
	}
	return Clock{c1.Hours - c2.Hours, c1.Minutes - c2.Minutes}
}

func (c Clock) Valid() error {
	if c.Hours > 23 {
		return fmt.Errorf("clock: invalid hours: %d", c.Hours)
	}
	if c.Minutes > 59 {
		return fmt.Errorf("clock: invalid minutes: %d", c.Minutes)
	}
	return nil
}
