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

func FromTime(t time.Time) Clock {
	return Clock{t.Hour(), t.Minute()}
}

func FromString(s string) (Clock, error) {
	c, err := parseClock(s)
	if err != nil {
		return c, err
	}
	return c, c.validate()
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

func (c Clock) validate() error {
	if c.Hours < 0 || c.Hours > 23 {
		return fmt.Errorf("clock: invalid hours: %d", c.Hours)
	}
	if c.Minutes < 0 || c.Minutes > 59 {
		return fmt.Errorf("clock: invalid minutes: %d", c.Minutes)
	}
	return nil
}

type matchedClock struct {
	hours   string
	minutes string
}

var clockRe = regexp.MustCompile(`(\d?\d):(\d\d)`)

func matchClock(s string) (matchedClock, error) {
	match := clockRe.FindStringSubmatch(s)
	if match == nil {
		return matchedClock{}, fmt.Errorf("clock: cannot match: %s", s)
	}
	return matchedClock{hours: match[1], minutes: match[2]}, nil
}

func buildClock(c matchedClock) (Clock, error) {
	h, err := strconv.Atoi(c.hours)
	if err != nil {
		return Clock{}, fmt.Errorf("clock: non-number hours: %s", c.hours)
	}
	m, err := strconv.Atoi(c.minutes)
	if err != nil {
		return Clock{}, fmt.Errorf("clock: non-number minutes: %s", c.minutes)
	}
	return Clock{h, m}, nil
}

func parseClock(s string) (Clock, error) {
	c, err := matchClock(s)
	if err != nil {
		return Clock{}, err
	}
	return buildClock(c)
}
