package handlers

import (
	"fmt"
	"html"
	"math/rand"
	"nechego/game"
	"nechego/teleutil"
	"regexp"

	tele "gopkg.in/telebot.v3"
)

type Infa struct{}

var infaRe = regexp.MustCompile("!инфа (.*)")

func (h *Infa) Regexp() *regexp.Regexp {
	return infaRe
}

func (h *Infa) Handle(c tele.Context) error {
	templates := [...]string{
		"Здравый смысл говорит мне о том, что %s с вероятностью %d%%",
		"Благодаря чувственному опыту я определил, что %s с вероятностью %d%%",
		"Я думаю, что %s с вероятностью %d%%",
		"Используя диалектическую логику, я пришел к выводу, что %s с вероятностью %d%%",
		"Проведя некие изыскания, я высяснил, что %s с вероятностью %d%%",
		"Я провел мысленный экперимент и выяснил, что %s с вероятностью %d%%",
		"Мои интеллектуальные потуги привели меня к тому, что %s с вероятностью %d%%",
		"С помощью фактов и логики я доказал, что %s с вероятностью %d%%",
		"Как показывает практика, %s с вероятностью %d%%",
		"Прикинув раз на раз, я определился с тем, что %s с вероятностью %d%%",
		"Уверяю вас в том, что %s с вероятностью %d%%",
	}
	tmpl := templates[rand.Intn(len(templates))]
	arg := h.Regexp().FindStringSubmatch(c.Message().Text)[1]
	return c.Send(fmt.Sprintf(tmpl, arg, rand.Intn(101)))
}

type Who struct {
	Universe *game.Universe
}

var whoRe = regexp.MustCompile("!кто (.*)")

func (h *Who) Regexp() *regexp.Regexp {
	return whoRe
}

func (h *Who) Handle(c tele.Context) error {
	w, err := h.Universe.World(c.Chat().ID)
	if err != nil {
		return err
	}
	w.Lock()
	defer w.Unlock()

	u := w.RandomUser()
	arg := h.Regexp().FindStringSubmatch(c.Message().Text)[1]
	msg := fmt.Sprintf("%s %s", teleutil.Mention(c, w.TGID, u.TUID), html.EscapeString(arg))
	return c.Send(msg, tele.ModeHTML)
}
