package fun

import (
	"fmt"
	"regexp"
	"time"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Date struct{}

var dateRe = regexp.MustCompile("^!дата")

func (h *Date) Match(c tele.Context) bool {
	return dateRe.MatchString(c.Text())
}

func (h *Date) Handle(c tele.Context) error {
	t := time.Now()
	y, m, d := t.Date()
	w := t.Weekday()
	mn := monthName(m)
	wn := weekdayName(w)
	s := fmt.Sprintf("📅 Сегодня %s, %d %s %d г.", wn, d, mn, y)
	return c.Send(s)
}

func weekdayName(w time.Weekday) string {
	days := [...]string{
		"воскресенье",
		"понедельник",
		"вторник",
		"среда",
		"четверг",
		"пятница",
		"суббота",
	}
	return days[w]
}

func monthName(m time.Month) string {
	months := [...]string{
		"января",
		"февраля",
		"марта",
		"апреля",
		"мая",
		"июня",
		"июля",
		"августа",
		"сентября",
		"октября",
		"ноября",
		"декабря",
	}
	return months[m-1]
}
