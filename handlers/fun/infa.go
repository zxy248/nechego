package fun

import (
	"fmt"
	"math/rand"
	"nechego/handlers"

	tele "gopkg.in/telebot.v3"
)

type Infa struct{}

var infaRe = handlers.NewRegexp("^!(инфа|вероятность)(.*)")

func (h *Infa) Match(c tele.Context) bool {
	return infaRe.MatchString(c.Text())
}

func (h *Infa) Handle(c tele.Context) error {
	return c.Send(formatInfa(parseInfa(c.Text())))
}

func parseInfa(s string) string {
	return infaRe.FindStringSubmatch(s)[2]
}

func formatInfa(s string) string {
	templates := [...]string{
		"Здравый смысл говорит мне о том, что",
		"Благодаря чувственному опыту я определил, что",
		"Я думаю, что",
		"Используя диалектическую логику, я пришел к выводу, что",
		"Проведя некие изыскания, я выяснил, что",
		"Я провёл мысленный эксперимент и выяснил, что",
		"Мои интеллектуальные потуги привели меня к тому, что",
		"С помощью фактов и логики я доказал, что",
		"Как показывает практика,",
		"Прикинув раз на раз, я определился с тем, что",
		"Уверяю вас в том, что",
	}
	t := templates[rand.Intn(len(templates))]
	p := rand.Intn(101)
	suffix := fmt.Sprintf(" с вероятностью %d%%", p)
	return t + s + suffix
}
