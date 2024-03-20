package fun

import (
	"fmt"
	"regexp"
	"time"

	tele "gopkg.in/zxy248/telebot.v3"
)

type NewYear struct{}

var newYearRe = regexp.MustCompile("^!новый год")

func (h *NewYear) Match(c tele.Context) bool {
	return newYearRe.MatchString(c.Text())
}

func (h *NewYear) Handle(c tele.Context) error {
	newYear := time.Date(time.Now().Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
	format := "<b>🎄 До Нового года осталось:<code>%s</code></b>"
	out := fmt.Sprintf(format, buildInterval(time.Until(newYear)))
	return c.Send(out, tele.ModeHTML)
}

func buildInterval(d time.Duration) string {
	var r string

	days := int(d.Hours()) / 24
	if days > 0 {
		r += fmt.Sprintf(" %d д.", days)
	}

	hr := int(d.Hours()) % 24
	if hr > 0 {
		r += fmt.Sprintf(" %d ч.", hr)
	}

	min := int(d.Minutes()) % 60
	if min > 0 {
		r += fmt.Sprintf(" %d мин.", min)
	}

	sec := int(d.Seconds()) % 60
	if sec > 0 {
		r += fmt.Sprintf(" %d сек.", sec)
	}

	return r
}
