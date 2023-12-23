package fun

import (
	"fmt"
	"regexp"
	"time"

	tele "gopkg.in/telebot.v3"
)

type NewYear struct{}

var newYearRe = regexp.MustCompile("^!новый год")

func (h *NewYear) Match(c tele.Context) bool {
	return newYearRe.MatchString(c.Text())
}

func (h *NewYear) Handle(c tele.Context) error {
	t1 := time.Now()
	t2 := time.Date(t1.Year()+1, time.January, 1, 0, 0, 0, 0, time.Local)
	d := t2.Sub(t1)

	r := ""
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
	s := fmt.Sprintf("<b>🎄 До Нового года осталось:<code>%s</code></b>", r)
	return c.Send(s, tele.ModeHTML)
}
