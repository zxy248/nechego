package fun

import (
	"fmt"
	"nechego/handlers"
	"strconv"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Time struct{}

func (*Time) Match(s string) bool {
	_, _, ok := parseTime(s)
	return ok
}

func (*Time) Handle(c tele.Context) error {
	h, m, _ := parseTime(c.Message().Text)
	dh, dm := clockUntil(h, m)
	return c.Send(formatClock(dh, dm))
}

func parseTime(s string) (hour, min int, ok bool) {
	re := handlers.Regexp(`!время до (\d?\d):(\d\d)`)
	ss := re.FindStringSubmatch(s)
	if ss == nil {
		return 0, 0, false
	}

	var err error
	hour, err = strconv.Atoi(ss[1])
	if err != nil {
		return 0, 0, false
	}
	min, err = strconv.Atoi(ss[2])
	if err != nil {
		return 0, 0, false
	}
	if hour > 23 || min > 59 {
		return 0, 0, false
	}
	return hour, min, true
}

func nextClock(h, m int) time.Time {
	t := time.Now()
	r := updateClock(t, h, m)
	if t.Hour() > h || t.Hour() == h && t.Minute() > m {
		return r.Add(24 * time.Hour)
	}
	return r
}

func updateClock(t time.Time, h, m int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), h, m, 0, 0, time.Local)
}

func clockDuration(d time.Duration) (hour, min int) {
	h, m := d/time.Hour, d/time.Minute%60
	return int(h), int(m)
}

func clockUntil(h, m int) (dh, dm int) {
	d := nextClock(h, m).Sub(time.Now().Truncate(time.Minute))
	return clockDuration(d)
}

func formatClock(h, m int) string {
	if h > 0 {
		return fmt.Sprintf("%d ч %d мин", h, m)
	}
	return fmt.Sprintf("%d мин", m)
}
