package handlers

import tele "gopkg.in/telebot.v3"

type Help struct{}

var helpRe = Regexp("^!(помощь|команды|документ)")

func (h *Help) Match(c tele.Context) bool {
	return helpRe.MatchString(c.Text())
}

func (h *Help) Handle(c tele.Context) error {
	return c.Send("📖 <b>Документация:</b> nechego.pages.dev.", tele.ModeHTML)
}
